package instant

import (
	"bytes"
	"encoding/xml"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	matchContent = regexp.MustCompile("((?:<p>|<figure[^>]*>).*?(?:</p>|</figure>))")
)

// Article struct
type Article struct {
	Prefix string `xml:"prefix,attr"`
	Lang   string `xml:"lang,attr"`

	Head head `xml:"head"`
	Body body `xml:"body"`
}

type head struct {
	Link link   `xml:"link"`
	Meta []Meta `xml:"meta"`
}

type body struct {
	Article article `xml:"article"`
}

type element struct {
	P      string  `xml:",innerxml"`
	Figure *Figure `xml:"figure,omitempty"`
}

type article struct {
	Header   header    `xml:"header"`
	Elements []element `xml:" ,"`
	Footer   footer    `xml:"footer"`
}

// Header represents instant article header
type header struct {
	H1      string    `xml:"h1"`
	Time    []Time    `xml:"time"`
	H2      string    `xml:"h2,omitempty"`
	H3      *h3       `xml:"h3,omitempty"`
	Address []address `xml:"address,omitempty"`
	Figure  []*Figure `xml:"figure,omitempty"`
}

// Footer represents instant article footer
type footer struct {
	Aside string `xml:"aside,omitempty"`
	Small string `xml:"small,omitempty"`
}

// link for canonical link tag
type link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

// Figure stuct
type Figure struct {
	Img        *Img    `xml:"img,omitempty"`
	IFrame     *IFrame `xml:"iframe,omitempty"`
	Video      *Video  `xml:"video,omitempty"`
	Figcaption string  `xml:"figcaption,omitempty"`
	Class      string  `xml:"class,attr,omitempty"`
}

// Video for article
type Video struct {
	Source source `xml:"source"`
}

// video source struct
type source struct {
	Src  string `xml:"src,attr"`
	Type string `xml:"type,attr"`
}

// IFrame selement
type IFrame struct {
	Src    string `xml:"src,attr,omitempty"`
	Height string `xml:"height,attr,omitempty"`
	Width  string `xml:"width,attr,omitempty"`
	Style  string `xml:"style,attr,omitempty"`
	Hidden string `xml:"hidden,attr,omitempty"`
	Text   string `xml:",innerxml"`
}

// Meta element
type Meta struct {
	Charset  string `xml:"charset,attr,omitempty"`
	Property string `xml:"property,attr,omitempty"`
	Content  string `xml:"content,attr,omitempty"`
}

// Address struct is a plaseholder for article author
type address struct {
	A    alink  `xml:"a"`
	Text string `xml:",chardata"`
}

// a link for Address struct
type alink struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Href string `xml:"href,attr,omitempty"`
	Text string `xml:",chardata"`
}

// h3 struct for article kick
type h3 struct {
	Class string `xml:"class,attr"`
	Text  string `xml:",chardata"`
}

// Time struct for publish date and modified date
type Time struct {
	Text     string `xml:",chardata"`
	Class    string `xml:"class,attr"`
	Datetime string `xml:"datetime,attr"`
}

// Img struct for images
type Img struct {
	Src string `xml:"src,attr"`
}

// MarshalXML for xml.Marshaler interface, marshal Article struct to Facebook Instant Article format.
func (a Article) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// check required fields
	if a.Body.Article.Header.H1 == "" {
		return errors.New("Article title <h1> is required")
	}
	if a.Head.Link.Href == "" {
		return errors.New("Canonical link is required")
	}

	html := struct {
		Prefix string `xml:"prefix,attr"`
		Lang   string `xml:"lang,attr"`
		Head   struct {
			S string `xml:",innerxml"`
		} `xml:"head"`
		Body struct {
			Article article `xml:"article"`
		} `xml:"body"`
	}{}

	html.Body.Article = a.Body.Article
	if a.Lang != "" {
		html.Lang = a.Lang
	} else {
		html.Lang = "en"
	}
	html.Prefix = "op: http://media.facebook.com/op#"
	html.Head.S = a.headString()

	start.Name.Local = "html" // rename root element from Article to html
	e.EncodeToken(xml.Directive("doctype html"))
	return e.EncodeElement(html, start)
}

// MarshalXML for elements in body article content
func (el *element) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	switch {
	case el.P != "":
		p := struct {
			S string `xml:",innerxml"`
		}{S: el.P}
		start.Name.Local = "p"
		return e.EncodeElement(p, start)
	case el.Figure != nil:
		start.Name.Local = "figure"
		return e.EncodeElement(el.Figure, start)
	}
	return nil
}

// headString returns instant article head section with self-closing <meta/> and <link/> tags.
// xml.Marshal by default doesn't produce self-closing tags
func (a *Article) headString() string {
	var buff bytes.Buffer

	// link element
	buff.WriteString("\n<link href=\"")
	buff.WriteString(a.Head.Link.Href)
	buff.WriteString("\" rel=\"")
	buff.WriteString(a.Head.Link.Rel)
	buff.WriteString("\" />")

	// meta elements
	var charset, markupVersion bool
	for _, m := range a.Head.Meta {
		buff.WriteString("\n<meta ")
		if m.Charset != "" {
			buff.WriteString("charset=\"")
			buff.WriteString(m.Charset)
			buff.WriteString("\" ")
			charset = true
		}
		if m.Property != "" {
			buff.WriteString("property=\"")
			buff.WriteString(m.Property)
			buff.WriteString("\" ")
			if m.Property == "op:markup_version" {
				markupVersion = true
			}
		}
		if m.Content != "" {
			buff.WriteString("content=\"")
			buff.WriteString(m.Content)
			buff.WriteString("\" ")
		}
		buff.WriteString("/>")
	}

	// write default meta elements if not set in meta tags
	if !charset {
		buff.WriteString("\n<meta charset=\"utf-8\" />")
	}
	if !markupVersion {
		buff.WriteString("\n<meta property=\"op:markup_version\" content=\"v1.0\" />")
	}

	return buff.String()
}

// SetTitle sets article title.
// Setting article title is mandatory.
func (a *Article) SetTitle(title string) {
	a.Body.Article.Header.H1 = title
}

// SetCanonical sets public web url for this article.
// Setting canonical link is required.
func (a *Article) SetCanonical(url string) {
	a.Head.Link.Rel = "canonical"
	a.Head.Link.Href = url
}

// SetLang sets two-chars language code for article. Default is set to en.
func (a *Article) SetLang(lang string) {
	a.Lang = lang
}

// SetSubtitle sets article subtitle.
func (a *Article) SetSubtitle(subtitle string) {
	a.Body.Article.Header.H2 = subtitle
}

// SetKick sets article kick text.
func (a *Article) SetKick(kick string) {
	a.Body.Article.Header.H3 = &h3{
		Class: "op-kicker",
		Text:  kick,
	}
}

// SetStyle set user defined style for this article.
// See https://developers.facebook.com/docs/instant-articles/guides/design#style for more info.
func (a *Article) SetStyle(style string) {
	a.Head.Meta = append(a.Head.Meta, Meta{
		Property: "fb:article_style",
		Content:  style,
	})
}

// SetContent of instant article.
// All html should be in <p> ... </p> elements. If no <p> elements found, entire HTML param will be added as one paragraph.
// Images and videos can be added later using InsertFigure() but they can also be contained in HTML param if formated properly.
// See https://developers.facebook.com/docs/instant-articles/reference for more info.
func (a *Article) SetContent(html string) {
	if matches := matchContent.FindAllString(html, -1); matches != nil {
		for _, m := range matches {
			switch {
			case strings.HasPrefix(m, "<p"):
				m = strings.TrimPrefix(m, "<p>")
				m = strings.TrimSuffix(m, "</p>")
				a.AddParagraph(m)
			case strings.HasPrefix(m, "<figure"):
				f := &Figure{}
				xml.Unmarshal([]byte(m), f)
				a.AddFigure(f)
			}
		}
	} else {
		a.AddParagraph(html)
	}
}

// AddParagraph to Instant Article.
func (a *Article) AddParagraph(html string) {
	a.Body.Article.Elements = append(a.Body.Article.Elements, element{P: html})
}

// AddAuthor adds article author.
// Link and description are optional (use "")
func (a *Article) AddAuthor(name, link, description string) {
	a.Body.Article.Header.Address = append(a.Body.Article.Header.Address, address{
		A:    alink{Text: name, Href: link},
		Text: description,
	})
}

// SetPublish sets published date.
func (a *Article) SetPublish(date time.Time) {
	// <time class="op-published" datetime="2014-11-11T04:44:16Z">November 11th, 4:44 PM</time>
	a.Body.Article.Header.Time = append(a.Body.Article.Header.Time, Time{
		Class:    "op-published",
		Datetime: date.Format("2006-01-02T15:04:05Z"),
		Text:     date.Format("2006-01-02 15:04:05"),
	})
}

// SetModified sets modified date if article has been modified.
func (a *Article) SetModified(date time.Time) {
	// <time class="op-modified" dateTime="2014-12-11T04:44:16Z">December 11th, 4:44 PM</time>
	a.Body.Article.Header.Time = append(a.Body.Article.Header.Time, Time{
		Class:    "op-modified",
		Datetime: date.Format("2006-01-02T15:04:05Z"),
		Text:     date.Format("2006-01-02 15:04:05"),
	})
}

// SetCoverImage of instant article.
// Caption can be empty string.
func (a *Article) SetCoverImage(url, caption string) {
	if url != "" {
		a.Body.Article.Header.Figure = append(a.Body.Article.Header.Figure, &Figure{
			Img:        &Img{Src: url},
			Figcaption: caption,
		})
	}
}

// SetCoverVideo of instant article.
// videoType in format video/mp4 See the list of supported formats here https://www.facebook.com/help/218673814818907
// Caption can be empty string.
func (a *Article) SetCoverVideo(url, videoType, caption string) {
	if url != "" {
		a.Body.Article.Header.Figure = append(a.Body.Article.Header.Figure, &Figure{
			Video: &Video{
				Source: source{
					Src:  url,
					Type: videoType,
				},
			},
			Figcaption: caption,
		})
	}
}

// SetFooter sets text in footer
// You can use <p> in credits
func (a *Article) SetFooter(credits, copyright string) {
	a.Body.Article.Footer.Aside = credits
	a.Body.Article.Footer.Small = copyright
}

// SetTrackerCode for 3rd party analytics (Google Analytics for example)
// Visit https://developers.facebook.com/docs/instant-articles/reference/analytics for more info.
func (a *Article) SetTrackerCode(code string) {
	f := &Figure{
		Class: "op-tracker",
		IFrame: &IFrame{
			Text: code,
		},
	}
	a.Body.Article.Elements = append(a.Body.Article.Elements, element{Figure: f})
}

// SetTrackerURL for 3rd party analytics that can be included with url.
// Visit https://developers.facebook.com/docs/instant-articles/reference/analytics for more info.
func (a *Article) SetTrackerURL(url string) {
	f := &Figure{
		Class: "op-tracker",
		IFrame: &IFrame{
			Src: url,
		},
	}
	a.Body.Article.Elements = append(a.Body.Article.Elements, element{Figure: f})
}

// switchAutomaticAd positioning by Facebook and choose to manually position ads in article content
func (a *Article) switchAutomaticAd(on bool) {
	a.Head.Meta = append(a.Head.Meta, Meta{
		Property: "fb:use_automatic_ad_placement",
		Content:  strconv.FormatBool(on),
	})
}

// SetAutomaticAd in header that Facebook will place automatically in article
func (a *Article) SetAutomaticAd(src string, width, height int, style, code string) {
	f := adFigure(src, width, height, style, code)
	a.Body.Article.Header.Figure = append(a.Body.Article.Header.Figure, f)
	a.switchAutomaticAd(true)
}

// InsertAd manually on position between paragraphs.
func (a *Article) InsertAd(position int, src string, width, height int, style, code string) {
	a.InsertFigure(position, adFigure(src, width, height, style, code))
	a.switchAutomaticAd(false)
}

// AddAd manually in article content.
func (a *Article) AddAd(position int, src string, width, height int, style, code string) {
	f := adFigure(src, width, height, style, code)
	a.switchAutomaticAd(false)
	a.AddFigure(f)
}

// AddFigure to article content
func (a *Article) AddFigure(f *Figure) {
	a.Body.Article.Elements = append(a.Body.Article.Elements, element{Figure: f})
}

// InsertFigure in content on specified position within existing elements (paragraphs)
func (a *Article) InsertFigure(position int, f *Figure) {
	e := element{Figure: f}
	if position >= len(a.Body.Article.Elements) {
		position = len(a.Body.Article.Elements)
	}
	a.Body.Article.Elements = append(a.Body.Article.Elements[:position], append([]element{e}, a.Body.Article.Elements[position:]...)...)
}

// adFigre create figure with ad
func adFigure(src string, width, height int, style, code string) *Figure {
	return &Figure{
		Class: "op-ad",
		IFrame: &IFrame{
			Width:  strconv.Itoa(width),
			Height: strconv.Itoa(height),
			Src:    src,
			Style:  style,
			Text:   code,
		},
	}
}

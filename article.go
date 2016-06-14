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
	regexContent = regexp.MustCompile("((?:<p>|<figure[^>]*>).*?(?:</p>|</figure>))")
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
	A    a      `xml:"a"`
	Text string `xml:",chardata"`
}

// a link for Address struct
type a struct {
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
func (ia Article) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// check required fields
	if ia.Body.Article.Header.H1 == "" {
		return errors.New("Article title <h1> is required")
	}
	if ia.Head.Link.Href == "" {
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

	html.Body.Article = ia.Body.Article
	if ia.Lang != "" {
		html.Lang = ia.Lang
	} else {
		html.Lang = "en"
	}
	html.Prefix = "op: http://media.facebook.com/op#"
	html.Head.S = ia.headString()

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
func (ia *Article) headString() string {
	var buff bytes.Buffer

	// link element
	buff.WriteString("\n<link href=\"")
	buff.WriteString(ia.Head.Link.Href)
	buff.WriteString("\" rel=\"")
	buff.WriteString(ia.Head.Link.Rel)
	buff.WriteString("\" />")

	// meta elements
	var charset, markupVersion bool
	for _, m := range ia.Head.Meta {
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
func (ia *Article) SetTitle(title string) {
	ia.Body.Article.Header.H1 = title
}

// SetCanonical sets public web url for this article.
// Setting canonical link is required.
func (ia *Article) SetCanonical(url string) {
	ia.Head.Link.Rel = "canonical"
	ia.Head.Link.Href = url
}

// SetLang sets two-chars language code for article. Default is set to en.
func (ia *Article) SetLang(lang string) {
	ia.Lang = lang
}

// SetSubtitle sets article subtitle.
func (ia *Article) SetSubtitle(subtitle string) {
	ia.Body.Article.Header.H2 = subtitle
}

// SetKick sets article kick text.
func (ia *Article) SetKick(kick string) {
	ia.Body.Article.Header.H3 = &h3{
		Class: "op-kicker",
		Text:  kick,
	}
}

// SetStyle set user deifined style for this article.
func (ia *Article) SetStyle(style string) {
	ia.Head.Meta = append(ia.Head.Meta, Meta{
		Property: "fb:article_style",
		Content:  style,
	})
}

// SetContent sets Instant article content.
// All html should be in <p> ... </p> elements. If no <p> elements found, entire HTML will be added as one and only paragraph.
// Images and videoos can be added later using InsertFigure() but <figure> can also be in html
// See more https://developers.facebook.com/docs/instant-articles/reference
func (ia *Article) SetContent(html string) {
	matches := regexContent.FindAllString(html, -1)
	if matches != nil {
		for _, m := range matches {
			switch {
			case strings.HasPrefix(m, "<p"):
				m = strings.TrimPrefix(m, "<p>")
				m = strings.TrimSuffix(m, "</p>")
				ia.AddParagraph(m)
			case strings.HasPrefix(m, "<figure"):
				f := &Figure{}
				xml.Unmarshal([]byte(m), f)
				ia.AddFigure(f)
			}
		}
	} else {
		ia.AddParagraph(html)
	}
}

// AddParagraph add text paragraph to Instant Article.
func (ia *Article) AddParagraph(html string) {
	ia.Body.Article.Elements = append(ia.Body.Article.Elements, element{P: html})
}

// AddAuthor adds article author.
// Link and description can be empty strings ""
func (ia *Article) AddAuthor(name, link, description string) {
	ia.Body.Article.Header.Address = append(ia.Body.Article.Header.Address, address{
		A:    a{Text: name, Href: link},
		Text: description,
	})
}

// SetPublish sets published date.
func (ia *Article) SetPublish(date time.Time) {
	// <time class="op-published" datetime="2014-11-11T04:44:16Z">November 11th, 4:44 PM</time>
	ia.Body.Article.Header.Time = append(ia.Body.Article.Header.Time, Time{
		Class:    "op-published",
		Datetime: date.Format("2006-01-02T15:04:05Z"),
		Text:     date.Format("2006-01-02 15:04:05"),
	})
}

// SetModified sets modified date if article has been modified.
func (ia *Article) SetModified(date time.Time) {
	// <time class="op-modified" dateTime="2014-12-11T04:44:16Z">December 11th, 4:44 PM</time>
	ia.Body.Article.Header.Time = append(ia.Body.Article.Header.Time, Time{
		Class:    "op-modified",
		Datetime: date.Format("2006-01-02T15:04:05Z"),
		Text:     date.Format("2006-01-02 15:04:05"),
	})
}

// SetCoverImage sets image url and image caption.
// Caption can be empty string.
func (ia *Article) SetCoverImage(url, caption string) {
	// override if cover video has been set
	if url != "" {
		ia.Body.Article.Header.Figure = append(ia.Body.Article.Header.Figure, &Figure{
			Img:        &Img{Src: url},
			Figcaption: caption,
		})
	}
}

// SetCoverVideo sets cover video.
// videoType in format video/mp4 See the list of supported formats here https://www.facebook.com/help/218673814818907
// Caption can be empty string.
func (ia *Article) SetCoverVideo(url, videoType, caption string) {
	// override cover image if set
	if url != "" {
		ia.Body.Article.Header.Figure = append(ia.Body.Article.Header.Figure, &Figure{
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
func (ia *Article) SetFooter(credits, copyright string) {
	ia.Body.Article.Footer.Aside = credits
	ia.Body.Article.Footer.Small = copyright
}

// SetTrackerCode for 3rd party analytics (Google Analytics for example)
// Visit https://developers.facebook.com/docs/instant-articles/reference/analytics for more info.
func (ia *Article) SetTrackerCode(html string) {
	f := &Figure{
		Class: "op-tracker",
		IFrame: &IFrame{
			Text: html,
		},
	}
	ia.Body.Article.Elements = append(ia.Body.Article.Elements, element{Figure: f})
}

// SetTrackerURL for 3rd party analytics that can be included with url.
// Visit https://developers.facebook.com/docs/instant-articles/reference/analytics for more info.
func (ia *Article) SetTrackerURL(url string) {
	f := &Figure{
		Class: "op-tracker",
		IFrame: &IFrame{
			Src: url,
		},
	}
	ia.Body.Article.Elements = append(ia.Body.Article.Elements, element{Figure: f})
}

// switchAutomaticAd positioning by Facebook and choose to manually position ads in article content
func (ia *Article) switchAutomaticAd(on bool) {
	ia.Head.Meta = append(ia.Head.Meta, Meta{
		Property: "fb:use_automatic_ad_placement",
		Content:  strconv.FormatBool(on),
	})
}

// SetAutomaticAd in header that Facebook will place automatically in article
func (ia *Article) SetAutomaticAd(src string, width, height int, style, code string) {
	f := adFigure(src, width, height, style, code)
	ia.Body.Article.Header.Figure = append(ia.Body.Article.Header.Figure, f)
	ia.switchAutomaticAd(true)
}

// InsertAd manually on position between paragraphs.
func (ia *Article) InsertAd(position int, src string, width, height int, style, code string) {
	ia.InsertFigure(position, adFigure(src, width, height, style, code))
	ia.switchAutomaticAd(false)
}

// AddAd manually in article content.
func (ia *Article) AddAd(position int, src string, width, height int, style, code string) {
	f := adFigure(src, width, height, style, code)
	ia.switchAutomaticAd(false)
	ia.AddFigure(f)
}

// AddFigure to article content
func (ia *Article) AddFigure(f *Figure) {
	ia.Body.Article.Elements = append(ia.Body.Article.Elements, element{Figure: f})
}

// InsertFigure in content on specified position within existing elements (paragraphs)
func (ia *Article) InsertFigure(position int, f *Figure) {
	e := element{Figure: f}
	if position >= len(ia.Body.Article.Elements) {
		position = len(ia.Body.Article.Elements)
	}
	ia.Body.Article.Elements = append(ia.Body.Article.Elements[:position], append([]element{e}, ia.Body.Article.Elements[position:]...)...)
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

// 	ia.Body.Article
// 	<figure class="op-ad">
//   <iframe width="320" height="50" style="border:0; margin:0;" src="https://www.facebook.com/adnw_request?placement=141956036215488_141956099548815&adtype=banner320x50"></iframe>
// </figure>
// }

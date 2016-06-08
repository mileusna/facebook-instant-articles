package instant

import (
	"bytes"
	"encoding/xml"
	"strings"
	"time"
)

const (
	helperMeta = "META HELPER FIX FOR SELF CLOSING TAGS"
	helperLink = "LINK HELPER FIX FOR SELF CLOSING TAGS"
)

type Article struct {
	Prefix string `xml:"prefix,attr"`
	Lang   string `xml:"lang,attr"`

	Head Head `xml:"head"`
	Body Body `xml:"body"`
}

type Head struct {
	Link link   `xml:"link"`
	Meta []Meta `xml:"meta"`
}

type Body struct {
	Article article `xml:"article"`
}

type article struct {
	Header Header   `xml:"header"`
	P      []string `xml:"p"`
	Footer Footer   `xml:"footer"`
}

// Header represents instant article header
type Header struct {
	H1      string    `xml:"h1"`
	Time    []Time    `xml:"time"`
	H2      string    `xml:"h2,omitempty"`
	H3      h3        `xml:"h3,omitempty"`
	Address []Address `xml:"address,omitempty"`
	Figure  *Figure   `xml:"figure,omitempty"`
}

// Footer represents instant article footer
type Footer struct {
	Aside string `xml:"aside,omitempty"`
	Small string `xml:"small,omitempty"`
}

// link for canonical link tag
type link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

type Figure struct {
	Img        *Img    `xml:"img,omitempty"`
	Iframe     *Iframe `xml:"iframe,omitempty`
	Video      *Video  `xml:"video,omitempty`
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

type Iframe struct {
	Src    string `xml:"src,attr"`
	Height string `xml:"height,attr"`
	Width  string `xml:"width,attr"`
	Hidden string `xml:"hidden,attr"`
}

type Meta struct {
	Charset  string `xml:"charset,attr,omitempty"`
	Property string `xml:"property,attr,omitempty"`
	Content  string `xml:"content,attr,omitempty"`
}

// Address struct is a plaseholder for article author
type Address struct {
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

type Img struct {
	Src string `xml:"src,attr"`
}

// NewArticle creates instant article with mandatory headers all set up
func NewArticle() *Article {
	a := &Article{
		Prefix: "op: http://media.facebook.com/op#",
		Lang:   "en",
		Head: Head{
			Meta: []Meta{
				Meta{
					Charset: "utf-8",
				},
				Meta{
					Property: "op:markup_version",
					Content:  "v1.0",
				},
			},
		},
	}
	return a
}

// MarshalXML for xml.Marshaler interface, marshal Article struct to Facebook Instant Article format.
func (ia *Article) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	html := struct {
		Prefix string `xml:"prefix,attr"`
		Lang   string `xml:"lang,attr"`
		Head   struct {
			S string `xml:",innerxml"`
		} `xml:"head"`
		Body Body `xml:"body"`
	}{}

	html.Body = ia.Body
	html.Lang = ia.Lang
	html.Prefix = ia.Prefix
	html.Head.S = ia.headString()

	e.EncodeToken(xml.Directive("doctype html"))
	e.EncodeToken(xml.CharData("\n"))

	start.Name.Local = "html" // rename root element from Article to html
	return e.EncodeElement(html, start)
}

// headString returns instant article head section with self-closing <meta/> and <link/> tags.
// xml.Marshal by default doesn't produce self-closing tags
func (ia *Article) headString() string {
	var buff bytes.Buffer

	// link element
	buff.WriteString("<link href=\"")
	buff.WriteString(ia.Head.Link.Href)
	buff.WriteString("\" rel=\"")
	buff.WriteString(ia.Head.Link.Rel)
	buff.WriteString("\" />")

	// meta elements
	for _, m := range ia.Head.Meta {
		buff.WriteString("<meta ")
		if m.Charset != "" {
			buff.WriteString("charset=\"")
			buff.WriteString(m.Charset)
			buff.WriteString("\" ")
		}
		if m.Property != "" {
			buff.WriteString("property=\"")
			buff.WriteString(m.Property)
			buff.WriteString("\" ")
		}
		if m.Content != "" {
			buff.WriteString("content=\"")
			buff.WriteString(m.Content)
			buff.WriteString("\" ")
		}
		buff.WriteString("/>")
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
	ia.Body.Article.Header.H3 = h3{
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

// SetText sets Instant article content.
// Text is plain text with new line (\n) for each paragraph.
func (ia *Article) SetText(text string) {
	paragraphs := strings.Split(text, "\n")
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			ia.AddParagraph(p)
		}
	}
}

// AddParagraph add text paragraph to Instant Article.
func (ia *Article) AddParagraph(p string) {
	ia.Body.Article.P = append(ia.Body.Article.P, p)
}

// AddAuthor adds article author.
// Link and description can be empty strings ""
func (ia *Article) AddAuthor(name, link, description string) {
	ia.Body.Article.Header.Address = append(ia.Body.Article.Header.Address, Address{
		A:    a{Text: name, Href: link},
		Text: description,
	})
}

// SetPublish sets published date.
// Use time.Parse() format for formatting user readable date.
func (ia *Article) SetPublish(date time.Time, format string) {
	// <time class="op-published" datetime="2014-11-11T04:44:16Z">November 11th, 4:44 PM</time>
	ia.Body.Article.Header.Time = append(ia.Body.Article.Header.Time, Time{
		Class:    "op-published",
		Datetime: date.Format("2006-01-02T15:04:05Z"),
		Text:     date.Format(format),
	})
}

// SetModified sets modified date if article has been modified.
// Use time.Parse() format for formatting user readable date.
func (ia *Article) SetModified(date time.Time, format string) {
	// <time class="op-modified" dateTime="2014-12-11T04:44:16Z">December 11th, 4:44 PM</time>
	ia.Body.Article.Header.Time = append(ia.Body.Article.Header.Time, Time{
		Class:    "op-modified",
		Datetime: date.Format("2006-01-02T15:04:05Z"),
		Text:     date.Format(format),
	})
}

// SetCoverImage sets image url and image caption.
// Caption can be empty string.
func (ia *Article) SetCoverImage(url, caption string) {
	// override if cover video has been set
	ia.Body.Article.Header.Figure = &Figure{
		Img:        &Img{Src: url},
		Figcaption: caption,
	}
}

// SetCoverVideo sets cover video.
// videoType in format video/mp4 See the list of supported formats here https://www.facebook.com/help/218673814818907
// Caption can be empty string.
func (ia *Article) SetCoverVideo(url, videoType, caption string) {
	// override cover image if set
	ia.Body.Article.Header.Figure = &Figure{
		Video: &Video{
			Source: source{
				Src:  url,
				Type: videoType,
			},
		},
		Figcaption: caption,
	}
}

// SetFooter sets text in footer
// You can use <p> in credits
func (ia *Article) SetFooter(credits, copyright string) {
	ia.Body.Article.Footer.Aside = credits
	ia.Body.Article.Footer.Small = copyright
}

func (ia *Article) SetVideo(url, videoType string) {

}

func (ia *Article) SetTracker() {

}

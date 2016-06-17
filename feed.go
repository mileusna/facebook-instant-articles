package instant

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
)

// Feed struct represents Facebook Instant Articles RSS feed.
// See https://developers.facebook.com/docs/instant-articles/publishing/setup-rss-feed for more info.
type Feed struct {
	Version string  `xml:"version,attr"`
	Content string  `xml:"xmlns:content,attr"`
	Channel channel `xml:"channel"`
}

type channel struct {
	Title         string `xml:"title"`
	LastBuildDate string `xml:"lastBuildDate"`
	Language      string `xml:"language"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Item          []item `xml:"item"`
}

type item struct {
	Title       string   `xml:"title"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	Author      []string `xml:"author"`
	PubDate     string   `xml:"pubDate"`
	Encoded     []byte   `xml:",innerxml"`
}

// SetTitle of feed. Optional.
func (f *Feed) SetTitle(s string) {
	f.Channel.Title = s
}

// SetLink for feed. Optional.
func (f *Feed) SetLink(url string) {
	f.Channel.Link = url
}

// SetDescription set feed description. Optional.
func (f *Feed) SetDescription(s string) {
	f.Channel.Description = s
}

// SetLanguage sets feed language, default is set to en-us.
func (f *Feed) SetLanguage(l string) {
	f.Channel.Language = l
}

// // SetLastBuildDate sets feed last build date.
// func (f *Feed) setLastBuildDate(d time.Time) {
// }

// AddArticle to feed. Md5 checksum of URL will be used as GUID.
func (f *Feed) AddArticle(a Article) error {
	return f.addArticle(a, "")
}

// AddArticleWithGUID to feed.
func (f *Feed) AddArticleWithGUID(a Article, guid string) error {
	return f.addArticle(a, guid)
}

func (f *Feed) addArticle(a Article, guid string) error {

	b, err := xml.Marshal(a)
	if err != nil {
		return err
	}

	if guid == "" {
		guid = fmt.Sprintf("%x", md5.Sum([]byte(a.Head.Link.Href)))
	}

	var buff bytes.Buffer
	buff.WriteString("\n<content:encoded><![CDATA[\n")
	buff.Write(b)
	buff.WriteString("\n]]></content:encoded>")

	i := item{
		Title:       a.Body.Article.Header.H1,
		Description: a.Body.Article.Header.H2,
		Link:        a.Head.Link.Href,
		GUID:        guid,
		Encoded:     buff.Bytes(),
	}

	// if no subtitle, set first para as description
	if i.Description == "" && a.Body.Article.Elements != nil {
		for _, e := range a.Body.Article.Elements {
			if e.P != "" {
				i.Description = e.P
				break
			}
		}
	}

	// set latest article date as pubDate
	for _, d := range a.Body.Article.Header.Time {
		if d.Datetime > i.PubDate {
			i.PubDate = d.Datetime
		}
	}

	// set latest article date to feed LastBuildDate
	if i.PubDate > f.Channel.LastBuildDate {
		f.Channel.LastBuildDate = i.PubDate
	}

	for _, auth := range a.Body.Article.Header.Address {
		i.Author = append(i.Author, auth.A.Text)
	}

	f.Channel.Item = append(f.Channel.Item, i)

	return nil
}

// MarshalXML for xml.Marshaler interface
func (f Feed) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	feed := struct {
		Version string  `xml:"version,attr"`
		Content string  `xml:"xmlns:content,attr"`
		Channel channel `xml:"channel"`
	}{
		Version: f.Version,
		Content: f.Content,
		Channel: f.Channel,
	}

	feed.Version = "2.0"
	feed.Content = "http://purl.org/rss/1.0/modules/content/"
	if feed.Channel.Language == "" {
		feed.Channel.Language = "en-us"
	}

	start.Name.Local = "rss" // rename root element from Article to html
	return e.EncodeElement(feed, start)
}

// RSS is synonym for xml.Marshal(f)
func (f Feed) RSS() ([]byte, error) {
	return xml.Marshal(f)
}

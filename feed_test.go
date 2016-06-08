package instant_test

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/mileusna/facebook-instant-articles"
)

func TestFeed(t *testing.T) {

	f := instant.NewFeed("My site title", "http://mysite.com", "News from all around the world")

	a := instant.NewArticle()

	// mandatory
	a.SetTitle("My article title")
	a.SetCanonical("http://mysite/url-to-this-article")

	// optional
	a.SetSubtitle("My article subtitle")
	a.SetKick("Set article kick")

	a.SetPublish(time.Now(), "02.01.2006") // 02.01.2006 is time.Parse format

	a.SetFooter("", "(C)MyComp 2016")
	a.AddAuthor("Michael", "http://facebook.com/mmichael", "Guest writter")

	a.SetText("Plain text\nPlain text\nPlain text")
	a.AddParagraph("End")

	f.AddArticle(a)
	f.AddArticle(a)

	_, err := xml.Marshal(f)
	if err != nil {
		t.Error(err)
	}
}

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

	a.SetPublish(time.Now())

	a.SetFooter("", "©2016 MyComp")
	a.AddAuthor("Michael", "http://facebook.com/mmichael", "Guest writter")

	a.SetContent("Some plain content")
	a.AddParagraph("End")

	f.AddArticle(a)
	f.AddArticle(a)

	_, err := xml.Marshal(f)
	if err != nil {
		t.Error(err)
	}
}

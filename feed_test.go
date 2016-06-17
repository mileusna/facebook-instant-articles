package instant_test

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/mileusna/facebook-instant-articles"
)

func TestFeed(t *testing.T) {

	var f instant.Feed

	// I don't believe Facebook cares about feed title, link and description, but you can set them
	f.SetTitle("Title feed")
	f.SetLink("http://www.mysite.com")
	f.SetDescription("Feed description")

	var a instant.Article
	a.SetTitle("My article title")
	a.SetCanonical("http://mysite/url-to-this-article")
	a.SetPublish(time.Now())
	a.SetContent("<p>Pragraph 1</p><p>paragraph 2</p>")

	f.AddArticle(a)
	f.AddArticleWithGUID(a, "12333") // add article and use for example mysql id as GUID

	_, err := xml.Marshal(f)
	if err != nil {
		t.Error(err)
	}
}

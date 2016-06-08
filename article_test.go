package instant_test

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/mileusna/facebook-instant-articles"
)

func TestArticle(t *testing.T) {

	a := instant.NewArticle()

	// mandatory
	a.SetTitle("My article title")
	a.SetCanonical("http://mysite/url-to-this-article")

	// optional
	a.SetSubtitle("My article subtitle")
	a.SetKick("Set article kick")

	a.SetPublish(time.Now(), "02.01.2006") // 02.01.2006 is time.Parse format

	a.SetFooter("", "(C)MyComp 2016")
	a.AddAuthor("Michael", "writer", "Guest writter")

	a.SetText("Plain text\nPlain text\nPlain text")
	a.AddParagraph("End")

	_, err := xml.Marshal(a)
	if err != nil {
		t.Error(err)
	}
}

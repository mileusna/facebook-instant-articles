# Facebook Instant Articles Go package [![GoDoc](https://godoc.org/github.com/mileusna/facebook-instant-articles?status.svg)](https://godoc.org/github.com/mileusna/facebook-instant-articles)

Package Instant enables creation and publishing of Facebook Instant Articles.

**Work in progress, things might change**

Facebook instant articles documentation can be found here
https://developers.facebook.com/docs/instant-articles

## Article

Struct instant.Article represents Facebook instant article as described on
https://developers.facebook.com/docs/instant-articles/guides/articlecreate
Use instant.Article and then use helper functions like SetTitle(), SetCoverImage() etc.
to easily create instant article without setting struct properties directly. Elemenets are
nested in inner struct, so using setter functions is much easier. Custom marshaler will generate
Facebook instant article valid html..

```Go
package main

import (
	"encoding/xml"
	"time"
	"github.com/mileusna/facebook-instant-articles"
)

func main() {
	a := instant.Article{}

	// required
	a.SetTitle("My article title")
	a.SetCanonical("http://mysite/url-to-this-article")
	a.SetPublish(time.Now())
	a.SetContent("<p>My content</p><p>Other paragraph</p>")

	// optional
	a.SetLang("fr") // default is en
	a.SetSubtitle("My article subtitle")
	a.SetKick("Exclusive")
	a.SetFooter("", "(C)2016 MyComp")
	a.AddAuthor("Michael", "http://facebook.com/mmichael", "Guest writter")

	html, err := xml.Marshal(a)
	// html, err := a.HTML() // synonym for xml.Marshal
	if err != nil {
		return
	}
	
    // html contains Facebook instant article html as []byte
}
```

##Feed

Struct instant.Feed represents Facebook instant article RSS feed as described on
https://developers.facebook.com/docs/instant-articles/publishing/setup-rss-feed
Use instant.Feed{} and AddArticles to add instant.Article to feed. Custom
marshaler provides valid Facebook instant article RSS feed.

```Go
package main

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/mileusna/facebook-instant-articles"
)

func main() {
	http.HandleFunc("/instant-articles/", handleInstantArticles)
	http.ListenAndServe(":8080", nil)
}

func handleInstantArticles(w http.ResponseWriter, r *http.Request) {
	var f instant.Feed
	// I don't believe Facebook cares about feed title, link and description, but if you like...
	f.SetTitle("Title feed")
	f.SetLink("http://www.mysite.com")
	f.SetDescription("Feed description")

	var a instant.Article
	a.SetTitle("My article title")
	a.SetCanonical("http://mysite/url-to-this-article")
	a.SetPublish(time.Now())
	a.SetContent("<p>Pragraph 1</p><p>paragraph 2</p>")
	f.AddArticle(a) // add article, GUID will be hash of URL

	a = instant.Article{}
	a.SetTitle("My other title")
	a.SetCanonical("http://mysite/url-to-this-article-number-two")
	a.SetPublish(time.Now())
	a.SetContent("<p>Pragraph 1</p><p>paragraph 2</p>")
	f.AddArticleWithGUID(a, "12333") // add article and use for example mysql id as GUID

	feed, err := xml.Marshal(f)
	// html, err := f.RSS() // synonym for xml.Marshal
	if err != nil {
		return
	}
	w.Write(feed)
}
```

###License

MIT
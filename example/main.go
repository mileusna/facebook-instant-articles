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
	// I don't believe Facebook cares about feed title, link and description, but you can set them
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
	if err != nil {
		return
	}

	w.Write(feed)
}

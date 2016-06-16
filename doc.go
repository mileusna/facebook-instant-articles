/*
Package instant enables creation and publishing of Facebook Instant Articles.

Full Facebook instant articles documentation can be found here
https://developers.facebook.com/docs/instant-articles

Struct instant.Article represents Facebook instant article as described on
https://developers.facebook.com/docs/instant-articles/guides/articlecreate
Use instant.Article and then use helper functions like SetTitle(), SetCoverImage() etc.
to easily create instant article without setting struct properties directly. Elemenets are
nested in inner struct, so using setter functions is much easier. Custom marshaler will generate
Facebook instant article valid html.

Example:
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
        a.SetSubtitle("My article subtitle")
        a.SetKick("Exclusive")
        a.SetFooter("", "(C)2016 MyComp")
        a.AddAuthor("Michael", "http://facebook.com/mmichael", "Guest writter")

        html, err := xml.Marshal(a)
        if err != nil {
            return
        }
        // html contains Facebook instant article html as []byte
    }

Struct instant.Feed represents Facebook instant article RSS feed as described on
https://developers.facebook.com/docs/instant-articles/publishing/setup-rss-feed
Use instant.NewFeed() to create initial struct with all headers set up and
use helper functions like AddArticles to add instant.Article to feed. Custom
marshaler provides valid Facebook instant article RSS feed.

Example:
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
*/
package instant

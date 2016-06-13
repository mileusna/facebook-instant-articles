/*
Package instant enables creation and publishing of Facebook Instant Articles.

Full Facebook instant articles documentation can be found here
https://developers.facebook.com/docs/instant-articles

Struct instant.Article represents Facebook instant article as described on
https://developers.facebook.com/docs/instant-articles/guides/articlecreate
Use instant.NewArticle() to create initial struct with all headers set up and
use helper functions like SetTitle(), SetCoverImage() to easily
create instant article without setting struct properties directly. Custom
marshaler provides Facebook instant article valid html.

Example:
    package main

    import (
        "encoding/xml"
        "time"
        "github.com/mileusna/facebook-instant-articles"
    )

    func main() {
        a := instant.NewArticle()

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
*/
package instant

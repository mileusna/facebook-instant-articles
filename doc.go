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

        // mandatory
        a.SetTitle("My article title")
        a.SetCanonical("http://mysite/url-to-this-article")

        // optional
        a.SetSubtitle("My article subtitle")
        a.SetKick("Set article kick")

        a.SetPublish(time.Now(), "02.01.2006") // 02.01.2006 is time.Parse format

        a.SetFooter("", "(C)MyComp 2016")
        a.AddAuthor("Michael", "http://facebook.com/mmichael", "Guest writter")

        a.SetContent("<p>my content</p><p>Use Paragraphs</p>")

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

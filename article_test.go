package instant_test

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/mileusna/facebook-instant-articles"
)

func TestArticle(t *testing.T) {

	a := instant.Article{}

	// required
	a.SetTitle("My article title")
	a.SetCanonical("http://mysite/url-to-this-article")

	// optional
	a.SetSubtitle("My article subtitle")
	a.SetKick("Exclusive!")
	a.SetLang("fr") // default is en

	a.SetPublish(time.Now())

	a.SetFooter("", "©MyComp 2016")
	a.AddAuthor("Michael", "http://facebook.com/mmichael", "Guest writter")

	a.SetTrackerCode(`<script>
		(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
		(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
		m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
		})(window,document,'script','https://www.google-analytics.com/analytics.js','ga');
		ga('create', 'UA-22233-16', 'auto');
		ga('send', 'pageview');
	</script>`)

	// <figure class="op-ad">
	//   <iframe width="320" height="50" style="border:0; margin:0;" src="https://www.facebook.com/adnw_request?placement=141956036215488_141956099548815&adtype=banner320x50"></iframe>
	// </figure>

	a.SetContent("<p>Plain text</p><figure><iframe>a = 0;</iframe></figure><p>Plain <b>text</b><p>Plain text 22</p>")
	a.AddParagraph("End <strong>The</strong>")
	a.SetAutomaticAd("https://www.facebook.com/adnw_request?placement=141956036215488_141956099548815&adtype=banner320x50", 320, 50, "border:0; margin:0;", "")
	//a.InsertAd(7, "https://www.facebook.com/adnw_request?placement=141956036215488_141956099548815&adtype=banner320x50", 320, 50, "border:0; margin:0;", "")

	content, err := xml.MarshalIndent(a, "", "    ")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(content))
}

package instant

/*
   Package Instant enables creation and publishing of Facebook Instant Articles.
   https://developers.facebook.com/docs/instant-articles

   Struct instant.Article represents Facebook instant article as described on
   https://developers.facebook.com/docs/instant-articles/guides/articlecreate
   Use instant.NewArticle() to creat initial struct with all headers set up and
   use helper functions line SetTitle(), SetCoverImage() to easily
   create instant article without setting struct properties directly. Custom
   marshaler provides Facebook instant article html.

   Struct instant.Feed represents Facebook instant article RSS feed as described on
   https://developers.facebook.com/docs/instant-articles/publishing/setup-rss-feed
   Use instant.NewFeed() to create initial struct with all headers set up and
   use helper functions line AddArticles to add instant.Article to feed. Custom
   marshaler provides valid Facebook instant article RSS feed.

*/

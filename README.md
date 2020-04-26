A quick Go hack to parse Twitter's JSON archive and delete all Tweets.

Made because all alternatives appear to be online services and that looks
like a very dumb idea to give access to some random website :-).
And I needed some Go project to start using it and Go is great to deal with
JSON.

**tweetstosql** parses the `tweet.js` file you get and dumps all the tweets
`id` to a Postgresql database (modify to whatever else you use)

**twitterwipe** queries the database and issues the API calls to delete the
tweets

And versions for likes (which use a different ID). Could be easily merged into
a single utility but this was originally a hack and just cleaned it a bit for
this release. You can practice your Go skills improving this :-).

Just request your Twitter archive, and point the SQL extractors to the `tweet.js`
and `like.js` files that are included in the archive.

You need a Twitter API key. I had an old one, not sure it's still easy to request
one or you need to justify and wait a bit more for it.

I was afraid of API rate limitations as described in documentation but I was
able to delete some 70k tweets without any problems. It took a while since
this is single threaded and single request since otherwise I could hit
some rate limitation on Twitter side.

Sorry, no go mod support yet, didn't have time yet to understand how to use
that stuff.

External main dependencies (and whatever else each project uses):
- [libpq](https://github.com/lib/pq)
- [oauth1](https://github.com/dghubble/oauth1)
- [go-twitter](https://github.com/dghubble/go-twitter/twitter)

Do
```
make deps
```
To install all dependencies (GOPATH used to everything self-contained)

To build you can use the Makefile
```
make all
```

or
```
make tweetstosql
make twitterwipe
make likestosql
make likeswipe
```

Edit the source files to add your Postgresql database information.

Have fun,

fG!

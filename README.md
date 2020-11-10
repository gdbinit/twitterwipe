### About

An utility to parse Twitter's JSON archives and delete tweets, likes and direct messages.

Made because all alternatives appear to be online services and that looks like a very dumb idea to give access to some random website :-).

Have fun,

fG!

### Requirements

Just request your Twitter archive (**Settings -> Your Account -> Download an archive of your data**), and point the import command to the proper files included in the archive.

* `tweet.js`: contains all your tweets and retweets
* `like.js`: contains all your likes
* `direct-messages.js`: contains all your direct messags threads

You need a Twitter API key from [Twitter developer](https://developer.twitter.com). I had an old one, not sure if it's still easy to request or you need to justify and wait for authorization.

You will need to create a `config.json` file with the API keys information. Hardcoded credentials is not a good idea :P.

```
{
    "consumerkey": "INSERT CONSUMER KEY HERE",
    "consumersecret": "INSERT CONSUMER SECRET HERE",
    "accesstoken": "INSERT ACCESS TOKEN HERE",
    "accesstokensecret": "INSERT ACCESS TOKEN SECRET HERE"
}
```

The `config.json` should be on the same folder as the main binary.

I was afraid of API rate limitations as described in documentation but I was able to delete some 70k tweets without any problems. It took a while since this is single threaded and single request since otherwise I could hit some rate limitation on Twitter side.

The `delete` command has a `-l` flag to limit the number of items deleted if you manage to hit limits.

### Building

To build you can use the Makefile
```
make
```

Or

```go
go build -o twitterwipe main.go
```

Tested only in macOS, should work in other OSes because Go promises that.

### Usage

#### Import

First step is to import the ids into the local database using the `import` command.

```bash
twitterwipe import tweets tweet.js
```

Available commands to `import` are `tweets`, `likes`, `dms`.

The files in the archive are not straight JSON parseable by Go `encoding/json` package but the code deals with it so you can use those files without a problem.

#### Delete

The next step is to remove the data from Twitter using the `delete` command.

```bash
twitterwipe delete tweets
```

Available commands to `delete` are `tweets`, `likes`, `dms`.

The `-l` or `--limit` flag is available to limit the maximum number of items to delete (in case if you hit Twitter rate limitation, which I never did).

#### Status

The `status` command displays how many ids currently exist in the database.

```bash
twitterwipe status
```

### TODO

Maybe...

- [] Delete by individual ID
- [] Import IDs via Twitter API (rate limited afaik!)

### Bill of Materials

- [oauth1](https://github.com/dghubble/oauth1)
- [go-twitter](https://github.com/dghubble/go-twitter/twitter)
- [progressbar](https://github.com/schollz/progressbar)
- [cobra](https://github.com/spf13/cobra)
- [boltdb](https://github.com/boltdb/bolt)

### References

- [Twitter Rate Limits](https://developer.twitter.com/en/docs/basics/rate-limits)
- [Tweet objects](https://developer.twitter.com/en/docs/tweets/data-dictionary/overview/user-object)
- [Sample Tweet Object](https://github.com/twitterdev/tweet-updates/blob/master/samples/initial/compatibility_extended_13996.json)

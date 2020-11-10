/*
 * Created by fG! on 08/11/2020.
 * Copyright © 2020 Pedro Vilaca. All rights reserved.
 * reverser@put.as - https://reverse.put.as
 *
 * All advertising materials mentioning features or use of this software must display
 * the following acknowledgement: This product includes software developed by
 * Pedro Vilaca.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted provided that the following conditions are met:
 
 * 1. Redistributions of source code must retain the above copyright notice, this list
 * of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice, this
 * list of conditions and the following disclaimer in the documentation and/or
 * other materials provided with the distribution.
 * 3. All advertising materials mentioning features or use of this software must
 * display the following acknowledgement: This product includes software developed
 * by Pedro Vilaca.
 * 4. Neither the name of the author nor the names of its contributors may be
 * used to endorse or promote products derived from this software without specific
 * prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR
 * CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
 * EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
 * PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package commands

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/boltdb/bolt"
)

// structure from tweet.js
type tweet struct {
    ID string `json:id`
}
type tweet_entry struct {
    Tweet tweet `json:"tweet"`
}

// structure from like.js
type like struct {
        TweetId string `json:"tweetId"`
}
type like_entry struct {
    Like like `json:"like"`
}

// structure from direct-messages.js
// each direct message has an unique id
type dm_id struct {
	ID string `json:"id"`
}
// each messages array entry is of this type
type dm_msg_create struct {
	MessageCreate dm_id `json:"messageCreate"`
}
// we are interested in the array of messages
type dm_msg struct {
	Messages []dm_msg_create `json:"messages"`
}
// the DMs json is an array of these
type dm_entry struct {
	DM dm_msg `json:"dmConversation"`
}

var importCmd = &cobra.Command{
  Use:   "import",
  Short: "Import tweets, likes, direct messages",
  Long:  `Import tweets, likes, dms to local database for further processing`,
}

var tweetsImportCmd = &cobra.Command{
	Use:	"tweets <path to tweet.js>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[+] Importing tweets...")

		var err error
		db, err := setupDB()
		if err != nil {
			fmt.Printf("[-] Failed to open database: %s\n", err.Error())
			return
		}

		defer db.Close()
		err = importTweets(db, args[0])
		if err != nil {
			fmt.Printf("[-] Failed to import tweets: %s\n", err.Error())
		}
	},
}

var likesImportCmd = &cobra.Command{
	Use:	"likes <path to like.js>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[+] Importing likes...")
		var err error
		db, err := setupDB()
		if err != nil {
			fmt.Printf("[-] Failed to open database: %s\n", err.Error())
			return
		}

		defer db.Close()
		err = importLikes(db, args[0])
		if err != nil {
			fmt.Printf("[-] Failed to import likes: %s\n", err.Error())
		}
	},
}

var dmsImportCmd = &cobra.Command{
	Use:	"dms <path to direct-messages.js",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[+] Importing direct messages...")
		var err error
		db, err := setupDB()
		if err != nil {
			fmt.Printf("[-] Failed to open database: %s\n", err.Error())
			return
		}

		defer db.Close()
		err = importDMs(db, args[0])
		if err != nil {
			fmt.Printf("[-] Failed to import direct messages: %s\n", err.Error())
		}
	},	
}

func init() {
  rootCmd.AddCommand(importCmd)
  importCmd.AddCommand(tweetsImportCmd)
  importCmd.AddCommand(likesImportCmd)
  importCmd.AddCommand(dmsImportCmd)
}

func setupDB() (*bolt.DB, error) {
    db, err := bolt.Open("tweets.db", 0600, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to open database, %v", err)
    }
    err = db.Update(func(tx *bolt.Tx) error {
        _, err = tx.CreateBucketIfNotExists([]byte("TWEETS"))
        if err != nil {
        	return fmt.Errorf("could not create tweets bucket: %v", err)
        }
        _, err = tx.CreateBucketIfNotExists([]byte("LIKES"))
        if err != nil {
        	return fmt.Errorf("could not create likes bucket: %v", err)
        }
        _, err = tx.CreateBucketIfNotExists([]byte("DMS"))
        if err != nil {
        	return fmt.Errorf("could not create dms bucket: %v", err)
    	}
        return nil
    })
    if err != nil {
        return nil, fmt.Errorf("could not set up buckets, %v", err)
    }
    return db, nil
}

func importTweets(db *bolt.DB, tweets_file string) error {
    var err error
    f, err := os.Open(tweets_file)
    if err != nil {
        return err
    }
    defer f.Close()

    // read and JSON decode the file - Go is nice for this!
    dat := bufio.NewReader(f)
    // the original tweet.js from the Twitter archive isn't pure JSON
    // Go can parse, so verify if that header exists and skip it
    // Peek doesn't advance the reader so we don't really change anything here
    header, err := dat.Peek(42)
    if err != nil {
    	return err
    }
    // if the header exists, discards those bytes aka advance the reader
    // so that the JSON parser will only see the good data
    // the first line appears to always be
    // window.YTD.tweet.part0 = [ {
    equalIndex := strings.Index(string(header), "=")
    bracketIndex := strings.Index(string(header), "[")
    // if this header exists we expect the = to happen before
    // the first square bracket
    // if not then it doesn't exist or it's after so we don't care
    // for example if user removed that from the js file already
    if equalIndex > 0 && equalIndex < bracketIndex {
    	// the index is where the character matches so advance one
    	// XXX: no bounds check? YOLO!
    	if _, err = dat.Discard(equalIndex+1); err != nil {
    		return err
    	}
    }
    decoder := json.NewDecoder(dat);
    var content []tweet_entry
    err = decoder.Decode(&content)
    if err != nil {
        return err
    }
    // iterate over the decoded array/slice and extract the tweet id for each entry
    // that's the only info we need to delete
    // for more complex delete actions just check the API documentation
    // and fix the tweet* structs to add the JSON fields you are interested in
    // and then fix the query and database table
    var count int
    err = db.Batch(func(tx *bolt.Tx) error {
    	b := tx.Bucket([]byte("TWEETS"))
    	for _, i := range content {
    		err = b.Put([]byte(i.Tweet.ID), []byte("1"))
    		if err != nil {
    			break 
    		}
    		count++
    	}
    	return err
    })
    if err != nil {
    	// fmt.Println("[-] Failed to update database")
    	return err
    }
    fmt.Printf("[+] Imported %d tweets.\n", count)
    return err
}

func importLikes(db *bolt.DB, input_file string) error {
    var err error
    f, err := os.Open(input_file)
    if err != nil {
       	return err
    }
    defer f.Close()

    dat := bufio.NewReader(f)
    // the original tweet.js from the Twitter archive isn't pure JSON
    // Go can parse, so verify if that header exists and skip it
    // Peek doesn't advance the reader so we don't really change anything here
    header, err := dat.Peek(42)
    if err != nil {
    	return err
    }
    // if the header exists, discards those bytes aka advance the reader
    // so that the JSON parser will only see the good data
    // the first line appears to always be
    // window.YTD.tweet.part0 = [ {
    equalIndex := strings.Index(string(header), "=")
    bracketIndex := strings.Index(string(header), "[")
    // if this header exists we expect the = to happen before
    // the first square bracket
    // if not then it doesn't exist or it's after so we don't care
    // for example if user removed that from the js file already
    if equalIndex > 0 && equalIndex < bracketIndex {
    	// the index is where the character matches so advance one
    	// XXX: no bounds check? YOLO!
    	if _, err = dat.Discard(equalIndex+1); err != nil {
    		return err
    	}
    }
    decoder := json.NewDecoder(dat);
    var content []like_entry
    err = decoder.Decode(&content)
    if err != nil {
        return err
    }

    var count int
    err = db.Batch(func(tx *bolt.Tx) error {
    	b := tx.Bucket([]byte("LIKES"))
    	for _, i := range content {
    		err = b.Put([]byte(i.Like.TweetId), []byte("1"))
    		if err != nil {
    			break 
    		}
    		count++
    	}
    	return err
    })
    if err != nil {
    	// fmt.Println("[-] Failed to update database")
    	return err
    }
    fmt.Printf("[+] Imported %d likes.\n", count)
    return nil
}

func importDMs(db *bolt.DB, input_file string) error {
    var err error
    f, err := os.Open(input_file)
    if err != nil {
    	return err
    }
    defer f.Close()

    dat := bufio.NewReader(f)
    // the original tweet.js from the Twitter archive isn't pure JSON
    // Go can parse, so verify if that header exists and skip it
    // Peek doesn't advance the reader so we don't really change anything here
    header, err := dat.Peek(42)
    if err != nil {
    	return err
    }
    // if the header exists, discards those bytes aka advance the reader
    // so that the JSON parser will only see the good data
    // the first line appears to always be
    // window.YTD.tweet.part0 = [ {
    equalIndex := strings.Index(string(header), "=")
    bracketIndex := strings.Index(string(header), "[")
    // if this header exists we expect the = to happen before
    // the first square bracket
    // if not then it doesn't exist or it's after so we don't care
    // for example if user removed that from the js file already
    if equalIndex > 0 && equalIndex < bracketIndex {
    	// the index is where the character matches so advance one
    	// XXX: no bounds check? YOLO!
    	if _, err = dat.Discard(equalIndex+1); err != nil {
    		return err
    	}
    }
    decoder := json.NewDecoder(dat);
    var content []dm_entry
    err = decoder.Decode(&content)
    if err != nil {
        return err
    }

    var count int
    err = db.Batch(func(tx *bolt.Tx) error {
    	b := tx.Bucket([]byte("DMS"))
    	for _, i := range content {
    		for _, x := range i.DM.Messages {
    			err = b.Put([]byte(x.MessageCreate.ID), []byte("1"))
    			if err != nil {
    				break 
    			}
				count++
    		}
    	}
    	return err
    })
    if err != nil {
    	// fmt.Println("[-] Failed to update database: ", err)
    	return err
    }
    fmt.Printf("[+] Imported %d direct messages.\n", count)
    return nil
}

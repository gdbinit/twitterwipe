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
	"github.com/spf13/cobra"
	"github.com/boltdb/bolt"
)

var statusCmd = &cobra.Command{
  	Use:   "status",
  	Short: "Display database information",
  	Long:  `Show amount of tweets, likes, dms in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("tweets.db", 0600, nil)
		if err != nil {
			fmt.Printf("Failed to open database: %s\n", err.Error())
			return
		}
		defer db.Close()

		var totalTweets int
		var totalLikes int
		var totalDMs int

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("TWEETS"))
			tweetStats := b.Stats()
			totalTweets = tweetStats.KeyN

			b = tx.Bucket([]byte("LIKES"))
			tweetStats = b.Stats()
			totalLikes = tweetStats.KeyN

			b = tx.Bucket([]byte("DMS"))
			tweetStats = b.Stats()
			totalDMs = tweetStats.KeyN

			return nil
		})
		fmt.Printf("[ Database stats ]\n")
		fmt.Printf("Tweets => %d\n", totalTweets)
		fmt.Printf("Likes  => %d\n", totalLikes)
		fmt.Printf("DMs    => %d\n", totalDMs)
	},
}

func init() {
  rootCmd.AddCommand(statusCmd)
}

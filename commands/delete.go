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
	"github.com/boltdb/bolt"
	"github.com/gdbinit/twitterwipe/twitter"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"strconv"
)

var limit int

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete tweets, likes, direct messages",
	Long:  `Delete all tweets, likes, dms`,
}

var tweetsDeleteCmd = &cobra.Command{
	Use: "tweets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[+] Deleting tweets...")
		db, err := bolt.Open("tweets.db", 0600, nil)
		if err != nil {
			fmt.Printf("[-] Failed to open database: %s\n", err.Error())
			return
		}
		defer db.Close()

		client, err := twitter.Open()
		if err != nil {
			fmt.Printf("[-] Failed to get Twitter client instance: %s\n", err.Error())
			return
		}

		totalWork := getBucketEntries(db, "TWEETS")
		// nothing to do
		if totalWork == 0 {
			return
		}

		limit, _ := cmd.Flags().GetInt("limit")
		if limit > 0 && totalWork >= limit {
			totalWork = limit
		}

		bar := progressbar.Default(int64(totalWork))

		var deletedCount int
		toDelete := make([]string, 0)
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("TWEETS"))
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				// key is []byte type
				id_str := string(k)
				id, err := strconv.Atoi(id_str)
				if err != nil {
					panic(err)
				}

				// bye bye tweet!!!!
				err = client.DeleteTweet(int64(id))
				bar.Add(1)
				if err != nil {
					fmt.Printf("[-] Error deleting tweet %s from twitter.com: %s\n", id_str, err.Error())
				} else {
					toDelete = append(toDelete, id_str)
					deletedCount++
					if deletedCount >= totalWork {
						break
					}
				}
			}
			return nil
		})

		if deletedCount > 0 {
			fmt.Printf("Deleting %d tweets from database...\n", deletedCount)
			err = db.Batch(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("TWEETS"))
				for _, k := range toDelete {
					b.Delete([]byte(k))
				}
				return nil
			})
		}
		fmt.Println("[+] All done!")
	},
}

var likesDeleteCmd = &cobra.Command{
	Use: "likes",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("[+] Deleting likes...")
		db, err := bolt.Open("tweets.db", 0600, nil)
		if err != nil {
			fmt.Printf("[-] Failed to open database: %s\n", err.Error())
			return
		}
		defer db.Close()

		client, err := twitter.Open()
		if err != nil {
			fmt.Printf("[-] Failed to get Twitter client instance: %s\n", err.Error())
			return
		}

		totalWork := getBucketEntries(db, "LIKES")
		// nothing to do
		if totalWork == 0 {
			return
		}

		limit, _ := cmd.Flags().GetInt("limit")
		if limit > 0 && totalWork >= limit {
			totalWork = limit
		}

		bar := progressbar.Default(int64(totalWork))

		var deletedCount int
		toDelete := make([]string, 0)
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("LIKES"))
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				id_str := string(k)
				id, err := strconv.Atoi(id_str)
				if err != nil {
					panic(err)
				}
				// This is different from tweets -> FavoriteDestroyParams
				bar.Add(1)
				err = client.DeleteLike(int64(id))
				if err != nil {
					fmt.Printf("[-] Error deleting like %s from twitter.com: %s\n", id_str, err.Error())
				} else {
					toDelete = append(toDelete, id_str)
					deletedCount++
					if deletedCount >= totalWork {
						break
					}
				}
			}
			return nil
		})

		if deletedCount > 0 {
			fmt.Printf("[+] Deleting %d likes from database\n", deletedCount)
			err = db.Batch(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("LIKES"))
				for _, k := range toDelete {
					b.Delete([]byte(k))
				}
				return nil
			})
		}
		fmt.Println("[+] All done!")
	},
}

var dmsDeleteCmd = &cobra.Command{
	Use: "dms",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("[+] Deleting direct messages...")
		db, err := bolt.Open("tweets.db", 0600, nil)
		if err != nil {
			fmt.Printf("[-] Failed to open database: %s\n", err.Error())
			return
		}
		defer db.Close()

		client, err := twitter.Open()
		_ = client
		if err != nil {
			fmt.Printf("[-] Failed to get Twitter client instance: %s\n", err.Error())
			return
		}

		totalWork := getBucketEntries(db, "DMS")
		// nothing to do
		if totalWork == 0 {
			return
		}

		limit, _ := cmd.Flags().GetInt("limit")
		if limit > 0 && totalWork >= limit {
			totalWork = limit
		}

		bar := progressbar.Default(int64(totalWork))

		var deletedCount int
		toDelete := make([]string, 0)
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("DMS"))
			c := b.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				id_str := string(k)
				// This is different from tweets -> FavoriteDestroyParams
				bar.Add(1)
				err = client.DeleteDM(id_str)
				if err != nil {
					fmt.Printf("[-] Error deleting DM %s from twitter.com: %s\n", id_str, err.Error())
				} else {
					toDelete = append(toDelete, id_str)
					deletedCount++
					if deletedCount >= totalWork {
						break
					}
				}
			}
			return nil
		})

		if deletedCount > 0 {
			fmt.Printf("[+] Deleting %d DMs from database\n", deletedCount)
			err = db.Batch(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("DMS"))
				for _, k := range toDelete {
					b.Delete([]byte(k))
				}
				return nil
			})
		}
		fmt.Println("[+] All done!")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(tweetsDeleteCmd)
	deleteCmd.AddCommand(likesDeleteCmd)
	deleteCmd.AddCommand(dmsDeleteCmd)
	deleteCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "limit the number of items to delete [default: unlimited]")
}

func getBucketEntries(db *bolt.DB, bucket string) int {
	var i int
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		s := b.Stats()
		i = s.KeyN
		return nil
	})
	return i
}

/*
 * Created by fG! on 06/08/2021.
 * Copyright © 2021 Pedro Vilaca. All rights reserved.
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
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Removes tweets, likes, direct messages from the database",
	Long:  `Removes tweets, likes, dms from local database`,
}

var tweetsEditCmd = &cobra.Command{
	Use:  "tweet tweetID",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetID := args[0]
		fmt.Printf("[+] Removing tweet with ID %s...\n", targetID)

		var err error
		db, err := bolt.Open("tweets.db", 0600, nil)
		if err != nil {
			fmt.Printf("[-] Failed to open database: %s\n", err.Error())
			return
		}
		defer db.Close()

		totalEntries := getBucketEntries(db, "TWEETS")
		// nothing to do
		if totalEntries == 0 {
			return
		}
		db.Update(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("TWEETS"))
			if err := b.Delete([]byte(targetID)); err != nil {
				fmt.Println("[-] Failed to delete tweet: ", err.Error())
				return err
			}
			return nil
		})
		// the key delete doesn't return an error if they key doesn't exist
		// so we compare the number of entries
		totalNewEntries := getBucketEntries(db, "TWEETS")
		if totalNewEntries < totalEntries {
			fmt.Printf("[+] Removed tweet with ID %s.\n", args[0])
		} else {
			fmt.Printf("[-] Tweet with ID %s does not exist in database.\n", args[0])
		}
	},
}

var likesEditCmd = &cobra.Command{
	Use:  "like likeID",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetID := args[0]
		fmt.Println("[+] Removing like with ID %s...", targetID)
		var err error
		db, err := bolt.Open("tweets.db", 0600, nil)
		if err != nil {
			fmt.Printf("[-] Failed to open database: %s\n", err.Error())
			return
		}
		defer db.Close()

		totalEntries := getBucketEntries(db, "LIKES")
		// nothing to do
		if totalEntries == 0 {
			return
		}
		db.Update(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("LIKES"))
			if err := b.Delete([]byte(targetID)); err != nil {
				fmt.Println("[-] Failed to delete like: ", err.Error())
				return err
			}
			return nil
		})
		// the key delete doesn't return an error if they key doesn't exist
		// so we compare the number of entries
		totalNewEntries := getBucketEntries(db, "LIKES")
		if totalNewEntries < totalEntries {
			fmt.Printf("[+] Removed like with ID %s.\n", args[0])
		} else {
			fmt.Printf("[-] Like with ID %s does not exist in database.\n", args[0])
		}
	},
}

var dmsEditCmd = &cobra.Command{
	Use:  "dm dmID",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetID := args[0]
		fmt.Println("[+] Removing direct message with ID %s...", targetID)
		var err error
		db, err := bolt.Open("tweets.db", 0600, nil)
		if err != nil {
			fmt.Printf("[-] Failed to open database: %s\n", err.Error())
			return
		}
		defer db.Close()

		totalEntries := getBucketEntries(db, "DMS")
		// nothing to do
		if totalEntries == 0 {
			return
		}
		db.Update(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("DMS"))
			if err := b.Delete([]byte(targetID)); err != nil {
				fmt.Println("[-] Failed to delete dm: ", err.Error())
				return err
			}
			return nil
		})
		// the key delete doesn't return an error if they key doesn't exist
		// so we compare the number of entries
		totalNewEntries := getBucketEntries(db, "DMS")
		if totalNewEntries < totalEntries {
			fmt.Printf("[+] Removed dm with ID %s.\n", args[0])
		} else {
			fmt.Printf("[-] Direct message with ID %s does not exist in database.\n", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.AddCommand(tweetsEditCmd)
	editCmd.AddCommand(likesEditCmd)
	editCmd.AddCommand(dmsEditCmd)
}

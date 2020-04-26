/*
 * Created by fG! on 18/02/2020.
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

package main

import(
	"log"
	"flag"
	"fmt"
	"strconv"

    "github.com/dghubble/go-twitter/twitter"
    "github.com/dghubble/oauth1"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
 	consumerKey = "INSERT CONSUMER KEY HERE"
 	consumerSecret = "INSERT CONSUMER SECRET HERE"
 	accessToken = "INSERT ACCESS TOKEN HERE"   
 	accessTokenSecret = "INSERT ACCESS TOKEN SECRET HERE"

 	db_host = "postgresql instance IP/hostname"
 	db_port = 5432
 	db_user = "database username"
 	db_password = "database password"
 	db_database = "database name"
)

var verbose bool

type tweet struct {
		ID string `json:id_str`
}

type tweet_entry struct {
	Tweet tweet `json:tweet`
}

func open_db() (*sql.DB, error) {
    psqlInfo := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        db_host, db_port, db_user, db_password, db_database)
    
    handle, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }
    return handle, nil
}

func main() {

	var amount int
 	flag.BoolVar(&verbose, "v", false, "verbose output")
 	flag.IntVar(&amount, "n", 20, "number of tweets to delete")
 	flag.Parse()

 	// open connection to database and retrieve the id of tweets to delete
 	// I wasn't sure about Twitter API rate limits since API docs says 
 	// they exist for delete
 	// but I managed to delete my whole timeline without rate limits
    db, err := open_db()
    if err != nil {
        log.Fatal("[-] error opening database: ", err)
    }

    sqlStatement := `SELECT id FROM tweets LIMIT $1`
    rows, err := db.Query(sqlStatement, amount)
    if err != nil {
    	log.Fatal("[-] error in select -> ", err)	
    }
    defer rows.Close()

    // connect to Twitter API
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	// Twitter client
	client := twitter.NewClient(httpClient)
	
	// Verify Credentials
	// Not really necessary - just debug so comment out after everything works
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("User's ACCOUNT:\n%+v\n", user)
	}

	// iterate over the rows we received from the database query
	// and issue delete command to Twitter API
    for rows.Next() {
    	var id string
    	err = rows.Scan(&id)
    	if err != nil {
    		panic(err)
    	}
    	fmt.Printf("Deleting Tweet with ID: %s\n", id)
    	// JSON doesn't like 64 bit integers (bahahahha)
    	// so they are stored as strings
    	x, err := strconv.Atoi(id)
    	if err != nil {
    		panic(err)
    	}
    	// bye bye tweet!!!!
    	destroyParams := &twitter.StatusDestroyParams{}
		tweet, _, err := client.Statuses.Destroy(int64(x), destroyParams)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		} else {
			fmt.Printf("Result:\n%+v\n", tweet)
		}
		// remove tweet from database so we don't query it again
		sqlStatement = `DELETE FROM tweets WHERE id=$1`
		_, err = db.Exec(sqlStatement, id)
		if err != nil {
			fmt.Printf("[-] Error deleting tweet %s from database -> %s\n", id, err.Error())
  			return
		}		
    }

    // debug query to check the API limits - never managed to have issues deleting 70k or so tweets
	rateLimits, _, err := client.RateLimits.Status(&twitter.RateLimitParams{Resources: []string{"statuses", "users"}})
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		for k, v := range rateLimits.Resources.Statuses {
			fmt.Printf("Limit Key: %s -> Value: %+v\n", k, v.Limit)
			fmt.Printf("Remaining Key: %s -> Value: %+v\n", k, v.Remaining)
			fmt.Printf("Reset Key: %s -> Value: %+v\n", k, v.Reset)
		}
	}
	db.Close()
}

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
    "os"
    "bufio"
    "log"
    "flag"
    "fmt"
    "encoding/json"
    "database/sql"
    _ "github.com/lib/pq"
)

/*
    CREATE TABLE tweets (id TEXT UNIQUE);
*/
const (
    db_host = "postgresql instance IP/hostname"
    db_port = 5432
    db_user = "database username"
    db_password = "database password"
    db_database = "database name"
)

var verbose bool

type like struct {
        TweetId string `json:tweetId`
}

type like_entry struct {
    Like like `json:like`
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
    var tweets_file string

    flag.StringVar(&tweets_file, "i", "", "path to tweets.js JSON file")
    flag.BoolVar(&verbose, "v", false, "verbose output")
    flag.Parse()
    // we have to do this by ourselves yoooooo... flag package is too basic!
    if tweets_file == "" {
        fmt.Println("[-] Error: missing tweets file")
        fmt.Println("Usage:")
        flag.PrintDefaults()
        os.Exit(1)
    }

    f, err := os.Open(tweets_file)
    if err != nil {
        log.Fatal("[-] failed to open tweets file: ", err)
    }
    defer f.Close()

    db, err := open_db()
    if err != nil {
        log.Fatal("[-] error opening database: ", err)
    }

    // read and JSON decode the file - Go is nice for this!
    dat := bufio.NewReader(f)
    decoder := json.NewDecoder(dat);
    var content []like_entry
    err = decoder.Decode(&content)
    if err != nil {
        log.Fatal(err)
    }
    // iterate over the decoded array/slice and extract the tweet id for each entry
    // that's the only info we need to delete
    // for more complex delete actions just check the API documentation
    // and fix the tweet* structs to add the JSON fields you are interested in
    // and then fix the query and database table
    for _, i := range content {
        fmt.Printf("ID: %s\n", i.Like.TweetId)
        sqlStatement := `INSERT INTO tweets (id) VALUES($1)`
        _, err = db.Exec(sqlStatement, i.Like.TweetId)
        if err != nil {
            fmt.Printf("[-] failed insert: %s\n", err.Error())
        }
    }
}

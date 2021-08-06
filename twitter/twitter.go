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

package twitter

import (
	"encoding/json"
	"fmt"
	t "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"io/ioutil"
)

type Client struct {
	client *t.Client
}

// read the tokens from a config file to avoid hardcoding secrets!
type Secrets struct {
	ConsumerKey       string `json:"consumerkey"`
	ConsumerSecret    string `json:"consumersecret"`
	AccessToken       string `json:"accesstoken"`
	AccessTokenSecret string `json:"accesstokensecret`
}

func Open() (*Client, error) {
	var err error
	dat, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.json")
	}

	s := new(Secrets)
	err = json.Unmarshal(dat, s)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config.json")
	}

	// connect to Twitter API
	config := oauth1.NewConfig(s.ConsumerKey, s.ConsumerSecret)
	token := oauth1.NewToken(s.AccessToken, s.AccessTokenSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	// Twitter client
	c := &Client{}
	c.client = t.NewClient(httpClient)

	// Verify Credentials
	// Not really necessary - just debug so comment out after everything works
	verifyParams := &t.AccountVerifyParams{
		SkipStatus:   t.Bool(true),
		IncludeEmail: t.Bool(true),
	}
	_, _, err = c.client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) DeleteTweet(id int64) error {
	// bye bye tweet!!!!
	destroyParams := &t.StatusDestroyParams{}
	_, _, err := c.client.Statuses.Destroy(id, destroyParams)
	return err
}

func (c *Client) DeleteLike(id int64) error {
	// This is different from tweets -> FavoriteDestroyParams
	destroyParams := &t.FavoriteDestroyParams{ID: id}
	_, _, err := c.client.Favorites.Destroy(destroyParams)
	return err
}

func (c *Client) DeleteDM(id string) error {
	// the Destroy() is deprecated and this is the right one to use
	// uses string instead of int64 as Destroy()
	_, err := c.client.DirectMessages.EventsDestroy(id)
	return err
}

// debug query to check the API limits - never managed to have issues deleting 70k or so tweets
func (c *Client) Limits() {
	rateLimits, _, err := c.client.RateLimits.Status(&t.RateLimitParams{Resources: []string{"statuses", "users"}})
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		for k, v := range rateLimits.Resources.Statuses {
			fmt.Printf("Limit Key: %s -> Value: %+v\n", k, v.Limit)
			fmt.Printf("Remaining Key: %s -> Value: %+v\n", k, v.Remaining)
			fmt.Printf("Reset Key: %s -> Value: %+v\n", k, v.Reset)
		}

		for k, v := range rateLimits.Resources.Users {
			fmt.Printf("Limit Key: %s -> Value: %+v\n", k, v.Limit)
			fmt.Printf("Remaining Key: %s -> Value: %+v\n", k, v.Remaining)
			fmt.Printf("Reset Key: %s -> Value: %+v\n", k, v.Reset)
		}
	}
}

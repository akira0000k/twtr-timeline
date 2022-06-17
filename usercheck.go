package main

import (
	"context"
	//"encoding/json"
	"fmt"
	"os"

	//"github.com/sivchari/gotwtr"
)

func name2id(screen_name string) (id string, err error) {
	id = ""
	res, err := twapi.client.RetrieveSingleUserWithUserName(context.Background(), screen_name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if res == nil {
		fmt.Fprintln(os.Stderr, "res == nil")
		return id, err
	}
	if res.Errors != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", res.Title, res.Detail)
		for _, e := range res.Errors {
			fmt.Fprintf(os.Stderr, "%s: %s(%s)\n", e.Title, e.Detail, e.Message)
		}
		return
	}
	//jsonraw, _ := json.Marshal(res) //test
	//fmt.Println(string(jsonraw))    //test
	user := res.User
	if user != nil {
		id = res.User.ID
	}
	return id, err
}


// func (c *Client) RetrieveSingleUserWithUserName(ctx context.Context, userName string, opt ...*RetrieveUserOption) (*UserResponse, error) {
//  	return retrieveSingleUserWithUserName(ctx, c.client, userName, opt...)
// }
//  
// type UserResponse struct {
//  	User     *User               `json:"data"`
//  	Includes *UserIncludes       `json:"includes,omitempty"`
//  	Errors   []*APIResponseError `json:"errors,omitempty"`
//  	Title    string              `json:"title,omitempty"`
//  	Detail   string              `json:"detail,omitempty"`
//  	Type     string              `json:"type,omitempty"`
// }
//  
// type User struct {
//  	ID              string             `json:"id"`
//  	Name            string             `json:"name"`
//  	UserName        string             `json:"username"`
//  	CreatedAt       string             `json:"created_at,omitempty"`
//  	Description     string             `json:"description,omitempty"`
//  	Entities        *UserEntity        `json:"entities,omitempty"`
//  	Location        string             `json:"location,omitempty"`
//  	PinnedTweetID   string             `json:"pinned_tweet_id,omitempty"`
//  	ProfileImageURL string             `json:"profile_image_url,omitempty"`
//  	Protected       bool               `json:"protected,omitempty"`
//  	PublicMetrics   *UserPublicMetrics `json:"public_metrics,omitempty"`
//  	URL             string             `json:"url,omitempty"`
//  	Verified        bool               `json:"verified,omitempty"`
//  	Withheld        *UserWithheld      `json:"withheld,omitempty"`
// }

// type APIResponseError struct {
//  	Title              string      `json:"title"`
//  	Detail             string      `json:"detail"`
//  	Type               string      `json:"type"`
//  	ResourceType       string      `json:"resource_type"`
//  	ResourceID         string      `json:"resource_id"`
//  	Parameter          string      `json:"parameter"`
//  	Parameters         Parameter   `json:"parameters"`
//  	Message            string      `json:"message"`
//  	Value              interface{} `json:"value"`
//  	Reason             string      `json:"reason,omitempty"`
//  	ClientID           string      `json:"client_id,omitempty"`
//  	RequiredEnrollment string      `json:"required_enrollment,omitempty"`
//  	RegistrationURL    string      `json:"registration_url,omitempty"`
//  	ConnectionIssue    string      `json:"connection_issue,omitempty"`
//  	Status             int         `json:"status,omitempty"`
// }

// {"data":null,
//  "errors":[{"title":"",
//  	    "detail":"",
//  	    "type":"",
//  	    "resource_type":"",
//  	    "resource_id":"",
//  	    "parameter":"",
//  	    "parameters":{"id":null,
//  			  "ids":null,
//  			  "username":["  "],
//  			  "usernames":null},
//  	    "message":"The `username` query parameter value [  ] does not match ^[A-Za-z0-9_]{1,15}$",
//  	    "value":null}],
//  "title":"Invalid Request",
//  "detail":"One or more parameters to your request was invalid.",
//  "type":"https://api.twitter.com/2/problems/invalid-request"}

// retrieve single user with user name: 400 Bad Request https://api.twitter.com/2/users/by/username/++

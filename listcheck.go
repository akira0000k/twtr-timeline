package main

import (
	"encoding/json"
	"context"
	"fmt"
	"strings"
	"os"
	
	"github.com/sivchari/gotwtr"
)

func listIDCheck(userID string, listid string, listname string) (returnID string) {
	fmt.Printf("userID=[%v] listid=[%v] listname=[%v]\n", userID, listid, listname)
	returnID = "0"
	if userID == "0" {
		if listid != "0" {
			return listid
		}
		fmt.Fprintln(os.Stderr, "no userid")
		return
	}
	var lists = []*gotwtr.List{}
	var onetime = 100
	var pagtoken = ""
	for {
		res, err := twapi.client.LookUpAllListsOwned(context.Background(), userID, &gotwtr.AllListsOwnedOption {
			ListFields: []gotwtr.ListField{gotwtr.ListFieldPrivate},
			MaxResults: onetime,
			PaginationToken: pagtoken,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			return
		}
		if res.Meta != nil && res.Meta.ResultCount == 0 {
			break
		}
		if res.Errors != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", res.Title, res.Detail)
			for _, e := range res.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s(%s)\n", e.Title, e.Detail, e.Message)
			}
			return
		}
		if res.Lists == nil || res.Meta == nil {
			jsonUser, _ := json.Marshal(res)
			fmt.Println(string(jsonUser))
			return
		}
		lists = append(lists, res.Lists...)
		if pagtoken = res.Meta.NextToken; pagtoken == "" {
			break
		}
	}
	if len(lists) <= 0 {
		fmt.Fprintln(os.Stderr, "no list in this user.")
		return
	}
	matchcount := 0
	for _, list := range lists {
		if listid != "0" && list.ID == listid ||
			listname != "" && strings.HasPrefix(list.Name, listname) {
			returnID = list.ID
			fmt.Fprintln(os.Stderr, "listId: ", list.ID, " Name: ", list.Name)
			matchcount += 1
		}
	}
	if matchcount == 1 {
		return returnID
	} else if matchcount > 1 {
		fmt.Fprintln(os.Stderr, "choose list id.")
	} else {
		if listid == "0" && listname == "" {
			fmt.Fprintln(os.Stderr, "need -listid or -listname.")
		} else {
			fmt.Fprintln(os.Stderr, "list id or list name unmatch.")
		}
		for _, list := range lists {
			fmt.Fprintln(os.Stderr, "listId: ", list.ID, " Name: ", list.Name)
		}
	}
	return "0"
}

// type AllListsOwnedOption struct {
//  	Expansions      []Expansion
//  	ListFields      []ListField
//  	MaxResults      int
//  	PaginationToken string
//  	UserFields      []UserField
// }
//const (
// 	ListFieldCreatedAt   ListField = "created_at"
// 	ListFollowerCount    ListField = "follower_count"
// 	ListMemberCount      ListField = "member_count"
// 	ListFieldPrivate     ListField = "private"
// 	ListFieldDescription ListField = "description"
// 	ListOwnerID          ListField = "owner_id"
//)
//  
// type AllListsOwnedResponse struct {
//  	Lists    []*List             `json:"data"`
//  	Includes *ListIncludes       `json:"includes,omitempty"`
//  	Errors   []*APIResponseError `json:"errors,omitempty"`
//  	Meta     *ListMeta           `json:"meta"`
//  	Title    string              `json:"title,omitempty"`
//  	Detail   string              `json:"detail,omitempty"`
//  	Type     string              `json:"type,omitempty"`
// }
//  
// type List struct {
//  	ID            string `json:"id"`
//  	Name          string `json:"name"`
//  	CreatedAt     string `json:"created_at,omitempty"`
//  	Private       bool   `json:"private,omitempty"`
//  	FollowerCount int    `json:"follower_count,omitempty"`
//  	MemberCount   int    `json:"member_count,omitempty"`
//  	OwnerID       string `json:"owner_id,omitempty"`
//  	Description   string `json:"description,omitempty"`
// }
//  
// type ListIncludes struct {
//  	Tweets []*Tweet
//  	Users  []*User
// }
//  
// type ListMeta struct {
//  	ResultCount   int    `json:"result_count"`
//  	PreviousToken string `json:"previous_token,omitempty"`
//  	NextToken     string `json:"next_token,omitempty"`
// }
//  
// // LookUpAllListsOwned returns all Lists owned by the specified user.
// func (c *Client) LookUpAllListsOwned(ctx context.Context, userID string, opt ...*AllListsOwnedOption) (*AllListsOwnedResponse, error) {
//  	return lookUpAllListsOwned(ctx, c.client, userID, opt...)
// }
//
// const	lookUpAllListsOwnedURL   = "https://api.twitter.com/2/users/%v/owned_lists"

// userID=[999999] listid=[123] listname=[]
//  
// {"data":null,
//  "errors":[{"title":"Not Found Error",
//  	    "detail":"Could not find user with id: [999999].",
//  	    "type":"https://api.twitter.com/2/problems/resource-not-found",
//  	    "resource_type":"user",
//  	    "resource_id":"999999",
//  	    "parameter":"id",
//  	    "parameters":{"id":null,
//  			  "ids":null,
//  			  "username":null,
//  			  "usernames":null},
//  	    "message":"",
//  	    "value":"999999"}],
//  "meta":null}

// {"data":null,"meta":{"result_count":0}}

// {"data":null
//  "errors":[{"title":""
//  	    "detail":""
//  	    "type":""
//  	    "resource_type":""
//  	    "resource_id":""
//  	    "parameter":""
//  	    "parameters":{"id":["3456 "]
//  			  "ids":null
//  			  "username":null
//  			  "usernames":null}
//  	    "message":"The `id` query parameter value [3456 ] is not valid"
//  	    "value":null}]
//  "meta":null
//  "title":"Invalid Request"
//  "detail":"One or more parameters to your request was invalid."
//  "type":"https://api.twitter.com/2/problems/invalid-request"}

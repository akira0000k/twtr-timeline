package main

import (
	"encoding/json"
	"context"
	"fmt"
	"strconv"
	"time"
	"os"
	"errors"

	"github.com/sivchari/gotwtr"
)


func twid2str(twid twidt) (string) {
	return strconv.FormatInt(int64(twid), 10)
}


type twSearchApi struct {
	client *gotwtr.Client
	t  tltype
	tlopt gotwtr.UserTweetTimelineOption
	tmopt gotwtr.UserMentionTimelineOption
	lsopt gotwtr.ListTweetsOption
	sropt gotwtr.SearchTweetsOption
	nextToken string
	seq int
	jsonp bool
}

func (ta *twSearchApi) getTL(userID string, maxresult int, max twidt, since twidt) (twres *gotwtr.TweetsResponse, count int, last bool, err error) {
	count = 0
	last = true
	err = nil
	client := ta.client
	switch ta.t {
	case tluser:
		if ta.seq == 0 {
			ta.tlopt = gotwtr.UserTweetTimelineOption {
				TweetFields: []gotwtr.TweetField{gotwtr.TweetFieldReferencedTweets, gotwtr.TweetFieldInReplyToUserID},
				Expansions:  []gotwtr.Expansion{gotwtr.ExpansionReferencedTweetsIDAuthorID},}
			if max > 0 {ta.tlopt.UntilID = twid2str(max)}
			if since > 0 {ta.tlopt.SinceID = twid2str(since)}
		} else {
			if ta.nextToken == "" {
				break
			}
			ta.tlopt.PaginationToken = ta.nextToken
		}
		if maxresult > 0 {ta.tlopt.MaxResults = maxresult}
		
		ta.seq++
		ta.nextToken = ""

		var res *gotwtr.UserTweetTimelineResponse
		res, err = client.UserTweetTimeline(context.Background(), userID, &ta.tlopt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			break
		}
		if ta.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Errors != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", res.Title, res.Detail)
			for _, e := range res.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s(%s)\n", e.Title, e.Detail, e.Message)
			}
			if err == nil {err = errors.New("see res.Errors")}
		}
		if res.Meta != nil {
			count = res.Meta.ResultCount
			if res.Meta.NextToken != "" {
				ta.nextToken = res.Meta.NextToken
				last = false
			}
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d\n", time.Now().Format("15:04:05"), count)
		tr := gotwtr.TweetsResponse{}
		tr.Tweets   = res.Tweets
		tr.Includes = res.Includes
		tr.Errors   = res.Errors
		tr.Title    = res.Title
		tr.Detail   = res.Detail
		tr.Type	    = res.Type
		return &tr, count, last, err
	case tlhome:
	case tlmention:
		if ta.seq == 0 {
			ta.tmopt = gotwtr.UserMentionTimelineOption {
				TweetFields: []gotwtr.TweetField{gotwtr.TweetFieldReferencedTweets, gotwtr.TweetFieldInReplyToUserID},
				Expansions:  []gotwtr.Expansion{gotwtr.ExpansionReferencedTweetsIDAuthorID},}
			if max > 0 {ta.tmopt.UntilID = twid2str(max)}
			if since > 0 {ta.tmopt.SinceID = twid2str(since)}
		} else {
			if ta.nextToken == "" {
				break
			}
			ta.tmopt.PaginationToken = ta.nextToken
		}
		if maxresult > 0 {ta.tmopt.MaxResults = maxresult}
		
		ta.seq++
		ta.nextToken = ""

		var res *gotwtr.UserMentionTimelineResponse
		res, err = client.UserMentionTimeline(context.Background(), userID, &ta.tmopt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			break
		}
		if ta.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Errors != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", res.Title, res.Detail)
			for _, e := range res.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s(%s)\n", e.Title, e.Detail, e.Message)
			}
			if err == nil {err = errors.New("see res.Errors")}
		}
		if res.Meta != nil {
			count = res.Meta.ResultCount
			if res.Meta.NextToken != "" {
				ta.nextToken = res.Meta.NextToken
				last = false
			}
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d\n", time.Now().Format("15:04:05"), count)
		tr := gotwtr.TweetsResponse{}
		tr.Tweets   = res.Tweets
		tr.Includes = res.Includes
		tr.Errors   = res.Errors
		tr.Title    = res.Title
		tr.Detail   = res.Detail
		tr.Type	    = res.Type
		return &tr, count, last, err
	case tlrtofme:
	case tllist:
		if ta.seq == 0 {
			ta.lsopt = gotwtr.ListTweetsOption {
				TweetFields: []gotwtr.TweetField{gotwtr.TweetFieldReferencedTweets, gotwtr.TweetFieldInReplyToUserID},
				Expansions:  []gotwtr.Expansion{gotwtr.ExpansionReferencedTweetsIDAuthorID},}
		} else {
			if ta.nextToken == "" {
				break
			}
			ta.lsopt.PaginationToken = ta.nextToken
		}
		if maxresult > 0 {ta.lsopt.MaxResults = maxresult}
		
		ta.seq++
		ta.nextToken = ""

		var res *gotwtr.ListTweetsResponse
		res, err = client.LookUpListTweets(context.Background(), userID, &ta.lsopt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			break
		}
		if ta.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Errors != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", res.Title, res.Detail)
			for _, e := range res.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s(%s)\n", e.Title, e.Detail, e.Message)
			}
			if err == nil {err = errors.New("see res.Errors")}
		}
		if res.Meta != nil {
			count = res.Meta.ResultCount
			if res.Meta.NextToken != "" {
				ta.nextToken = res.Meta.NextToken
				last = false
			}
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d\n", time.Now().Format("15:04:05"), count)
		tr := gotwtr.TweetsResponse{}
		tr.Tweets   = res.Tweets
		tr.Includes = &gotwtr.TweetIncludes{}
		if res.Includes != nil {
			tr.Includes.Tweets = res.Includes.Tweets
			tr.Includes.Users = res.Includes.Users
		}
		tr.Errors   = res.Errors
		tr.Title    = res.Title
		tr.Detail   = res.Detail
		tr.Type	    = res.Type
		return &tr, count, last, err
	case tlsearch: fallthrough
	case tlsearcha:
		query := userID
		if ta.seq == 0 {
			ta.sropt = gotwtr.SearchTweetsOption {
				TweetFields: []gotwtr.TweetField{gotwtr.TweetFieldReferencedTweets, gotwtr.TweetFieldInReplyToUserID},
				Expansions:  []gotwtr.Expansion{gotwtr.ExpansionReferencedTweetsIDAuthorID},}
			if max > 0 {ta.sropt.UntilID = twid2str(max)}
			if since > 0 {ta.sropt.SinceID = twid2str(since)}
		} else {
			if ta.nextToken == "" {
				break
			}
			ta.sropt.NextToken = ta.nextToken
		}
		if maxresult > 0 {ta.sropt.MaxResults = maxresult}
		
		ta.seq++
		ta.nextToken = ""

		var res *gotwtr.SearchTweetsResponse
		switch ta.t {
		case tlsearch:
			res, err = client.SearchRecentTweets(context.Background(), query, &ta.sropt)
		case tlsearcha:
			res, err = client.SearchAllTweets(context.Background(), query, &ta.sropt)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			break
		}
		if ta.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Errors != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", res.Title, res.Detail)
			for _, e := range res.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s(%s)\n", e.Title, e.Detail, e.Message)
			}
			if err == nil {err = errors.New("see res.Errors")}
		}
		if res.Meta != nil {
			count = res.Meta.ResultCount
			if res.Meta.NextToken != "" {
				ta.nextToken = res.Meta.NextToken
				last = false
			}
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d\n", time.Now().Format("15:04:05"), count)
		tr := gotwtr.TweetsResponse{}
		tr.Tweets   = res.Tweets
		tr.Includes = res.Includes
		tr.Errors   = res.Errors
		tr.Title    = res.Title
		tr.Detail   = res.Detail
		tr.Type	    = res.Type
		return &tr, count, last, err
	}
	return nil, 0, true, err
}

func (ta *twSearchApi) rewindQuery() {
	ta.seq = 0
}

// type UserTweetTimelineResponse struct {
//  	Tweets   []*Tweet            `json:"data"`
//  	Includes *TweetIncludes      `json:"includes,omitempty"`
//  	Errors   []*APIResponseError `json:"errors,omitempty"`
//  	Meta     *UserTimelineMeta   `json:"meta"`
//  	Title    string              `json:"title,omitempty"`
//  	Detail   string              `json:"detail,omitempty"`
//  	Type     string              `json:"type,omitempty"`
// }
//  
// type UserMentionTimelineResponse struct {
//  	Tweets   []*Tweet            `json:"data"`
//  	Includes *TweetIncludes      `json:"includes,omitempty"`
//  	Errors   []*APIResponseError `json:"errors,omitempty"`
//  	Meta     *UserTimelineMeta   `json:"meta"`
//  	Title    string              `json:"title,omitempty"`
//  	Detail   string              `json:"detail,omitempty"`
//  	Type     string              `json:"type,omitempty"`
// }
//  
// type UserTimelineMeta struct {
//  	ResultCount int    `json:"result_count"`
//  	NewestID    string `json:"newest_id"`
//  	OldestID    string `json:"oldest_id"`
//  	NextToken   string `json:"next_token"`
// }
//  
// type ListTweetsResponse struct {
//  	Tweets   []*Tweet            `json:"data"`
//  	Includes *ListIncludes       `json:"includes,omitempty"`
//  	Errors   []*APIResponseError `json:"errors,omitempty"`
//  	Meta     *ListMeta           `json:"meta"`
//  	Title    string              `json:"title,omitempty"`
//  	Detail   string              `json:"detail,omitempty"`
//  	Type     string              `json:"type,omitempty"`
// }
//  
// type ListMeta struct {
//  	ResultCount   int    `json:"result_count"`
//  	PreviousToken string `json:"previous_token,omitempty"`
//  	NextToken     string `json:"next_token,omitempty"`
// }

// type SearchTweetsResponse struct {
//  	Tweets   []*Tweet            `json:"data"`
//  	Includes *TweetIncludes      `json:"includes,omitempty"`
//  	Meta     *SearchTweetsMeta   `json:"meta"`
//  	Errors   []*APIResponseError `json:"errors,omitempty"`
//  	Title    string              `json:"title,omitempty"`
//  	Detail   string              `json:"detail,omitempty"`
//  	Type     string              `json:"type,omitempty"`
// }
//  
// type SearchTweetsMeta struct {
//  	ResultCount int    `json:"result_count"`
//  	NewestID    string `json:"newest_id"`
//  	OldestID    string `json:"oldest_id"`
//  	NextToken   string `json:"next_token,omitempty"`
// }

// {"data":null,
//  "errors":[{"title":"Authorization Error",
//  	    "detail":"Sorry,
//  you are not authorized to see the user with id: [1234567890].",
//  	    "type":"https://api.twitter.com/2/problems/not-authorized-for-resource",
//  	    "resource_type":"user",
//  	    "resource_id":"1234567890",
//  	    "parameter":"id",
//  	    "parameters":{"id":null,
//  			  "ids":null,
//  			  "username":null,
//  			  "usernames":null},
//  	    "message":"",
//  	    "value":"1234567890"}],
//  "meta":null}

// {"data":null,
//  "errors":[{"title":"",
//  	    "detail":"",
//  	    "type":"",
//  	    "resource_type":"",
//  	    "resource_id":"",
//  	    "parameter":"",
//  	    "parameters":{"id":["1234567 890"],
//  			  "ids":null,
//  			  "username":null,
//  			  "usernames":null},
//  	    "message":"The `id` query parameter value [1234567 890] is not valid",
//  	    "value":null}],
//  "meta":null,
//  "title":"Invalid Request",
//  "detail":"One or more parameters to your request was invalid.",
//  "type":"https://api.twitter.com/2/problems/invalid-request"}

// user tweet timeline: 400 Bad Request https://api.twitter.com/2/users/12345+67890/tweets?expansions=referenced_tweets.id.author_id&max_results=5&tweet.fields=referenced_tweets%2Cin_reply_to_user_id
// Invalid Request: One or more parameters to your request was invalid.
// : (The `id` query parameter value [12345 67890] is not valid)

// {"data":null,"meta":{"result_count":0,"newest_id":"","oldest_id":"","next_token":""}}

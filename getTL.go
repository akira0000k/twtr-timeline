package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"os"

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
	sropt gotwtr.SearchTweetsOption
	usermeta *gotwtr.UserTimelineMeta
	srchmeta *gotwtr.SearchTweetsMeta
	seq int
}

func (ta *twSearchApi) getTL(userID string, maxresult int, max twidt, since twidt) (twres *gotwtr.TweetsResponse, err error) {
	client := ta.client
	switch ta.t {
	case tluser:
		if ta.seq == 0 {
			ta.tlopt = gotwtr.UserTweetTimelineOption {
				TweetFields: []gotwtr.TweetField{gotwtr.TweetFieldReferencedTweets,gotwtr.TweetFieldInReplyToUserID},
				Expansions:  []gotwtr.Expansion{gotwtr.ExpansionReferencedTweetsIDAuthorID},}
			if max > 0 {ta.tlopt.UntilID = twid2str(max)}
			if since > 0 {ta.tlopt.SinceID = twid2str(since)}
		} else {
			ta.tlopt.PaginationToken = ta.usermeta.NextToken
		}
		if maxresult > 0 {ta.tlopt.MaxResults = maxresult}
		
		ta.seq++

		res, err := client.UserTweetTimeline(context.Background(), userID, &ta.tlopt)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d meta(%d)\n", time.Now().Format("15:04:05"), len(res.Tweets), res.Meta.ResultCount)
		tr := gotwtr.TweetsResponse{}
		tr.Tweets   = res.Tweets
		tr.Includes = res.Includes
		tr.Errors   = res.Errors
		tr.Title    = res.Title
		tr.Detail   = res.Detail
		tr.Type	    = res.Type
		ta.usermeta = res.Meta
		ta.usermeta = res.Meta
		return &tr, err
	case tlhome:
	case tlmention:
		if ta.seq == 0 {
			ta.tmopt = gotwtr.UserMentionTimelineOption {
				TweetFields: []gotwtr.TweetField{gotwtr.TweetFieldReferencedTweets,gotwtr.TweetFieldInReplyToUserID},
				Expansions:  []gotwtr.Expansion{gotwtr.ExpansionReferencedTweetsIDAuthorID},}
			if max > 0 {ta.tmopt.UntilID = twid2str(max)}
			if since > 0 {ta.tmopt.SinceID = twid2str(since)}
		} else {
			ta.tmopt.PaginationToken = ta.usermeta.NextToken
		}
		if maxresult > 0 {ta.tmopt.MaxResults = maxresult}
		
		ta.seq++

		res, err := client.UserMentionTimeline(context.Background(), userID, &ta.tmopt)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d meta(%d)\n", time.Now().Format("15:04:05"), len(res.Tweets), res.Meta.ResultCount)
		tr := gotwtr.TweetsResponse{}
		tr.Tweets   = res.Tweets
		tr.Includes = res.Includes
		tr.Errors   = res.Errors
		tr.Title    = res.Title
		tr.Detail   = res.Detail
		tr.Type	    = res.Type
		ta.usermeta = res.Meta
		return &tr, err
	case tlrtofme:
	case tllist:
	case tlsearch: fallthrough
	case tlsearcha:
		query := userID
		if ta.seq == 0 {
			ta.sropt = gotwtr.SearchTweetsOption {
				TweetFields: []gotwtr.TweetField{gotwtr.TweetFieldReferencedTweets,gotwtr.TweetFieldInReplyToUserID},
				Expansions:  []gotwtr.Expansion{gotwtr.ExpansionReferencedTweetsIDAuthorID},}
			if max > 0 {ta.sropt.UntilID = twid2str(max)}
			if since > 0 {ta.sropt.SinceID = twid2str(since)}
		} else {
			ta.sropt.NextToken = ta.srchmeta.NextToken
		}
		if maxresult > 0 {ta.sropt.MaxResults = maxresult}
		
		ta.seq++
		var res *gotwtr.SearchTweetsResponse
		switch ta.t {
		case tlsearch:
			res, err = client.SearchRecentTweets(context.Background(), query, &ta.sropt)
		case tlsearcha:
			res, err = client.SearchAllTweets(context.Background(), query, &ta.sropt)
		}
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d meta(%d)\n", time.Now().Format("15:04:05"), len(res.Tweets), res.Meta.ResultCount)
		tr := gotwtr.TweetsResponse{}
		tr.Tweets   = res.Tweets
		tr.Includes = res.Includes
		tr.Errors   = res.Errors
		tr.Title    = res.Title
		tr.Detail   = res.Detail
		tr.Type	    = res.Type
		ta.srchmeta = res.Meta
		return &tr, err
	}
	return nil, err
}

func (ta *twSearchApi) rewindQuery() {
	ta.seq = 0
}

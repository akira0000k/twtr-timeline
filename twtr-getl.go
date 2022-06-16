package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"flag"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/sivchari/gotwtr"
)

var exitcode int = 0

type twidt int64
var twid_def twidt = 0
var next_max twidt = twid_def
var next_since twidt = twid_def
func print_id() {
	if uniqid != nil {
		err := uniqid.write()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	fmt.Fprintf(os.Stderr, "--------\n-since_id=%d\n", next_since)
	fmt.Fprintf(os.Stderr,   "-max_id=%d\n", next_max)
}

const onetimedefault = 10
const onetimemax = 100
const onetimemin_t = 5
const onetimemin_s = 10
var onetimemin = onetimemin_t
const sleepdot = 5

// tweet/user hash for v2 api
type tweetHash map[string] *gotwtr.Tweet
type userHash map[string] *gotwtr.User

// TL type "enum"
type tltype int
const (
	tlnone tltype = iota
	tluser
	tlhome
	tlmention
	tlrtofme
	tllist
	tlsearch
	tlsearcha
)

type revtype bool
const (
	reverse revtype = true
	forward revtype = false
)

const (
	rsrecent string = "recent"
	rsall string = "all"
)

type idCheck map[string]bool
var uniqid idCheck = nil

func (c idCheck) checkID(id string) (exist bool) {
	if c[id] {
		return true
	} else {
		c[id] = true
		return false
	}
}

func (c idCheck) write() (err error) {
	bytes, _ := json.Marshal(c)
	err = ioutil.WriteFile("tempids.json", bytes, os.FileMode(0600))
	return err
}

func (c *idCheck) read() (err error) {
	*c = idCheck{}
	raw, err := ioutil.ReadFile("tempids.json")
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return nil
	}
	json.Unmarshal(raw, c)
	return nil
}

var twapi twSearchApi

func main(){
	var err error
	tLtypePtr := flag.String("get", "", "TLtype: user, mention, search")
	screennamePtr := flag.String("user", "", "twitter @ screenname")
	useridPtr := flag.String("userid", "0", "integer user Id")
	// listnamePtr := flag.String("listname", "", "list name")
	// listIDPtr := flag.Int64("listid", 0, "list ID")
	queryPtr := flag.String("query", "", "Query String")
	resulttypePtr := flag.String("restype", "", "result type: [recent]/all")
	countPtr := flag.Int("count", 0, "tweet count. 5-800 ?")
	eachPtr := flag.Int("each", 0, "req count for each loop 5-100")
	max_idPtr := flag.Int64("max_id", 0, "starting tweet id")
	since_idPtr := flag.Int64("since_id", 0, "reverse start tweet id")
	reversePtr := flag.Bool("reverse", false, "reverse output. wait newest TL")
	loopsPtr := flag.Int("loops", 0, "get loop max")
	waitPtr := flag.Int64("wait", 0, "wait second for next loop")
	// nortPtr := flag.Bool("nort", false, "not include retweets")
	flag.Parse()
	var tLtype = *tLtypePtr
	var screenname = *screennamePtr
	var userid = *useridPtr
	// var listname = *listnamePtr
	// var listID = *listIDPtr
	var queryString = *queryPtr
	var resulttype = *resulttypePtr
	var count = *countPtr
	var eachcount = *eachPtr
	var max_id = twidt(*max_idPtr)
	var since_id = twidt(*since_idPtr)
	var reverseflag = *reversePtr
	var max_loop = *loopsPtr
	var waitsecond = *waitPtr
	// var includeRTs = ! *nortPtr
	
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "positional argument no need [%s]\n", flag.Arg(0))
		os.Exit(2)
	}

	var t tltype
	switch tLtype {
	case "user":    t = tluser
	//case "home":    t = tlhome
	case "mention": t = tlmention
	//case "rtofme":  t = tlrtofme
	//case "list":    t = tllist
	case "search":  t = tlsearch
	case "":
		if userid != "0" || screenname != "" {
			t = tluser
			tLtype = "user"
			fmt.Fprintln(os.Stderr, "assume -get=user")
		} else if queryString != "" {
			t = tlsearch
			tLtype = "search"
			fmt.Fprintln(os.Stderr, "assume -get=search")
		} else {
			fmt.Fprintf(os.Stderr, "invalid type -get=%s\n", tLtype)
			os.Exit(2)
		}
	default:
		fmt.Fprintf(os.Stderr, "invalid type -get=%s\n", tLtype)
		os.Exit(2)
	}
	fmt.Fprintf(os.Stderr, "-get=%s\n", tLtype)
	
	twapi.client = connectTwitterApi()

	switch t {
	case tluser: fallthrough
	case tlmention:
		if userid != "0" {
			fmt.Fprintf(os.Stderr, "user id=%s\n", userid)
			if (screenname != "") {
				fmt.Fprintln(os.Stderr, "screen name ignored.")
			}
		} else if screenname != "" {
			userid, err = name2id(screenname)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(2)
			}
			fmt.Fprintf(os.Stderr, "user id=%s %s\n", userid, screenname)
		} else {
			fmt.Fprintf(os.Stderr, "no user id\n")
			os.Exit(2)
		}
	default:
		if userid != "0" || screenname != "" {
			fmt.Fprintf(os.Stderr, "-get=%s no need userid/screenname\n", tLtype)
			os.Exit(2)
		}
	}

	switch t {
	case tlsearch:
		if queryString == "" {
			fmt.Fprintln(os.Stderr, "-query not specified")
			os.Exit(2)
		}
		switch {
		case strings.HasPrefix(rsrecent, resulttype):
			//default
		case strings.HasPrefix(rsall, resulttype):
			t = tlsearcha
		default:
			fmt.Fprintf(os.Stderr, "invalid -restype=%s\n", resulttype)
			os.Exit(2)
		}
		onetimemin = onetimemin_s
		err = uniqid.read()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		userid = queryString
	default:
		if queryString != "" {
			fmt.Fprintf(os.Stderr, "-get=%s no need -query\n", tLtype)
			os.Exit(2)
		}
		if resulttype != "" {
			fmt.Fprintf(os.Stderr, "-get=%s no need -restype=%s\n", tLtype, resulttype)
			os.Exit(2)
		}
	}		
	twapi.t = t

	fmt.Fprintf(os.Stderr, "count=%d\n", count)
	fmt.Fprintf(os.Stderr, "each=%d\n", eachcount)
	fmt.Fprintf(os.Stderr, "reverse=%v\n", reverseflag)
	fmt.Fprintf(os.Stderr, "loops=%d\n", max_loop)
	fmt.Fprintf(os.Stderr, "max_id=%d\n", max_id)
	fmt.Fprintf(os.Stderr, "since_id=%d\n", since_id)
	fmt.Fprintf(os.Stderr, "wait=%d\n", waitsecond)

	sgchn.sighandle()
	
	if reverseflag {
		if max_id != 0 {
			fmt.Fprintf(os.Stderr, "max id ignored when reverse\n")
		}
		if waitsecond <= 0 {
			waitsecond = 60
			fmt.Fprintf(os.Stderr, "wait default=%d (reverse)\n", waitsecond)
		}
		getReverseTLs(userid, count, max_loop, waitsecond, since_id)
	} else {
		if max_loop == 0 && since_id == 0 && count == 0 {
			count = onetimemin
			fmt.Fprintf(os.Stderr, "set forward default count=%d\n", count)
		}
		if count != 0 && count < onetimemin {
			count = onetimemin
			fmt.Fprintf(os.Stderr, "set count=%d\n", count)
		}
		if eachcount != 0 && eachcount < onetimemin {
			eachcount = onetimemin
			fmt.Fprintf(os.Stderr, "set eachcount=%d\n", eachcount)
		} else if eachcount > onetimemax {
			eachcount = onetimemax
			fmt.Fprintf(os.Stderr, "set eachcount=%d\n", eachcount)
		}
			
		if max_id > 0 && max_id <= since_id {
			fmt.Fprintf(os.Stderr, "sincd id ignored when max<=since\n")
		}
		if waitsecond <= 0 {
			waitsecond = 5
			fmt.Fprintf(os.Stderr, "wait default=%d (forward)\n", waitsecond)
		}
		getFowardTLs(userid, count, eachcount, max_loop, waitsecond, max_id, since_id)
	}
	print_id()
	os.Exit(exitcode)
}

func getFowardTLs(userid string, count int, eachcount int, loops int, waitsecond int64, max twidt, since twidt) {
	var countlim bool = true
	if count <= 0 {
		countlim = false
	}
	if eachcount == 0 {
		if count > 0 {
			eachcount = count
			if eachcount > onetimemax {
				eachcount = onetimemax
			}
			fmt.Fprintf(os.Stderr, "-each=%d assumed\n", eachcount)
		}
	}
	if max > 0 {
		if max <= since {
			since = 0
		}
	}
	until := since
	if until > 0 {
		until -= 1
	}
	for i := 1; ; i++ {

		res, err := twapi.getTL(userid, eachcount, max, until)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			print_id()
			os.Exit(2)
		}
		// jsonTweets, _ := json.Marshal(tweets) //test
		// fmt.Println(string(jsonTweets))       //test
		
		c := len(res.Tweets)
		if c == 0 {
			exitcode = 1
			break
		}
		gotfirstid := str2twid(res.Tweets[0].ID)
		if next_max > 0 && gotfirstid >= next_max {
			fmt.Fprintln(os.Stderr, "same record. break")
			fmt.Fprintln(os.Stderr, "get[0]ID=", gotfirstid, " next_max=", next_max)
			break
		}

		var userhash = userHash{}
		var tweethash = tweetHash{}
		for _, u := range res.Includes.Users {
			//fmt.Println(u.ID, u.UserName, u.Name)
			userhash[u.ID] = u
		}
		for _, t := range res.Includes.Tweets {
			//fmt.Println(t.ID, userhash[t.AuthorID].UserName)
			tweethash[t.ID] = t
		}
		
		firstid, lastid, nout := printTL(res.Tweets, userhash, tweethash, count, forward)
		// fmt.Println("printTL id:", firstid, "-", lastid)
		if next_since == twid_def {
			next_since = firstid
		}
		next_max = lastid

		if lastid <= since {
			break  // break by since_id
		}
		if countlim {
			count -= nout
			if count <= 0 { break }
		}
		if loops > 0 && i >= loops {
			break
		}

		sleep(waitsecond) //?
	}
	return
}

func getReverseTLs(userid string, count int, loops int, waitsecond int64, since twidt) {
	var tweets []*gotwtr.Tweet
	var userhash = userHash{}
	var tweethash = tweetHash{}
	var zeror bool
	var countlim bool = true
	if count <=  0 {
		countlim = false
	}
	var sinceid = since
	var zerocount int = 0
	const maxzero int = 1
	next_since = sinceid //default: same sinceid
	if sinceid <= 0 {
		fmt.Fprintf(os.Stderr, "since=%d. get %d tweet\n", sinceid, onetimemin)
		res, err := twapi.getTL(userid, onetimemin, 0, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			print_id()
			os.Exit(2)
		}
		c := len(res.Tweets)
		if c == 0 {
			fmt.Fprintln(os.Stderr, "Not 1 record available")
			os.Exit(2)
		}

		for _, u := range res.Includes.Users {
			//fmt.Println(u.ID, u.UserName, u.Name)
			userhash[u.ID] = u
		}
		for _, t := range res.Includes.Tweets {
			//fmt.Println(t.ID, userhash[t.AuthorID].UserName)
			tweethash[t.ID] = t
		}
 
		firstid, lastid, _ := printTL(res.Tweets, userhash, tweethash, 0, reverse)
		next_max = firstid
		next_since = lastid
		sinceid = lastid
		sleep(5)
	} else {
		fmt.Fprintf(os.Stderr, "since=%d. start from this record.\n", sinceid)
	}
	for i:=1; ; i+=1 {
		tweets, userhash, tweethash, zeror = getTLsince(userid, sinceid)
 
		c := len(tweets)
		if c > 0 {
			zerocount = 0
			minid := str2twid(tweets[len(tweets) - 1].ID)
			if minid <= sinceid {
				//指定ツイートまで取れたのでダブらないように削除する
				tweets = tweets[: len(tweets) - 1]
				c = len(tweets)
			} else {
				//gap
				fmt.Fprintf(os.Stderr, "Gap exists\n")
			}
			if c > 0 {
				firstid, lastid, nout := printTL(tweets, userhash, tweethash, 0, reverse)
				if next_max == 0 {
					next_max = firstid
				}
				next_since = lastid
				sinceid = lastid
				if countlim {
					count -= nout
					if count <= 0 { break }
				}
			}
			if zeror {
				//accident. no response
				zerocount += 1
				if zerocount == maxzero {
					exitcode = 1
					break
				}
			}
		} else {
			//accident. no response
			zerocount += 1
			if zerocount == maxzero {
				exitcode = 1
				break
			}
		}
		if loops > 0 && i >= loops {
			break
		}
		sleep(waitsecond)
	}
	return
}
 
func getTLsince(userid string, since twidt) (tweets []*gotwtr.Tweet, userhash userHash, tweethash tweetHash, zeror bool) {
	tweets = []*gotwtr.Tweet{}
	userhash = userHash{}
	tweethash = tweetHash{}
	zeror = false
	var max_id twidt = 0
	twapi.rewindQuery()
	for i := 0; ; i++ {
 
		res, err := twapi.getTL(userid, onetimemax, max_id, since - 1)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			print_id()
			os.Exit(2)
		}
		c := len(res.Tweets)
		if c == 0 {
			zeror = true
			break
		}
 
		lastid := str2twid(res.Tweets[c - 1].ID)
 
		tweets = append(tweets, res.Tweets...)

		for _, u := range res.Includes.Users {
			//fmt.Println(u.ID, u.UserName, u.Name)
			userhash[u.ID] = u
		}
		for _, t := range res.Includes.Tweets {
			//fmt.Println(t.ID, userhash[t.AuthorID].UserName)
			tweethash[t.ID] = t
		}
 
		if lastid <= since {
			break
		}
		// 一度で取りきれなかった
		fmt.Fprintln(os.Stderr, "------continue")
		max_id = lastid - 1
 
		sleep(10) //??
	}
	return tweets, userhash, tweethash, zeror
}

func str2twid(tweetID string) (twidt) {
	id64, _  := strconv.ParseInt(tweetID, 10, 64)
	return twidt(id64)
}


func printTL(tweets []*gotwtr.Tweet, userhash userHash, tweethash tweetHash, count int, revs revtype) (firstid twidt, lastid twidt, nout int) {

	firstid = twid_def
	lastid = twid_def
	imax := len(tweets)
	is := 0
	ip := 1
	if revs {
		is = imax - 1
		ip = -1
	}
	nout = 0
	for i := is; 0 <= i && i < imax; i += ip {
		tweet := tweets[i]
		id := str2twid(tweet.ID)
		if i == is {
			firstid = id
			lastid = id
		}
		twtype, rt := ifRetweeted(tweet, tweethash)
		//  RT > Reply > Mention > tweet
		var done bool
		if rt != nil {
			twtype2, _ := ifRetweeted(rt, tweethash)
			done = printTweet(twtype, tweet, twtype2, rt, userhash)
		} else {
			done = printTweet("or", tweet, twtype, tweet, userhash)
		}
		if done {
			nout++
		}

		lastid = id
		
		if count > 0 && nout >= count {
			break
		}
	}
	return firstid, lastid, nout
}

func ifRetweeted(t *gotwtr.Tweet, rhash tweetHash) (twtype string, rt *gotwtr.Tweet) {
	twtype = "tw"
	if t.InReplyToUserID != "" {
		twtype = "Mn"
	}
	rt = nil
	var ok bool
	for _, r := range t.ReferencedTweets {
		switch r.Type {
		case "retweeted":
			rt, ok = rhash[r.ID]
			if ok {
				twtype = "RT"
			}
		case "replied_to":
			twtype = "Re"
		case "quoted":
		}
	}
	return twtype, rt
}

func printTweet(twtype1 string, tweet1 *gotwtr.Tweet, twtype2 string, tweet2 *gotwtr.Tweet, userhash userHash) (bool) {
	tweetid := tweet1.ID
	tweetuser := userhash[tweet1.AuthorID].UserName

	origiid := tweet2.ID
	origiuser := userhash[tweet2.AuthorID].UserName
	
	firstp := true
	idst := "*Id:"
	if uniqid != nil {
		if uniqid.checkID(origiid) {
			firstp = false
			idst = "_id:"
		}
	}

	if tweetid == origiid {
		fmt.Fprintln(os.Stderr, idst, tweetid)
	} else {
		fmt.Fprintln(os.Stderr, idst, tweetid, origiid)
	}

	if firstp {
		fmt.Printf("%s\t@%s\t%s\t%s\t@%s\t%s\t\"%s\"\n",
			tweetid, tweetuser, twtype1,
			origiid, origiuser, twtype2, quoteText(tweet2.Text))
	}
	return firstp
}

func quoteText(fulltext string) (qtext string) {
	quoted1 := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fulltext, "\n", `\n`), "\r", `\r`), "\"", `\"`)
	qtext = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(quoted1, `&amp;`, `&`), `&lt;`, `<`), `&gt;`, `>`)
	return qtext
}


func name2id(screen_name string) (id string, err error) {
	ures, err := twapi.client.RetrieveSingleUserWithUserName(context.Background(), screen_name)
	if err != nil {
		return "", err
	}
	//jsonUser, _ := json.Marshal(users[0])
	//fmt.Println(string(jsonUser))
	//os.Exit(9)
	
	var userinfo *gotwtr.User = ures.User
 
	id = userinfo.ID
	return id, nil
}

func connectTwitterApi() (client *gotwtr.Client) {
	usr, _ := user.Current()
	raw, error := ioutil.ReadFile(usr.HomeDir + "/twitter/twitterBearerToken.json")
	if error != nil {
		fmt.Fprintln(os.Stderr, error.Error())
		os.Exit(2)
	}
	var twitterBearerToken TwitterBearerToken
	json.Unmarshal(raw, &twitterBearerToken)

	// raw, error = ioutil.ReadFile(usr.HomeDir + "/twitter/twitterAccount.json")
	// if error != nil {
	//  	fmt.Fprintln(os.Stderr, error.Error())
	//  	os.Exit(2)
	// }
	// var twitterAccount TwitterAccount
	// json.Unmarshal(raw, &twitterAccount)

	client =  gotwtr.New(twitterBearerToken.BearerToken)
	// client =  gotwtr.New(twitterBarerToken.BarerToken,
	//  	gotwtr.WithConsumerKey(twitterAccount.ConsumerKey),
	//  	gotwtr.WithConsumerSecret(twitterAccount.ConsumerSecret))
	return client
}

type TwitterAccount struct {
	AccessToken string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	ConsumerKey string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
}

type TwitterBearerToken struct {
	BearerToken string `json:"bearerToken"`
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

// type UserMentionTimelineResponse struct {
//  	Tweets   []*Tweet            `json:"data"`
//  	Includes *TweetIncludes      `json:"includes,omitempty"`
//  	Errors   []*APIResponseError `json:"errors,omitempty"`
//  	Meta     *UserTimelineMeta   `json:"meta"`
//  	Title    string              `json:"title,omitempty"`
//  	Detail   string              `json:"detail,omitempty"`
//  	Type     string              `json:"type,omitempty"`
// }

// type TweetsResponse struct {
//  	Tweets   []*Tweet            `json:"data"`
//  	Includes *TweetIncludes      `json:"includes,omitempty"`
//  	Errors   []*APIResponseError `json:"errors,omitempty"`
//  	Title    string              `json:"title,omitempty"`
//  	Detail   string              `json:"detail,omitempty"`
//  	Type     string              `json:"type,omitempty"`
// }

// type TweetIncludes struct {
//  	Media  []*Media
//  	Places []*Place
//  	Polls  []*Poll
//  	Tweets []*Tweet
//  	Users  []*User
// }

// type UserTimelineMeta struct {
//  	ResultCount int    `json:"result_count"`
//  	NewestID    string `json:"newest_id"`
//  	OldestID    string `json:"oldest_id"`
//  	NextToken   string `json:"next_token"`
// }

// type UserTweetTimelineOption struct {
//  	EndTime         time.Time
//  	Exclude         []Exclude
//  	Expansions      []Expansion
//  	MaxResults      int
//  	MediaFields     []MediaField
//  	PaginationToken string
//  	PlaceFields     []PlaceField
//  	PollFields      []PollField
//  	SinceID         string
//  	StartTime       time.Time
//  	TweetFields     []TweetField
//  	UntilID         string
//  	UserFields      []UserField
// }

// type UserMentionTimelineOption struct {
//  	EndTime         time.Time
//  	Expansions      []Expansion
//  	MaxResults      int
//  	MediaFields     []MediaField
//  	PaginationToken string
//  	PlaceFields     []PlaceField
//  	PollFields      []PollField
//  	SinceID         string
//  	StartTime       time.Time
//  	TweetFields     []TweetField
//  	UntilID         string
//  	UserFields      []UserField
// }

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

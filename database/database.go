package database

import (
	"fmt"
	"log"
	"regexp"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

const hostURL = "mongodb://159.89.28.157:27017"
const databaseName = "social_net"
const collName = "tweets"

var db *mgo.Database
var coll *mgo.Collection
var session *mgo.Session

//Tweet is the main structure for a Tweet in this database
type Tweet struct {
	ID       int    `bson:"id"`
	Polarity int    `bson:"polarity"`
	Date     string `bson:"date"`
	Query    string `bson:"query"`
	User     string `bson:"user"`
	Text     string `bson:"text"`
}

func init() {
	var err error
	session, err = GetSession()
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(databaseName)
	coll = db.C(collName)

}

//GetSession returns the mgo MongoDB session.
func GetSession() (*mgo.Session, error) {
	var err error
	session, err = mgo.Dial(hostURL)
	if err != nil {
		log.Fatal(err)
	}

	return session, err
}

//GetUserCount - this returns the amount of users in the database.
func GetUserCount() int {
	var result []string
	err := coll.Find(nil).Distinct("user", &result)
	if err != nil {
		log.Fatal(err)
	}

	return len(result)

}

//GetTopTaggers - this methods returns the top 10 taggers.
// I use regex to match all tweets which start with @.
// When it is done, it group the matches together by unique user with a sum value for each user.
// This is done each time a user matches and has tagged another user.
// At last it sorts and limits the query to 5 results.
func GetTopTaggers() {
	var result []bson.M

	pipeline := []bson.M{
		{"$match": bson.M{"text": bson.M{"$regex": bson.RegEx{`@\w`, ""}}}},
		{"$group": bson.M{"_id": "$user",
			"matches": bson.M{"$sum": 1},
		},
		},
		{"$sort": bson.M{"matches": -1}}, //1: Ascending, -1: Descending
		{"$limit": 10},
	}

	err := coll.Pipe(pipeline).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	for i, user := range result {
		fmt.Println(i+1, user["_id"], "has tagged others:", user["matches"], "times")

	}
}

//GetMostTagged - this method returns top 5 of most tagged users
func GetMostTagged() {
	var result []bson.M

	pipeline := []bson.M{
		{"$match": bson.M{"text": bson.M{"$regex": bson.RegEx{`@\w+`, ""}}}},
		{"$group": bson.M{"_id": "$text",
			"matches": bson.M{"$sum": 1},
		},
		},
		{"$sort": bson.M{"matches": -1}}, //1: Ascending, -1: Descending
		{"$limit": 5},
	}

	regEx, _ := regexp.Compile(`@\w+`)

	err := coll.Pipe(pipeline).AllowDiskUse().All(&result)
	if err != nil {
		log.Fatal(err)
	}

	for i, user := range result {
		fmt.Println(i+1, regEx.FindString(user["_id"].(string)))

	}

}

//GetMostActive - this method returns the 10 most active twitter users
func GetMostActive() {
	var result []bson.M

	pipeline := []bson.M{
		{"$match": bson.M{"user": bson.M{"$regex": bson.RegEx{`.*`, ""}}}}, //This RegEx grabs everything in the 'user'-field.
		{"$group": bson.M{"_id": "$user",
			"matches": bson.M{"$sum": 1},
		},
		},
		{"$sort": bson.M{"matches": -1}}, //1: Ascending, -1: Descending
		{"$limit": 10},
	}

	err := coll.Pipe(pipeline).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	for i, user := range result {
		fmt.Println(i+1, user["_id"], "has made:", user["matches"], "tweets")

	}
}

//GetGrumpiest - this methos returns the fice most grumpiest user tweets. Based on some negative words.
func GetGrumpiest() {
	var result []bson.M

	const negativeWords = "(shit|fuck|damn|bitch|crap|piss|dick|darn|asshole|bastard|douche|sad|angry|stupid)"

	pipeline := []bson.M{
		{"$match": bson.M{"text": bson.M{"$regex": bson.RegEx{negativeWords, ""}}}},
		{"$group": bson.M{"_id": "$user",
			"matches": bson.M{"$sum": 1},
		},
		},
		{"$sort": bson.M{"matches": -1}}, //1: Ascending, -1: Descending
		{"$limit": 5},
	}

	err := coll.Pipe(pipeline).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Most negative tweeters:")
	for i, user := range result {
		fmt.Println(i+1, user["_id"], "-", user["matches"], "negative tweets.")
	}

}

//GetHappiest - this methos returns the fice most happiest user tweets. Based on some positive words.
func GetHappiest() {
	var result []bson.M

	const positiveWords = "(love|happy|amazing|beautiful|yay|joy|pleasure|smile|win|winning|smiling|healthy|delight|paradise|positive|fantastic|blessed|splendid|sweetheart|great|funny)"

	pipeline := []bson.M{
		{"$match": bson.M{"text": bson.M{"$regex": bson.RegEx{positiveWords, ""}}}},
		{"$group": bson.M{"_id": "$user",
			"matches": bson.M{"$sum": 1},
		},
		},
		{"$sort": bson.M{"matches": -1}}, //1: Ascending, -1: Descending
		{"$limit": 5},
	}

	err := coll.Pipe(pipeline).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Most positive tweeters:")
	for i, user := range result {
		fmt.Println(i+1, user["_id"], "-", user["matches"], "positive tweets.")
	}
}

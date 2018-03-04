package main

import (
	"fmt"
	"log"

	database "./database"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	session, err := database.GetSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	//Switching the session to monotonic behavior. Not necessary.
	session.SetMode(mgo.Monotonic, true)

	//Uncomment the method you want to run. Remember to build the project first (described in the README).

	//1. Count Users
	// fmt.Println("Counting users...")
	// fmt.Println(database.GetUserCount())

	// //2. Top Taggers
	// fmt.Println("Counting top taggers...")
	// database.GetTopTaggers()

	// //3. Most Tagged
	// fmt.Println("Counting most tagged users...")
	// database.GetMostTagged()

	// //4. Most Active
	// fmt.Println("Counting most active users...")
	// database.GetMostActive()

	// //5 Most Grumpiest
	// fmt.Println("Getting grumpiest users...")
	// database.GetGrumpiest()

	// //6. Most Happiest
	fmt.Println("Getting happiest users...")
	database.GetHappiest()

}

func printUsage() {
	fmt.Println("usage:\n--CountUsers\n--TopTaggers\n--MostTagged\n--MostActive\n--Grumpiest\n--Happiest")
}

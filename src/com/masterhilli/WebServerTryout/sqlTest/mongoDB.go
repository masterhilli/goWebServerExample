package main
import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	//startMySQL()
	StartMongoDB()
}

type Person struct {
	Name string
	Phone string
}
func StartMongoDB() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}

	index := mgo.Index{
		Key: []string{"name"},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("people")
	c.EnsureIndex(index)
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})

	if (mgo.IsDup(err)) {
		fmt.Println("Ale + Cla is duplicate!" + err.Error())
	}
	err = c.Insert(&Person{"Martin", "Yep that is my number ;) "})
	if (mgo.IsDup(err)) {
		fmt.Println("Martin is duplicate!" + err.Error())
	}
	err = c.Insert(&Person{"Martin", "Yep thats is my number ;) "})
	if (mgo.IsDup(err)) {
		fmt.Println("Martin2 is duplicate!" + err.Error())
	}

	err = c.Insert(&Person{"Martin22", "Yep thats is my number ;) "})
	if (mgo.IsDup(err)) {
		fmt.Println("Martin22 is duplicate!" + err.Error())
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Martin"}).One(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println("Phone:", result.Phone)

	var results []Person
	err = c.Find(nil).All(&results)
	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)
}
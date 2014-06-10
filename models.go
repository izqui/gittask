package main

import (
	"fmt"

	"github.com/izqui/helpers"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Id bson.ObjectId `json:"_id" bson:"_id"`

	Email       string `json:"email,omitempty" bson:"email"`
	Username    string `json:"login" bson:"username"`
	Location    string `json:"location,omitempty" bson:"location"`
	AccessToken string `json:"" bson:"token" map:"token"`
}

func (u *User) Update() {

	userCollection := DB.C("users")

	change := mgo.Change{

		ReturnNew: true,
		Update:    helpers.StructToBSONMap(u),
	}
	fmt.Println(change.Update)
	if _, err := userCollection.Find(bson.M{"_id": u.Id}).Apply(change, u); err != nil {

		panic(err)
	}

	return
}

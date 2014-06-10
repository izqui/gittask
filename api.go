package main

import (
	"fmt"
	//"github.com/izqui/helpers"
	"github.com/izqui/oauth2"
	//"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Api struct{}

func (a *Api) Login(tokens oauth2.Tokens) {

	userCollection := DB.C("users")
	user := &User{}
	token := tokens.Access()

	//Is a user with that token in the database?
	if err := userCollection.Find(bson.M{"token": token}).One(&user); err != nil || user == nil {

		//Get user and save it if not in database
		github := &Github{AccessToken: token}
		user = github.UserInfo("me")

		//Check if that github user is already in the database with a different token
		dbuser := &User{}
		if err := userCollection.Find(bson.M{"username": user.Username}).One(&dbuser); err != nil || dbuser == nil {

			user.Id = bson.NewObjectId()
			//Save user
			if err := userCollection.Insert(user); err != nil {

				panic(err)
			}
		} else {

			dbuser.AccessToken = token
			dbuser.Update()
		}

	} else {

		fmt.Printf("Db user %v", *user)

	}
}

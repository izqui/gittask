package main

import (
	//"fmt"

	"github.com/izqui/oauth2"

	//"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Api struct{}

func (a *Api) Login(tokens oauth2.Tokens) {

	userCollection := DB.C("users")
	user := &User{}

	//Is a user with that token in the database?

	if err := userCollection.Find(bson.M{"token": tokens.Access()}).One(&user); err != nil {

		panic(err)
	}
	if user == nil {

		//Get user and save it if not in database
		github := Github{AccessToken: tokens.Access()}
		user = github.UserInfo("me")
		user.Id = bson.NewObjectId()

		if err := userCollection.Insert(user); err != nil {

			panic(err)
		}

	}
}

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
	AccessToken string `json:"" bson:"token"`
}

type Repo struct {
	Id     bson.ObjectId `bson:"_id"`
	UserId bson.ObjectId `bson:"user_id"`

	GithubId int    `json:"id,omitempty" bson:"github_id"`
	Name     string `json:"name,omitempty" bson:"name"`
	FullName string `json:"full_name,omitempty" bson:"full_name"`
	Language string `json:"language,omitempty" bson:"language"`

	Tasks []Task `bson:"tasks"`
}

type Task struct {
	Name     string `bson:"name"`
	Priority int    `bson:"priority"`
	Status   string `bson:"status"`
}

func CurrentUser(token string) (user *User) {

	userCollection := DB.C("users")

	if err := userCollection.Find(bson.M{"token": token}).One(&user); err != nil {

		panic(err)

	} else {

		return user
	}
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

func (u *User) GetRepos() (err error, repos []Repo) {

	repoCollection := DB.C("repos")

	err = repoCollection.Find(bson.M{"user_id": u.Id}).All(&repos)
	return
}

func (u *User) NewRepo(name string) (err error, repo *Repo) {

	repoCollection := DB.C("repos")

	//I probably should implement a better error handler
	err = repoCollection.Find(bson.M{"user_id": u.Id, "name": name}).One(&repo)

	if repo == nil {

		g := &Github{AccessToken: u.AccessToken}
		err, repo = g.GetRepo(u.Username, name)

		//Saving Repo to DB
		repo.Id = bson.NewObjectId()
		repo.UserId = u.Id
		err = repoCollection.Insert(repo)
	}

	return
}

func (r *Repo) Update() {

	repoCollection := DB.C("repos")

	change := mgo.Change{

		ReturnNew: true,
		Update:    helpers.StructToBSONMap(r),
	}
	fmt.Println(change.Update)
	if _, err := repoCollection.Find(bson.M{"_id": r.Id}).Apply(change, r); err != nil {

		panic(err)
	}

	return
}

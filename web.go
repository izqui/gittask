package main

import (
	"fmt"
	"net/http"

	"github.com/izqui/helpers"
	"github.com/izqui/oauth2"

	"github.com/codegangsta/martini-contrib/render"
	"labix.org/v2/mgo/bson"
)

type Website struct{}

func (w *Website) Index(tokens oauth2.Tokens) string {

	//It is needed to check if a value is nil because of martini's dependency injector
	if !helpers.IsNil(tokens) {

		user := CurrentUser(tokens.Access())

		return fmt.Sprintf("You're logged in %v :)", user.Username)
	}
	return "You're not logged in :("
}

func (w *Website) RepoList(tokens oauth2.Tokens) string {

	user := CurrentUser(tokens.Access())

	return fmt.Sprintf("%s repositories", user.Username)
}

func (w *Website) NewRepoGet(tokens oauth2.Tokens, r render.Render) {

	user := CurrentUser(tokens.Access())

	g := &Github{AccessToken: tokens.Access()}
	repos := g.UserRepos("me")

	data := struct {
		User  *User
		Repos []Repo
	}{
		user,
		repos,
	}

	r.HTML(200, "new_repo", data)
}

func (w *Website) NewRepoPost(tokens oauth2.Tokens, request *http.Request) string {

	repoCollection := DB.C("repos")

	project := request.FormValue("project")

	user := CurrentUser(tokens.Access())
	repo := &Repo{}

	//I probably should implement a better error handler
	if err := repoCollection.Find(bson.M{}).One(&repo); repo == nil || err != nil {

		g := &Github{AccessToken: tokens.Access()}
		repo := g.GetRepo(user.Username, project)

		//Saving Repo to DB
		repo.Id = bson.NewObjectId()
		repo.UserId = user.Id
		if err := repoCollection.Insert(repo); err != nil {

			panic(err)

		} else {

			return "Saved repo to DB"
		}

	} else {

		return "Repo already exists in database"
	}

	return "WTF"
}

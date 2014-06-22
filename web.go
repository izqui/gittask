package main

import (
	"fmt"
	"net/http"

	"github.com/izqui/helpers"
	"github.com/izqui/oauth2"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"labix.org/v2/mgo/bson"
)

type Website struct{}

func (w *Website) Index(tokens oauth2.Tokens) string {

	//It is needed to check if a value is nil because of martini's dependency injector
	if !helpers.IsNil(tokens) {

		//user := CurrentUser(tokens.Access())
		return fmt.Sprintf("You're logged")
	}
	return "You're not logged in :("
}

func (w *Website) RepoList(tokens oauth2.Tokens, r render.Render) {

	user := CurrentUser(tokens.Access())

	repoCollection := DB.C("repos")

	var repos []Repo
	if err := repoCollection.Find(bson.M{"user_id": user.Id}).All(&repos); err != nil {

		panic(err)
	} else {

		data := struct {
			User  *User
			Repos []Repo
		}{
			user,
			repos,
		}

		r.HTML(200, "repos", data)
	}

}

func (w *Website) NewRepoGet(tokens oauth2.Tokens, r render.Render) {

	user := CurrentUser(tokens.Access())

	g := &Github{AccessToken: tokens.Access()}
	repos := g.UserRepos("me")

	/*
		//Don't present repos already in the db
		repoCollection := DB.C("repos")

		var dbrepos []Repo
		if err := repoCollection.Find(bson.M{"user_id": user.Id}).All(&dbrepos); dbrepos != nil { //|| err != nil {

			for _, r := range dbrepos {
				for _, rs := range repos {
					if r == rs {

					}
				}
			}
		}
	*/
	data := struct {
		User  *User
		Repos []Repo
	}{
		user,
		repos,
	}

	r.HTML(200, "new_repo", data)
}

func (w *Website) NewRepoPost(tokens oauth2.Tokens, request *http.Request, r render.Render) string {

	repoCollection := DB.C("repos")

	project := request.FormValue("project")

	user := CurrentUser(tokens.Access())
	repo := &Repo{}

	//I probably should implement a better error handler
	if err := repoCollection.Find(bson.M{"user_id": user.Id, "name": project}).One(&repo); repo == nil || err != nil {

		g := &Github{AccessToken: tokens.Access()}
		repo := g.GetRepo(user.Username, project)

		//Saving Repo to DB
		repo.Id = bson.NewObjectId()
		repo.UserId = user.Id
		if err := repoCollection.Insert(repo); err != nil {

			panic(err)

		} else {
			r.Redirect("/repos")
		}

	} else {

		return "Repo already exists in database"
	}

	return "WTF"
}

func (w *Website) GetRepo(params martini.Params, tokens oauth2.Tokens, r render.Render) {

	username := params["user"]
	reponame := params["repo"]
	fullname := fmt.Sprintf("%s/%s", username, reponame)

	repoCollection := DB.C("repos")
	repo := new(Repo)
	if err := repoCollection.Find(bson.M{"full_name": fullname}).One(&repo); err != nil || repo == nil {

		r.HTML(404, "404", nil)

	} else {

		isowner := false
		var user *User
		//If there is an authenticated user, find out whether he is the owner of the project
		if !helpers.IsNil(tokens) {

			user = CurrentUser(tokens.Access())
			isowner = (user.Id == repo.UserId)
		}

		data := struct {
			Owner bool
			User  *User
			Repo  *Repo
		}{
			isowner,
			user,
			repo,
		}

		r.HTML(200, "repo", data)

	}
}

func (w *Website) NewTask(params martini.Params, tokens oauth2.Tokens, r render.Render, request *http.Request) {

	username := params["user"]
	reponame := params["repo"]
	fullname := fmt.Sprintf("%s/%s", username, reponame)
	taskname := request.FormValue("task")

	repoCollection := DB.C("repos")
	repo := new(Repo)
	if err := repoCollection.Find(bson.M{"full_name": fullname}).One(&repo); err != nil || repo == nil {

		r.HTML(404, "404", nil)

	} else {

		task := Task{Name: taskname}

		repo.Tasks = append(repo.Tasks, task)
		repo.Update()

		r.Redirect("/repo/" + fullname)
	}
}

package main

import (
	"fmt"
	"net/http"

	"github.com/izqui/helpers"
	"github.com/izqui/oauth2"

	"github.com/codegangsta/martini-contrib/render"
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

	project := request.FormValue("project")
	return project
}

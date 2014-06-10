package main

import (
	"fmt"

	"github.com/izqui/oauth2"

	"github.com/izqui/helpers"
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

	g := &Github{AccessToken: tokens.Access()}
	repos := g.UserRepos("me")

	fmt.Println(repos)

	return fmt.Sprintf("%s repositories", user.Username)
}

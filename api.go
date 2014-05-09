package main

import (
	"fmt"
	//"github.com/go-martini/martini"
	"github.com/izqui/oauth2"
)

type Api struct{}

func (a *Api) Login(tokens oauth2.Tokens) {

	github := Github{AccessToken: tokens.Access()}

	user := github.UserInfo("me")
	fmt.Printf("Github User Authenticated: %v", user)

}
func (a *Api) Logout(tokens oauth2.Tokens) {

	fmt.Printf("Log out callback")
}

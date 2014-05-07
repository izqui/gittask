package main

import (
	"fmt"

	//"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
)

type Website struct{}

func (w *Website) Index(tokens oauth2.Tokens) string {

	return fmt.Sprintf("Hello %s", tokens.Access())
}

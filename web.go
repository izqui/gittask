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

		return fmt.Sprintf("You're logged in %v :)", tokens.Access())
	}
	return "You're not logged in :("
}

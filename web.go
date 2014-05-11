package main

import (
	"fmt"
	"reflect"
	//"github.com/go-martini/martini"
	"github.com/izqui/oauth2"
)

type Website struct{}

func (w *Website) Index(tokens oauth2.Tokens) string {

	if !isNil(tokens) {

		return fmt.Sprintf("You're logged in %v :)", tokens.Access())
	}
	return "You're not logged in :("
}

func isNil(v interface{}) bool {

	//It is needed to check if a value is nil because of martini's dependency injector
	return reflect.ValueOf(v).IsNil()
}

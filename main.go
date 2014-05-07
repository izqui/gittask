package main

import (
	"flag"
	//"fmt"
	"os"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
)

var host string = "http://localhost:"
var server *martini.Martini

func init() {

	setupServer(server)
}

func main() {

	server.Run()
}

func setupServer(s *martini.Martini) {

	//PORT
	port := flag.String("port", "9000", "Port for the HTTP server")
	flag.Parse()
	os.Setenv("PORT", *port)
	host = host + *port

	//MARTINI
	s = martini.New()
	s.Use(martini.Logger())
	s.Use(martini.Recovery())
	s.Use(sessions.Sessions("sessionbro", sessions.NewCookieStore([]byte("olakase"))))
	s.Use(oauth2.Github(&oauth2.Options{
		
		//From keys_.go file
		ClientId:     ClientId,
		ClientSecret: ClientSecret,
		RedirectURL:  host + "/oauth2callback",
		Scopes:       []string{"user", "repo"},
	}))
	//s.Use(martini.Static("static"))

	router := martini.NewRouter()
	website := &Website{}

	router.Get("/", oauth2.LoginRequired, website.Index)

	s.Action(router.Handle)

	server = s

}

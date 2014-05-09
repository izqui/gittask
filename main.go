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
var server *martini.ClassicMartini

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
	s = martini.Classic()
	s.Use(sessions.Sessions("sessionbro", sessions.NewCookieStore([]byte("olakase"))))
	s.Use(oauth2.Github(&oauth2.Options{
		
		//From keys_.go file
		ClientId:     ClientId,
		ClientSecret: ClientSecret,
		RedirectURL:  host + "/oauth2callback",
		Scopes:       []string{"user", "repo"},
	}))
	//s.Use(martini.Static("static"))

	website := &Website{}

	s.Get("/", oauth2.LoginRequired, website.Index)

	server = s

}

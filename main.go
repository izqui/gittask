package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-martini/martini"
	"github.com/izqui/oauth2"
	"github.com/martini-contrib/sessions"
)

var server *martini.Martini

func init() {

	setupServer(server)
}

func main() {

	server.Run()
}

func setupServer(s *martini.Martini) {

	//MODULES
	website := &Website{}
	api := &Api{}

	//PORT
	port := flag.String("port", "9000", "Port for the HTTP server")
	flag.Parse()
	os.Setenv("PORT", *port)

	baseurl := fmt.Sprintf("%s://%s:%s", GittaskProtocol, GittaskHost, *port)

	//MARTINI
	s = martini.New()
	s.Use(martini.Logger())
	s.Use(martini.Recovery())
	s.Use(sessions.Sessions("sessionbro", sessions.NewCookieStore([]byte("olakase"))))
	s.Use(oauth2.Github(&oauth2.Options{

		//From conf_.go file
		ClientId:      GithubClientId,
		ClientSecret:  GithubClientSecret,
		RedirectURL:   baseurl + "/oauth2callback",
		Scopes:        []string{"user", "repo"},
		LoginCallback: api.Login,
	}))
	//s.Use(martini.Static("static"))

	router := martini.NewRouter()

	router.Get("/", oauth2.LoginRequired, website.Index)

	s.Action(router.Handle)

	server = s

}

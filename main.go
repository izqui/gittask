package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-martini/martini"
	"github.com/izqui/oauth2"
	"github.com/martini-contrib/sessions"

	"labix.org/v2/mgo"
)

var server *martini.Martini = martini.New()
var DB *mgo.Database = &mgo.Database{}

func init() {

	setupDB(DB)
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
	host := flag.String("host", GittaskHost, "Host")

	flag.Parse()
	os.Setenv("PORT", *port)

	baseurl := fmt.Sprintf("%s://%s:%s", GittaskProtocol, *host, *port)

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
	router.Get("/repos", oauth2.LoginRequired, website.RepoList)

	s.Action(router.Handle)

	server = s
}

func setupDB(db *mgo.Database) {

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{DatabaseHost},
		Username: DatabaseUser,
		Password: DatabasePasswd,
		Database: DatabaseName,
	})

	if err != nil {
		panic(err)
	}

	db = session.DB(DatabaseName)
	DB = db

}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"io"
	"io/ioutil"

	"net/http"
	"net/url"
	"os"
	"strconv"

	//"github.com/izqui/helpers"
)

const (
	baseurl = "https://api.github.com"
)

type Github struct {
	AccessToken string
}

type jsonCallback chan []byte

func (g *Github) UserInfo(username string) *User {

	callback := make(jsonCallback)

	path := "/user"

	if username != "" && username != "me" {

		path = "/users/" + username
	}

	go g.request("GET", path, nil, nil, callback)

	select {

	case body := <-callback:

		user := new(User)
		json.Unmarshal(body, &user)

		return user
	}

}

func (g *Github) UserRepos(username string) []Repo {

	callback := make(jsonCallback)

	path := "/user/repos"

	if username != "" && username != "me" {

		path = "/users/" + username + "/repos"
	}

	go g.request("GET", path, map[string]string{"type": "public", "sort": "updated"}, nil, callback)

	select {

	case body := <-callback:

		var repos []Repo
		if err := json.Unmarshal(body, &repos); err != nil {

			panic(err)
		}
		return repos
	}

}

func (g *Github) GetRepo(owner string, repo string) *Repo {

	cb := make(jsonCallback)
	path := fmt.Sprintf("/repos/%s/%s", owner, repo)

	go g.request("GET", path, nil, nil, cb)

	select {

	case body := <-cb:

		repo := new(Repo)
		if err := json.Unmarshal(body, &repo); err != nil {
			panic(err)
		}

		return repo
	}
}
func (g *Github) request(method string, path string, params map[string]string, body map[string]string, cb jsonCallback) {

	//Append URL params
	if params != nil {

		p := url.Values{}
		for key, value := range params {

			p.Add(key, value)
		}

		path = path + "?" + p.Encode()
	}

	fmt.Printf("%s %s", method, path)

	//Body path
	b := bytes.NewBufferString("")
	length := 0

	if method != "GET" && body != nil {

		p := url.Values{}

		for key, value := range body {

			p.Add(key, value)
		}

		b = bytes.NewBufferString(p.Encode())
		length = len(p.Encode())
	}

	client := &http.Client{}
	req, _ := http.NewRequest(method, baseurl+path, b)
	req.Header.Add("Accept", "application/vnd.github.v3.full+json")
	req.Header.Add("Authorization", "token "+g.AccessToken)
	if length > 0 {
		req.Header.Add("Content-Length", strconv.Itoa(length))
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}

		cb <- body
	}
}

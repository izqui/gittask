package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	//"io/ioutil"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	baseurl = "https://api.github.com"
)

type Github struct {
	AccessToken string
}

type User struct {
	Id bson.ObjectId `json:"_id,omitempty" bson: "_id"`

	Email       string `json:"email,omitempty" bson:"email"`
	Username    string `json:"login" bson:"username"`
	Location    string `json:"location,omitempty" bson:"location"`
	AccessToken string `json:""bson:"token"`
}

type userCallback chan *User

func (g *Github) UserInfo(username string) *User {

	callback := make(userCallback)

	path := "/user"

	if username != "" && username != "me" {

		path = "/users/" + username
	}

	go g.request("GET", path, nil, nil, callback)

	select {

	case user := <-callback:
		return user

	}

}

func (g *Github) request(method string, path string, params map[string]string, body map[string]string, cb userCallback) {

	//Append URL params
	if params != nil {

		p := url.Values{}
		for key, value := range params {

			p.Add(key, value)
		}

		path = path + "?" + p.Encode()
	}

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

		defer resp.Body.Close()

		user, _ := decodeUser(resp.Body)
		user.AccessToken = g.AccessToken

		cb <- user
	}
}

func decodeUser(r io.Reader) (u *User, err error) {

	u = new(User)
	err = json.NewDecoder(r).Decode(u)
	return
}

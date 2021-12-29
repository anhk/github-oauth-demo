package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

// GitHub App ClientID
const (
	clientId     = "dc7c2a2035486454e400"
	clientSecret = "9bbb968691e817399a5000bf7d83d3f532c9712e"
)

//
func httpGet(url, token string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %v", token))
	}

	if res, err := (&http.Client{}).Do(req); err != nil {
		panic(err)
	} else if data, err := ioutil.ReadAll(res.Body); err != nil {
		panic(err)
	} else {
		res.Body.Close()
		return data
	}
}

//
func httpPost(url, body string) []byte {
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Accept", "application/json")

	if res, err := (&http.Client{}).Do(req); err != nil {
		panic(err)
	} else if data, err := ioutil.ReadAll(res.Body); err != nil {
		panic(err)
	} else {
		res.Body.Close()
		return data
	}
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GithubUser struct {
	Login      string `json:"login"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Id         string `json:"id"`
	NodeId     string `json:"node_id"`
	AvatarUrl  string `json:"avatar_url"`
	GravatarId string `json:"gravatar_id"`
	HtmlUrl    string `json:"html_url"`
}

// processApi
func processApi(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// 获取令牌
	data := httpPost(fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%v&client_secret=%v&code=%v",
		clientId, clientSecret, code), "")
	token := &Token{}
	json.Unmarshal(data, token)

	user := &GithubUser{}
	data2 := httpGet("https://api.github.com/user", token.AccessToken)
	json.Unmarshal(data2, user)
	w.Write(data2)
}

// processHtml
func processHtml(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("./html/index.html"))
	tpl.Execute(w, nil)
}

// 启用Web
func startWeb() {
	router := mux.NewRouter()
	router.HandleFunc("/api/go", processApi)
	router.PathPrefix("/css").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./html/css"))))
	router.HandleFunc("/", processHtml)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

//main
func main() {
	startWeb()
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
		_ = res.Body.Close()
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
		_ = res.Body.Close()
		return data
	}
}

type Token struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	Scope       string `json:"scope,omitempty"`

	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorUri         string `json:"error_uri,omitempty"`
}

type GithubUser struct {
	Login             string    `json:"login"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// processApi
func processApi(ctx *gin.Context) {
	code := ctx.Query("code")

	// 获取令牌
	data := httpPost(fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%v&client_secret=%v&code=%v",
		clientId, clientSecret, code), "")
	token := &Token{}
	_ = json.Unmarshal(data, token)

	if token.Error != "" {
		ctx.IndentedJSON(403, token)
		return
	}

	user := &GithubUser{}
	data2 := httpGet("https://api.github.com/user", token.AccessToken)
	_ = json.Unmarshal(data2, user)
	ctx.IndentedJSON(200, user)
}

// processHtml
func processHtml(ctx *gin.Context) {
	_, _ = ctx.Writer.Write([]byte(`
<head>
    <title>HelloWorld</title>
</head>

<body>
    <a href="https://github.com/login/oauth/authorize?client_id=dc7c2a2035486454e400">Github 登录</a>
</body>
`))
}

// 启用Web
func startWeb() {
	router := gin.Default()
	router.GET("/", processHtml)
	router.GET("/api/go", processApi) // 在Github上注册应用的时候，回调接口设置的是 http://127.1:8080/api/go
	_ = router.Run(":8080")
}

//main
func main() {
	startWeb()
}

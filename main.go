package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GitHub App ClientID
const (
	// GitHub
	//AccessTokenUrl = "https://github.com/login/oauth/access_token"
	//UserInfoUrl    = "https://api.github.com/user"
	//AuthorizeUrl   = "https://github.com/login/oauth/authorize"
	//ClientId       = "dc7c2a2035486454e400"
	//ClientSecret   = "9bbb968691e817399a5000bf7d83d3f532c9712e"

	// AliYun
	//AccessTokenUrl = "https://eiam-api-cn-hangzhou.aliyuncs.com/v2/idaas_abt3pfwojojcq323si6g5tx7e4/app_mkvmrajqr7e6sehnu7vlynet6a/oauth2/token"
	//UserInfoUrl    = "https://eiam-api-cn-hangzhou.aliyuncs.com/v2/idaas_abt3pfwojojcq323si6g5tx7e4/app_mkvmrajqr7e6sehnu7vlynet6a/oauth2/userinfo"
	//AuthorizeUrl   = "https://0hyuelcn.aliyunidaas.com/login/app/app_mkvmrajqr7e6sehnu7vlynet6a/oauth2/authorize"
	//ClientId       = "app_mkvmrajqr7e6sehnu7vlynet6a"
	//ClientSecret   = "CS2mfMVnRXJ11SkAdDYN3aedWdoTb9zjCiMVXdkn3LcvEW"

	// JDCloud
	AccessTokenUrl = "https://oauth2.jdcloud.com/token"
	UserInfoUrl    = "https://oauth2.jdcloud.com/userinfo"
	AuthorizeUrl   = "https://oauth2.jdcloud.com/authorize"
	ClientId       = "9241655362939860"
	ClientSecret   = "CS2mfMVnRXJ1"
)

//
func httpGet(url, token string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
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

type UserInfo struct {
	Login             string `json:"login"` // ==> Github
	Name              string `json:"name"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
}

// processApi
func processApi(ctx *gin.Context) {
	code := ctx.Query("code")

	// 获取令牌
	data := httpPost(fmt.Sprintf("%s?grant_type=authorization_code&client_id=%v&client_secret=%v&code=%v&redirect_uri=",
		AccessTokenUrl, ClientId, ClientSecret, code), "")
	token := &Token{}
	_ = json.Unmarshal(data, token)

	if token.Error != "" {
		ctx.IndentedJSON(403, token)
		return
	}

	user := &UserInfo{}
	data2 := httpGet(UserInfoUrl, token.AccessToken)
	fmt.Println(string(data2))
	_ = json.Unmarshal(data2, user)
	ctx.IndentedJSON(200, user)
}

var HtmlFmt = `
<head>
    <title>HelloWorld</title>
</head>

<body>
    <a href="%s?client_id=%s">OAuth2 登录</a>
</body>
`

// processHtml
func processHtml(ctx *gin.Context) {
	_, _ = ctx.Writer.Write([]byte(fmt.Sprintf(HtmlFmt, AuthorizeUrl, ClientId)))
}

func main() {
	router := gin.Default()
	router.GET("/", processHtml)
	router.GET("/api/go", processApi) // 在Github上注册应用的时候，回调接口设置的是 http://127.1:8080/api/go
	_ = router.Run(":8080")
}

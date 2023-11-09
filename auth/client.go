package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type Client struct {
	Email           string
	Password        string
	Persona         Persona
	CookieSessionId string
	AccessCode      string
	AccessToken     string
	UTSessionId     string
	client          http.Client
	fid             string
	loginurl        url.URL
	Identity        Pid
}

func NewClient(email string, password string) *Client {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		panic("cookie jar fail")
	}

	return &Client{
		Email:    email,
		Password: password,
		client: http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *Client) UTASSession() string {
	return c.UTSessionId
}
func (c *Client) promptOneTimeCode() (string, error) {
	var oneTimeCode string

	fmt.Print("Please enter your 2FA code: ")
	_, err := fmt.Scanln(&oneTimeCode)
	if err != nil {
		return "", err
	}

	oneTimeCode = strings.TrimSpace(oneTimeCode)

	return oneTimeCode, nil
}
func (c *Client) maskEmail(email string) string {
	parts := strings.Split(email, "@")
	localPart := parts[0]
	domainPart := parts[1]

	maskedLocalPart := localPart[:2] + strings.Repeat("*", len(localPart)-2)

	return maskedLocalPart + "@" + domainPart
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (c *Client) initialVisit() error {
	c.client.CheckRedirect = nil
	url, err := url.Parse("https://www.ea.com/de-de/ea-sports-fc/ultimate-team/web-app/")
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req)
	if err != nil {
		return err
	}

	url, err = url.Parse("https://accounts.ea.com/connect/auth?response_type=token&redirect_uri=nucleus%3Arest&prompt=none&client_id=ORIGIN_JS_SDK")
	if err != nil {
		return err
	}

	req, err = http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req)
	if err != nil {
		return err
	}

	url, err = url.Parse("https://accounts.ea.com/connect/auth?client_id=FC24_JS_WEB_APP&response_type=token&display=web2/login&locale=en_US&machineProfileKey=0&redirect_uri=nucleus:rest&prompt=none&release_type=prod&scope=basic.identity+offline+signin+basic.entitlement+basic.persona")
	if err != nil {
		return err
	}

	req, err = http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req)
	if err != nil {
		return err
	}

	c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return nil
}
func (c *Client) initAuthentication() error {
	v := url.Values{}

	url, err := url.Parse("https://accounts.ea.com/connect/auth")
	if err != nil {
		return err
	}
	v.Add("accessToken", "")
	v.Add("client_id", "FC24_JS_WEB_APP")
	v.Add("display", "web2/login")
	v.Add("hide_create", "true")
	v.Add("locale", "en_US")
	v.Add("prompt", "login")
	v.Add("redirect_uri", "https://www.ea.com/de-de/ea-sports-fc/ultimate-team/web-app/auth.html")
	v.Add("release_type", "prod")
	v.Add("response_type", "token")
	v.Add("scope", "basic.identity offline signin basic.entitlement basic.persona")
	url.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	url, err = url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return err
	}
	c.fid = url.Query().Get("fid")

	c.client.CheckRedirect = nil

	req, err = http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}
	resp, err = c.Do(req)
	if err != nil {
		return err
	}

	url, err = url.Parse(resp.Header.Get("Selflocation"))
	if err != nil {
		return err
	}
	c.loginurl = *url
	return nil
}

func (c *Client) loginEmailPassword() error {
	formData := map[string]string{
		"email":                     c.Email,
		"regionCode":                "DE",
		"phoneNumber":               "",
		"password":                  c.Password,
		"_eventId":                  "submit",
		"cid":                       "6bQaSXWmBnlRRsKNQsvgIcWGeZ7Ltgau,GR5E0ATIu78fhdJFy5IpLnPLlG9PoRkv",
		"showAgeUp":                 "true",
		"thirdPartyCaptchaResponse": "",
		"loginMethod":               "emailPassword",
		"_rememberMe":               "on",
		"rememberMe":                "on",
	}

	// Encode the form data
	data := url.Values{}
	for key, value := range formData {
		data.Add(key, value)
	}

	req, err := http.NewRequest(http.MethodPost, c.loginurl.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	formData = map[string]string{
		"codeType":          "EMAIL",
		"maskedDestination": c.maskEmail(c.Email),
		"_eventId":          "submit",
	}
	// Encode the form data
	data = url.Values{}
	for key, value := range formData {
		data.Add(key, value)
	}

	req, err = http.NewRequest(http.MethodPost, resp.Request.URL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = c.Do(req)
	if err != nil {
		return err
	}

	oneTimeCode, err := c.promptOneTimeCode()
	if err != nil {
		return err
	}

	formData = map[string]string{
		"oneTimeCode":      oneTimeCode,
		"_trustThisDevice": "on",
		"trustThisDevice":  "on",
		"_eventId":         "submit",
	}

	// Encode the form data
	data = url.Values{}
	for key, value := range formData {
		data.Add(key, value)
	}
	req, err = http.NewRequest(http.MethodPost, resp.Request.URL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = c.Do(req)
	if err != nil {
		return err
	}

	authUrl, err := url.Parse("https://accounts.ea.com/connect/auth?hide_create=true&display=web2%2Flogin&scope=basic.identity+offline+signin+basic.entitlement+basic.persona&release_type=prod&response_type=token&redirect_uri=https%3A%2F%2Fwww.ea.com%2Fde-de%2Fea-sports-fc%2Fultimate-team%2Fweb-app%2Fauth.html&accessToken=&locale=en_US&prompt=login&client_id=FC24_JS_WEB_APP&fid=")
	if err != nil {
		return err
	}

	v := authUrl.Query()
	v.Set("fid", c.fid)
	authUrl.RawQuery = v.Encode()

	req, err = http.NewRequest(http.MethodGet, authUrl.String(), nil)
	if err != nil {
		return err
	}

	c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err = c.Do(req)
	if err != nil {
		return err
	}

	accesstokenurl, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return err
	}
	frag := accesstokenurl.Fragment
	query, err := url.ParseQuery(frag)
	if err != nil {
		return err
	}
	c.AccessToken = query.Get("access_token")
	query = url.Values{}
	query.Add("token", c.AccessToken)
	accesstokenurl.RawQuery = query.Encode()

	req, err = http.NewRequest(http.MethodGet, accesstokenurl.String(), nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) authCodeOriginJSSDK() error {
	v := url.Values{}

	url, err := url.Parse("https://accounts.ea.com/connect/auth")
	if err != nil {
		return err
	}

	v.Add("response_type", "token")
	v.Add("redirect_uri", "nucleus:rest")
	v.Add("prompt", "none")
	v.Add("client_id", "ORIGIN_JS_SDK")

	url.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	var response AccessTokenResponse
	jsonstring, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonstring, &response)
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) me() error {
	url, err := url.Parse("https://gateway.ea.com/proxy/identity/pids/me")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	var response IdentityMeResponse
	jsonstring, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonstring, &response)
	if err != nil {
		return err
	}
	c.Identity = response.Pid
	return nil
}

func (c *Client) accountInfo() error {
	v := url.Values{}

	url, err := url.Parse("https://utas.mob.v2.fut.ea.com/ut/game/fc24/v2/user/accountinfo")
	if err != nil {
		return err
	}

	v.Add("filterConsoleLogin", "true")
	v.Add("sku", "FUT24WEB")
	v.Add("returningUserGameYear", "2023")
	v.Add("clientVersion", "1")

	url.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Nucleus-Access-Code", c.AccessCode)
	req.Header.Add("Easw-Session-Data-Nucleus-Id", fmt.Sprintf("%d", c.Identity.PidId))
	req.Header.Add("Nucleus-Redirect-Url", "nucleus:rest")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	var response AccountInfoResponse
	jsonstring, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonstring, &response)
	if err != nil {
		return err
	}
	c.Persona = response.UserAccountInfo.Personas[0]
	return nil
}
func (c *Client) accessCode(clientSequence string) error {
	v := url.Values{}

	url, err := url.Parse("https://accounts.ea.com/connect/auth")
	if err != nil {
		return err
	}

	v.Add("client_id", "FUTWEB_BK_OL_SERVER")
	v.Add("redirect_uri", "nucleus:rest")
	v.Add("response_type", "code")
	v.Add("access_token", c.AccessToken)
	v.Add("release_type", clientSequence)
	v.Add("client_sequence", "ut-auth")
	url.RawQuery = v.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	var response AccessCodeResponse
	jsonstring, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonstring, &response)
	if err != nil {
		return err
	}

	c.AccessCode = response.Code
	return nil
}

func (c *Client) utasauth() error {
	data := UTASAuthPayload{
		ClientVersion: 1,
		Ds:            "7b1a102ff145b21177c55e9213312b8722e8f527b2a32874fda55b4efd8c7bf5/2f31", // dynamic value?
		GameSku:       "FFA24PS5",
		Identification: UTASAuthIdentification{
			AuthCode:    c.AccessCode,
			RedirectUrl: "nucleus:rest",
		},
		IsReadOnly:       false,
		Locale:           "en-US",
		Method:           "authcode",
		NucleusPersonaId: c.Persona.PersonaId,
		PriorityLevel:    4,
		Sku:              "FUT24WEB",
	}
	jsondata, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url, err := url.Parse("https://utas.mob.v2.fut.ea.com/ut/auth")
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(jsondata))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	var response UTASAuthResponse
	jsonstring, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonstring, &response)
	if err != nil {
		return err
	}

	c.UTSessionId = response.SID
	return nil
}

func (c *Client) Login() error {
	if err := c.initialVisit(); err != nil {
		return err
	}

	if err := c.initAuthentication(); err != nil {
		return err
	}

	if err := c.loginEmailPassword(); err != nil {
		return err
	}

	if err := c.authCodeOriginJSSDK(); err != nil {
		return err
	}

	if err := c.me(); err != nil {
		return err
	}

	if err := c.accessCode("shard5"); err != nil {
		return err
	}

	if err := c.accountInfo(); err != nil {
		return err
	}

	if err := c.accessCode("prod"); err != nil {
		return err
	}

	if err := c.utasauth(); err != nil {
		return err
	}
	return nil
}

/*
Copyright [2016] [mercadolibre.com]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


**

This package allows you to interact with the Mercadolibre open platform API.
The are two main structures:
1) Client
2) Authorization

1) - This structure keeps within the secret to be used for generating the token to be sent when calling to the private APIs.
     This also provides several methods to call either public and private APIs

2) - This structure keeps the tokens and their expiration time and has to be passed by param each time a call has to be performed to any private API.
*/

package sdk

import (
    "net/url"
    "strconv"
    "bytes"
    "net/http"
    "fmt"
    "io"
    "encoding/json"
    "io/ioutil"
    "log"
    "strings"
    "errors"
    "time"
)

const (

    MLA = "https://auth.mercadolibre.com.ar" // Argentina
    MLB = "https://auth.mercadolivre.com.br" // Brasil
    MCO = "https://auth.mercadolibre.com.co" // Colombia
    MCR = "https://auth.mercadolibre.com.cr" // Costa Rica
    MEC = "https://auth.mercadolibre.com.ec" // Ecuador
    MLC = "https://auth.mercadolibre.cl"      // Chile
    MLM = "https://auth.mercadolibre.com.mx" // Mexico
    MLU = "https://auth.mercadolibre.com.uy" // Uruguay
    MLV = "https://auth.mercadolibre.com.ve" // Venezuela
    MPA = "https://auth.mercadolibre.com.pa" // Panama
    MPE = "https://auth.mercadolibre.com.pe" // Peru
    MPT = "https://auth.mercadolivre.pt"      // Portugal
    MRD = "https://auth.mercadolibre.com.do" // Dominicana

    AUTHORIZATION_CODE = "authorization_code"
    REFRESH_TOKEN = "refresh_token"

)
var ANONYMOUS = Authorization{}
//var refTokenMux = &sync.Mutex{}

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Client struct {
    apiUrl string
    clientId int64
    clientSecret string
}
/*
clientId and clientSecret are  generated when you create your application
*/
func NewClient(clientId int64, clientSecret string) *Client {

    client := new(Client)
    client.apiUrl = "https://api.mercadolibre.com"
    client.clientId = clientId
    client.clientSecret = clientSecret

    return client
}

func (client *Client) SetApiURL(url string) {
    client.apiUrl = url
}

func (client Client) GetAuthURL(base_site, callback string ) string {

    authURL := NewAuthorizationURL(base_site  + "/authorization")
    authURL.addResponseType("code")
    authURL.addClientId(client.clientId)
    authURL.addRedirectUri(callback)

    return authURL.string()
}

func (client Client) Authorize(code, redirectUri string) (*Authorization, error) {

    authURL := NewAuthorizationURL(client.apiUrl + "/oauth/token")
    authURL.addGrantType(AUTHORIZATION_CODE)
    authURL.addClientId(client.clientId)
    authURL.addClientSecret(client.clientSecret)
    authURL.addCode(code)
    authURL.addRedirectUri(redirectUri)

    resp, err := http.Post(authURL.string(), "application/json", *(new(io.Reader)))

    if err != nil {
        fmt.Printf("Error when posting: %s", err)
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("There was an error while authorizing. Check wether your code has not expired.")
    }

    body, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    authorization := new(Authorization)
    if err := json.Unmarshal(body, authorization); err != nil {
        log.Printf("Error while receiving the authorization %s %s", err.Error(), body)
        return nil, err
    }

    authorization.Received_at = time.Now().Unix()
    return authorization, nil
}

func (client Client) Get(resource_path string, authorization Authorization) (*http.Response, error) {

    apiUrl := NewAuthorizationURL(client.apiUrl + resource_path)

    if authorization != ANONYMOUS {
        apiUrl.addAccessToken(authorization.Access_token)
    }

    resp, err := http.Get(apiUrl.string())
    if err != nil {
        fmt.Printf("Error while calling url: %s \n Error: %s", apiUrl.string(), err.Error())
        return nil, err
    }

   /* if resp.StatusCode == http.StatusNotFound {
        err := refreshTokenIfNeeded(client, authorization)
        //err := client.RefreshToken(authorization)
        if err != nil {
            log.Printf("Error while refreshing token %s\n", err.Error())
            return nil, err
        }

        apiUrl := NewAuthorizationURL(client.apiUrl + resource_path)
        apiUrl.addAccessToken(authorization.Access_token)
        resp, err = http.Get(apiUrl.string())

        if err != nil {
            log.Printf("Error while calling API %s\n", err.Error())
            return nil, err
        }
    }*/

    return resp, nil
}


func (client Client) Post(resource_path string, authorization Authorization, body string) (*http.Response, error){

    apiUrl := NewAuthorizationURL(client.apiUrl + resource_path)

    if authorization != ANONYMOUS {
        apiUrl.addAccessToken(authorization.Access_token)
    }

    resp, err := http.Post(apiUrl.string(), "application/json", bytes.NewReader([]byte(body)))

    if err != nil {
        fmt.Printf("Error while calling url: %s \n Error: %s", apiUrl.string(), err)
        return nil, err
    }

  /*  if resp.StatusCode == http.StatusNotFound {

        err := refreshTokenIfNeeded(client, authorization)
        //err :=client.RefreshToken(authorization)
        if err != nil {
            log.Printf("Error while refreshing token %s\n", err.Error())
            return nil, err
        }

        apiUrl := NewAuthorizationURL(client.apiUrl + resource_path)
        apiUrl.addAccessToken(authorization.Access_token)

        resp, err = http.Post(apiUrl.string(), "application/json", bytes.NewReader([]byte(body)))

        if err != nil {
            log.Printf("Error while calling API %s\n", err)
            return nil, err
        }
    }*/

    return resp, nil
}

func (client Client) Put(resource_path string, authorization Authorization, body *string) (*http.Response, error){

    apiUrl := NewAuthorizationURL(client.apiUrl + resource_path)

    if authorization != ANONYMOUS {
        apiUrl.addAccessToken(authorization.Access_token)
    }

    req, err := http.NewRequest(http.MethodPut, apiUrl.string(), strings.NewReader(*body))
    if err != nil {
        log.Printf("Error when creating PUT request %d.", err)
        return nil, err
    }

    req.Header.Add("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        fmt.Printf("Error while calling url: %s\n Error: %s", apiUrl.string(), err)
        return nil, err
    }

   /* if resp.StatusCode == http.StatusNotFound {

        //err := client.RefreshToken(authorization)
        err := refreshTokenIfNeeded(client, authorization)
        if err != nil {
            log.Printf("Error while refreshing token %s\n", err.Error())
            return nil, err
        }

        apiUrl := NewAuthorizationURL(client.apiUrl + resource_path)
        apiUrl.addAccessToken(authorization.Access_token)

        req, err = http.NewRequest(http.MethodPut, apiUrl.string(), strings.NewReader(*body))
        if err != nil {
            log.Printf("Error when creating PUT request %d.", err)
            return nil, err
        }

        req.Header.Add("Content-Type", "application/json")
        resp, err = http.DefaultClient.Do(req)
    }*/

    return resp, nil
}

func (client Client) Delete(resource_path string, authorization Authorization) (*http.Response, error) {

    apiUrl := NewAuthorizationURL(client.apiUrl + resource_path)

    if authorization != ANONYMOUS {
        apiUrl.addAccessToken(authorization.Access_token)
    }

    req, err := http.NewRequest(http.MethodDelete, apiUrl.string(), nil)
    if err != nil {
        log.Printf("Error when creating PUT request %d.", err)
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        fmt.Printf("Error while calling url: %s \n Error: %s", apiUrl.string(), err)
        return nil, err
    }

    /*if resp.StatusCode == http.StatusNotFound {

        //err := client.RefreshToken(authorization)
        err := refreshTokenIfNeeded(client, authorization)
        if err != nil {
            log.Printf("Error while refreshing token %s\n", err.Error())
            return nil, err
        }

        apiUrl := NewAuthorizationURL(client.apiUrl + resource_path)
        apiUrl.addAccessToken(authorization.Access_token)

        req, err = http.NewRequest(http.MethodDelete, apiUrl.string(), nil)
        if err != nil {
            log.Printf("Error when creating PUT request %d.", err)
            return nil, err
        }

        resp, err = http.DefaultClient.Do(req)
    }*/

    return resp, nil
}


func (client Client) RefreshToken(authorization Authorization) (*Authorization, error) {

    authorizationURL := NewAuthorizationURL(client.apiUrl + "/oauth/token")
    authorizationURL.addGrantType(REFRESH_TOKEN)
    authorizationURL.addClientId(client.clientId)
    authorizationURL.addClientSecret(client.clientSecret)
    authorizationURL.addRefreshToken(authorization.Refresh_token)

    log.Printf(authorizationURL.string())
    resp, err := http.Post(authorizationURL.string(), "application/json", *(new(io.Reader)))

    if err != nil {
        log.Printf("Error while refreshing token: %s\n", err.Error())
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("Refreshing token returned status code " + resp.Status)
    }

    body, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    newAuth := authorization
    if err := json.Unmarshal(body, &newAuth); err != nil {
        log.Printf("Error while receiving the authorization %s %s", err.Error(), body)
        return nil, err
    }

    newAuth.Received_at = time.Now().Unix()

    return &newAuth, nil
}

/*
If a refresh token is present in the authorization code exchange, then it may be used to obtain a new access tokens at any time.
*/

type Authorization struct {
    Access_token string
    Token_type string
    Expires_in int16
    Received_at int64
    Refresh_token string
    Scope string
}

func (auth Authorization) isExpired() bool {

    return ((auth.Received_at + int64(auth.Expires_in)) <= time.Now().Unix())
}


/*
This struct allows adding all the params needed to the URL to be sent
*/
type AuthorizationURL struct{
    url bytes.Buffer
}

func (u *AuthorizationURL) addGrantType(value string) {
    u.add("grant_type=" + value)
}

func (u *AuthorizationURL) addClientId(value int64) {
    u.add("client_id=" + strconv.FormatInt(value, 10))
}

func (u *AuthorizationURL) addClientSecret(value string) {
    u.add("client_secret=" + url.QueryEscape(value))
}

func (u *AuthorizationURL) addCode(value string) {
    u.add("code=" + url.QueryEscape(value))
}

func (u *AuthorizationURL) addRedirectUri(uri string) {
    u.add("redirect_uri=" + url.QueryEscape(uri))
}

func (u *AuthorizationURL) addRefreshToken(t string) {
    u.add("refresh_token=" + url.QueryEscape(t))
}

func (u *AuthorizationURL) addResponseType(value string) {
    u.add("response_type=" + url.QueryEscape(value))
}

func (u *AuthorizationURL) addAccessToken(t string){
    u.add("access_token=" + url.QueryEscape(t))
}

func (u *AuthorizationURL) string() string {
    return u.url.String()
}

func (u *AuthorizationURL) add(value string) {

    if !strings.Contains(u.url.String(), "?"){
        u.url.WriteString("?" + value)
    } else if strings.LastIndex("&",u.url.String()) >= u.url.Len(){
        u.url.WriteString(value)
    } else {
        u.url.WriteString("&" + value)
    }
}

func NewAuthorizationURL(baseURL string) *AuthorizationURL{
    authURL := new(AuthorizationURL)
    authURL.url.WriteString(baseURL)
    return authURL
}
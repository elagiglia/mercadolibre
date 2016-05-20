//This class allows you to interact with the Mercadolibre open platform API.
//

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
)

const (

	MLA = "https://auth.mercadolibre.com.ar" // Argentina
	/*MLB = "https://auth.mercadolivre.com.br" // Brasil
	MCO = "https://auth.mercadolibre.com.co" // Colombia
	MCR = "https://auth.mercadolibre.com.cr" // Costa Rica
	MEC = "https://auth.mercadolibre.com.ec" // Ecuador
	MLC = "https://auth.mercadolibre.cl" 	 // Chile
	MLM = "https://auth.mercadolibre.com.mx" // Mexico
	MLU = "https://auth.mercadolibre.com.uy" // Uruguay
	MLV = "https://auth.mercadolibre.com.ve" // Venezuela
	MPA = "https://auth.mercadolibre.com.pa" // Panama
	MPE = "https://auth.mercadolibre.com.pe" // Peru
	MPT = "https://auth.mercadolivre.pt" 	 // Portugal
	MRD = "https://auth.mercadolibre.com.do" // Dominicana*/
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Client struct {
	apiUrl string
	clientId int64
	clientSecret string
}

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

	var buffer bytes.Buffer
	buffer.WriteString("/authorization?response_type=code&client_id=")
	buffer.WriteString(strconv.FormatInt(client.clientId, 10))
	buffer.WriteString("&redirect_uri=")

	base_url := base_site + buffer.String()

	encoded_callback := url.QueryEscape(callback)

	full_url := base_url + encoded_callback

	return full_url
}

func (client Client) Authorize(code, redirectUri string) (Authorization, error) {

	var params bytes.Buffer
	params.WriteString("grant_type=authorization_code")
	params.WriteString("&client_id=")
	params.WriteString(strconv.FormatInt(client.clientId, 10))
	params.WriteString("&client_secret=" + url.QueryEscape(client.clientSecret))
	params.WriteString("&code=" + url.QueryEscape(code))
	params.WriteString("&redirect_uri=" + url.QueryEscape(redirectUri))

	final_url := client.apiUrl + "/oauth/token?" + params.String()

	authorization := new(Authorization)
	resp, err := http.Post(final_url, "application/json", *(new(io.Reader)))

	if err != nil {
		fmt.Printf("Error when posting: %s", err)
		return *authorization, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, authorization); err != nil {
		log.Printf("Error while receiving the authorization %s %s", err, body)
		return *authorization, err
	}

	return *authorization, nil
}

func (client Client) Get(resource_path string, authorization *Authorization) (*http.Response, error) {

	base_url := client.apiUrl + resource_path
	final_url := base_url

	if authorization != nil {
		final_url = base_url + "?access_token=" + url.QueryEscape(authorization.Access_token)
	}

	resp, err := http.Get(final_url)
	if err != nil {
		fmt.Printf("Error while calling url: %s \n Error: %s", final_url, err)
		return resp, err
	}

	if resp.StatusCode == http.StatusNotFound {

		client.RefreshToken(authorization)

		resp, err = http.Get(base_url + "?access_token=" + url.QueryEscape(authorization.Access_token))

		if err != nil {
			log.Printf("Error while calling API %s\n", err)
			return resp, err
		}
	}

	return resp, nil
}

//TODO: Try to return an Authorization object instead of changing the original one passed by param.
func (client Client) RefreshToken(authorization *Authorization) error {

	log.Printf("Refreshing token\n")

	var base_url bytes.Buffer
	base_url.WriteString(client.apiUrl)
	base_url.WriteString("/oauth/token?")

	base_url.WriteString("grant_type=refresh_token")
	base_url.WriteString("&client_id=")
	base_url.WriteString(strconv.FormatInt(client.clientId, 10))
	base_url.WriteString("&client_secret=" + url.QueryEscape(client.clientSecret))
	base_url.WriteString("&refresh_token=" + url.QueryEscape(authorization.Refresh_token))

	resp, err := http.Post(base_url.String(), "application/json", *(new(io.Reader)))

	if err != nil || resp.StatusCode != http.StatusOK {

		log.Printf("Error while refreshing token: http status: %s err: %s\n", resp.StatusCode, err)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, authorization); err != nil {
		log.Printf("Error while receiving the authorization %s %s", err, body)
		return err
	}

	return nil
}

func (client Client) Post(resource_path string, authorization *Authorization, body string) (*http.Response, error){

	base_url := client.apiUrl + resource_path
	final_url := base_url

	if authorization != nil {
		final_url = base_url + "?access_token=" + url.QueryEscape(authorization.Access_token)
	}


	resp, err := http.Post(final_url,"application/json", bytes.NewReader([]byte(body)))

	if err != nil {
		fmt.Printf("Error while calling url: %s \n Error: %s", final_url, err)
		return resp, err
	}

	if resp.StatusCode == http.StatusNotFound {

		client.RefreshToken(authorization)
		resp, err = http.Post(base_url + "?access_token=" + url.QueryEscape(authorization.Access_token), "application/json", bytes.NewReader([]byte(body)))

		if err != nil {
			log.Printf("Error while calling API %s\n", err)
			return resp, err
		}
	}

	return resp, nil
}
func (client Client) Put(resource_path string, authorization *Authorization, body string) (*http.Response, error){

	base_url := client.apiUrl + resource_path
	final_url := base_url

	if authorization != nil {
		final_url = base_url + "?access_token=" + url.QueryEscape(authorization.Access_token)
	}


	req, err := http.NewRequest(http.MethodPut, final_url, strings.NewReader(body))
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)


	if err != nil {
		fmt.Printf("Error while calling url: %s \n Error: %s", final_url, err)
		return resp, err
	}

	if resp.StatusCode == http.StatusNotFound {

		client.RefreshToken(authorization)
		req, err = http.NewRequest(http.MethodPut, base_url + "?access_token=" + url.QueryEscape(authorization.Access_token), strings.NewReader(body))
		if err != nil {
			log.Printf("Error when creating PUT request %d.", err)
			return nil, err
		}

		req.Header.Add("Content-Type", "application/json")
		resp, err = http.DefaultClient.Do(req)
	}

	return resp, nil
}

func (client Client) Delete(resource_path string, authorization *Authorization) (*http.Response, error) {

	base_url := client.apiUrl + resource_path
	final_url := base_url

	if authorization != nil {
		final_url = base_url + "?access_token=" + url.QueryEscape(authorization.Access_token)
	}

	req, err := http.NewRequest(http.MethodDelete, final_url, nil)
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)


	if err != nil {
		fmt.Printf("Error while calling url: %s \n Error: %s", final_url, err)
		return resp, err
	}

	if resp.StatusCode == http.StatusNotFound {

		client.RefreshToken(authorization)
		req, err = http.NewRequest(http.MethodDelete, base_url + "?access_token=" + url.QueryEscape(authorization.Access_token), nil)
		if err != nil {
			log.Printf("Error when creating PUT request %d.", err)
			return nil, err
		}

		resp, err = http.DefaultClient.Do(req)
	}

	return resp, nil
}
/*
func refreshIfNeeded(resp http.Response, ) {

	if resp.StatusCode == http.StatusNotFound {

		client.RefreshToken(authorization)

		resp, err = http.Get(base_url + "?access_token=" + url.QueryEscape(authorization.Access_token))

		if err != nil {
			log.Printf("Error while calling API %s\n", err)
			return resp, err
		}
	}
}*/


type Authorization struct {
	Access_token string
	Token_type string
	Expires_in int16
	Refresh_token string
	Scope string
}
package sdk

import (
	"testing"
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
)

const (
	API_TEST = "http://localhost:3000"
	CLIENT_ID = 123456
	CLIENT_SECRET = "client secret"
)

func Test_URL_for_authentication_is_properly_return(t *testing.T) {

	expectedUrl := "https://auth.mercadolibre.com.ar/authorization?response_type=code&client_id=123456&redirect_uri=http%3A%2F%2Fsomeurl.com"

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	url := client.GetAuthURL(MLA, "http://someurl.com")

	if url != expectedUrl {
		log.Printf("Error: The URL is different from the one that was expected.")
		log.Printf("expected %s", expectedUrl)
		log.Printf("obtained %s", url)
		t.FailNow()
	}

}

func Test_That_Authorization_Process_Works(t *testing.T) {

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)

	resp, err := client.Authorize("valid code with refresh token", "http://someurl.com")

	fmt.Printf("Access_token: %s \nRefresh_token: %s", resp.Access_token, resp.Refresh_token)

	if err != nil {
		t.FailNow()
	}
	if resp.Access_token != "valid token" {
		fmt.Errorf("Access_token is not what was expected")
		t.FailNow()
	}

	if resp.Refresh_token != "valid refresh token" {
		fmt.Errorf("Refresh_token is not what was expected")
		t.FailNow()
	}

}

func Test_GET_public_API_sites_works_properly ( t *testing.T){

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)

	//Public APIs do not need Authorization
	resp, err := client.Get("/sites", new (Authorization))

	if err != nil {
		t.FailNow()
	}

	if resp.StatusCode != http.StatusOK {
		t.FailNow()
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil || string(body) == ""{
		t.FailNow()
	}

	fmt.Printf("body %s", body)
}

func Test_GET_private_API_users_works_properly (t *testing.T){

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)

	authorization := Authorization{Access_token:"expired token", Refresh_token:"valid refresh token"}

	resp, err := client.Get("/users/me", &authorization)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		t.FailNow()
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("\nThe expected http status code is 200 but the received one was: %s\n", resp.Status)
		t.FailNow()
	}
}

func Test_GET_private_API_users_returns_an_error_when_refresh_token_is_not_valid (t *testing.T){

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)
	//client := Client{clientId:123456, clientSecret:"client secret", apiUrl:"http://localhost:3000"}
	authorization := Authorization{Access_token:"expired token", Refresh_token:"no valid"}

	resp, err := client.Get("/users/me", &authorization)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		t.FailNow()
	}

	if resp.StatusCode != http.StatusNotFound {
		fmt.Printf("\nThe expected http status code is 200 but the received one was: %s\n", resp.Status)
		t.FailNow()
	}
}

func Test_POST_a_new_item_works_properly_when_token_IS_EXPIRED(t *testing.T){

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)
	authorization := Authorization{Access_token:"expired token", Refresh_token:"valid refresh token"}

	body := "{\"foo\":\"bar\"}"
	resp, err := client.Post("/items", &authorization, body)

	if err != nil {
		log.Printf("Error while posting a new item %s\n", err)
		t.FailNow()
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Error while posting a new item status code: %s\n", resp.StatusCode)
		t.FailNow()
	}
}

func Test_POST_a_new_item_works_properly_when_token_IS_NOT_EXPIRED (t *testing.T){

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)
	authorization := Authorization{Access_token:"valid token", Refresh_token:"valid refresh token"}

	body := "{\"foo\":\"bar\"}"
	resp, err := client.Post("/items", &authorization, body)

	if err != nil {
		log.Printf("Error while posting a new item %s\n", err)
		t.FailNow()
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Error while posting a new item status code: %s\n", resp.StatusCode)
		t.FailNow()
	}
}

func Test_PUT_a_new_item_works_properly_when_token_IS_NOT_EXPIRED (t *testing.T){

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)
	authorization := Authorization{Access_token:"valid token", Refresh_token:"valid refresh token"}

	body := "{\"foo\":\"bar\"}"
	resp, err := client.Put("/items/123", &authorization, body)

	if err != nil {
		log.Printf("Error while posting a new item %s\n", err)
		t.FailNow()
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error while putting a new item. Status code: %s\n", resp.StatusCode)
		t.FailNow()
	}
}

func Test_PUT_a_new_item_works_properly_when_token_IS_EXPIRED (t *testing.T){

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)
	authorization := Authorization{Access_token:"expired token", Refresh_token:"valid refresh token"}

	body := "{\"foo\":\"bar\"}"
	resp, err := client.Put("/items/123", &authorization, body)

	if err != nil {
		log.Printf("Error while posting a new item %s\n", err)
		t.FailNow()
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error while putting a new item. Status code: %s\n", resp.StatusCode)
		t.FailNow()
	}
}

func Test_DELETE_an_item_returns_200_when_token_IS_NOT_EXPIRED (t *testing.T){

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)
	authorization := Authorization{Access_token:"valid token", Refresh_token:"valid refresh token"}

	//body := "{\"foo\":\"bar\"}"
	resp, err := client.Delete("/items/123", &authorization)

	if err != nil {
		log.Printf("Error while deleting an item %s\n", err)
		t.FailNow()
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error while putting a new item. Status code: %s\n", resp.StatusCode)
		t.FailNow()
	}
}

func Test_DELETE_an_item_returns_200_when_token_IS_EXPIRED (t *testing.T){

	client := NewClient(CLIENT_ID, CLIENT_SECRET)
	client.SetApiURL(API_TEST)
	authorization := Authorization{Access_token:"expired token", Refresh_token:"valid refresh token"}

	//body := "{\"foo\":\"bar\"}"
	resp, err := client.Delete("/items/123", &authorization)

	if err != nil {
		log.Printf("Error while deleting an item %s\n", err)
		t.FailNow()
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error while putting a new item. Status code: %s\n", resp.StatusCode)
		t.FailNow()
	}
}
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
*/

package sdk

import (
    "testing"
    "log"
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
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

func Test_That_Authorization_Process_Returns_An_Error_When_Code_Has_Expired(t *testing.T) {

    client := NewClient(CLIENT_ID, CLIENT_SECRET)
    client.SetApiURL(API_TEST)

    _, err := client.Authorize("bad code", "http://someurl.com")

    if err == nil {
        t.FailNow()
    }
    if err.Error() != "There was an error while authorizing. Check wether your code has not expired." {
        t.FailNow()
    }

}

func Test_GET_public_API_sites_works_properly ( t *testing.T){

    client := NewClient(CLIENT_ID, CLIENT_SECRET)
    client.SetApiURL(API_TEST)

    //Public APIs do not need Authorization
    resp, err := client.Get("/sites", ANONYMOUS)

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

    resp, err := client.Get("/users/me", authorization)

    if err != nil {
        fmt.Printf("Error: %s\n", err)
        t.FailNow()
    }

    if resp.StatusCode != http.StatusOK {
        newAuth, _ := client.RefreshToken(authorization)
        _, err := client.Get("/users/me", *newAuth)
        if err != nil {
            t.FailNow()
        }
    }
}

/*func Test_GET_private_API_users_returns_an_error_when_refresh_token_is_not_valid (t *testing.T){

    client := NewClient(CLIENT_ID, CLIENT_SECRET)
    client.SetApiURL(API_TEST)

    authorization := Authorization{Access_token:"expired token", Refresh_token:"no valid"}

    _, err := client.Get("/users/me", authorization)

    if err == nil {
        fmt.Printf("Error should not be nil")
        t.FailNow()
    }
}*/

func Test_POST_a_new_item_works_properly_when_token_IS_EXPIRED(t *testing.T){

    client := NewClient(CLIENT_ID, CLIENT_SECRET)
    client.SetApiURL(API_TEST)
    authorization := Authorization{Access_token:"expired token", Refresh_token:"valid refresh token"}

    body := "{\"foo\":\"bar\"}"
    resp, err := client.Post("/items", authorization, body)

    if err != nil {
        log.Printf("Error while posting a new item %s\n", err)
        t.FailNow()
    }

    if resp.StatusCode != http.StatusOK {
        newAuth, _ := client.RefreshToken(authorization)
        resp, err = client.Post("/items", *newAuth, body)
        if err != nil {
            t.FailNow()
        }
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
    resp, err := client.Post("/items", authorization, body)

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
    resp, err := client.Put("/items/123", authorization, &body)

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
    resp, err := client.Put("/items/123", authorization, &body)

    if resp.StatusCode != http.StatusOK {
        newAuth, _ := client.RefreshToken(authorization)
        resp, err = client.Put("/items/123", *newAuth, &body)
        if err != nil {
            t.FailNow()
        }
    }
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

    resp, err := client.Delete("/items/123", authorization)

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

    resp, err := client.Delete("/items/123", authorization)

    if err != nil {
        log.Printf("Error while deleting an item %s\n", err)
        t.FailNow()
    }

    if resp.StatusCode != http.StatusOK {
        newAuth, _ := client.RefreshToken(authorization)
        resp, err = client.Delete("/items/123", *newAuth)
        if err != nil {
            t.FailNow()
        }
    }
    if resp.StatusCode != http.StatusOK {
        log.Printf("Error while putting a new item. Status code: %s\n", resp.StatusCode)
        t.FailNow()
    }
}

func Test_AuthorizationURL_adds_a_params_separator_when_needed(t *testing.T)  {
    client := NewClient(1234, "abcdedfadafas")
    auth := NewAuthorizationURL(client.apiUrl + "/authorizationauth")
    auth.addGrantType(AUTHORIZATION_CODE)

    url := client.apiUrl + "/authorizationauth?" + "grant_type=" + AUTHORIZATION_CODE

    if strings.Compare(url, auth.string()) != 0 {
        log.Printf("url was different from what was expected\n expected: %s \n obtained: %s \n", url, auth.string())
        t.FailNow()
    }
}

func Test_AuthorizationURL_adds_a_query_param_separator_when_needed(t *testing.T)  {
    client := NewClient(1234, "abcdedfadafas")
    auth := NewAuthorizationURL(client.apiUrl + "/authorizationauth")
    auth.addGrantType(AUTHORIZATION_CODE)
    auth.addClientId(1213213)

    url := client.apiUrl + "/authorizationauth?" + "grant_type=" + AUTHORIZATION_CODE + "&client_id=1213213"

    if strings.Compare(url, auth.string()) != 0 {
        log.Printf("url was different from what was expected\n expected: %s \n obtained: %s \n", url, auth.string())
        t.FailNow()
    }
}

func Test_only_one_token_refresh_call_is_done_when_several_threads_are_executed(t *testing.T){

    client := NewClient(CLIENT_ID, CLIENT_SECRET)
    client.SetApiURL(API_TEST)
    authorization := Authorization{Access_token:"expired token", Refresh_token:"valid refresh token"}

    for i := 0; i< 10 ; i++ {

        go client.Get("/users/me", authorization)

    }
}
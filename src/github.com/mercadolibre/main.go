package main

import (
	"github.com/mercadolibre/sdk"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
)


func main() {


	/*Example 1)
	  Getting the URL to call for authenticating purposes
	  Once you generate the URL and call it, you will be redirect to a ML login page where your credentials will be asked. Then, after
	  entering your credentials you will obtained a CODE which will be used to get all the authorization tokens.
	*/
	client := sdk.NewClient(396051416295796, "qM66avGpv5rcQxNWF4sno5oH7Cjph0I7")
	url := client.GetAuthURL(sdk.MLA,"https://www.example.com")
	fmt.Printf("url:%s\n", url)


	/* Example 2)
	  To get all the tokens which will allow you to access the APIs, you need to call the Authorize method. As parameter you need to
	  use the  CODE returned in the previous example.
	*/

	authorization, err := client.Authorize("TG-574473f0e4b07d84ecb42f00-214509008","https://www.example.com")

	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	js, err := json.Marshal(authorization)
	fmt.Printf("Token:%s\n", js)

	/*
	Example 3)
	Call a private API by using the authorization tokes obatained in the previous example.
	*/

	resp, err := client.Get("/users/me", authorization)

	if err != nil {
		log.Printf("Error %s\n", err.Error())
	}
	userInfo, _:= ioutil.ReadAll(resp.Body)
	fmt.Printf("response:%s\n", userInfo)


	/*
	 Example 4) Refreshing token
	  This becomes necessary when an access token is no longer valid and when you need to make it valid again.
	  As it is stated in OpenPlatform doc: to make this possible, you need to set up offline_access on the application manager.
	*/

	//ATTENTION: authorization param is going to be modified by RefreshToken method.
	client.RefreshToken(authorization)
	js, err = json.Marshal(authorization)
	fmt.Printf("Token:%s\n", js)

	/*
	  Example 5)
	  This example shows you how to POST (publish) a new Item.
	 */

	body :=	"{\"title\":\"Item de test - No Ofertar\",\"category_id\":\"MLA1912\",\"price\":10,\"currency_id\":\"ARS\",\"available_quantity\":1,\"buying_mode\":\"buy_it_now\",\"listing_type_id\":\"bronze\",\"condition\":\"new\",\"description\": \"Item:,  Ray-Ban WAYFARER Gloss Black RB2140 901  Model: RB2140. Size: 50mm. Name: WAYFARER. Color: Gloss Black. Includes Ray-Ban Carrying Case and Cleaning Cloth. New in Box\",\"video_id\": \"YOUTUBE_ID_HERE\",\"warranty\": \"12 months by Ray Ban\",\"pictures\":[{\"source\":\"http://upload.wikimedia.org/wikipedia/commons/f/fd/Ray_Ban_Original_Wayfarer.jpg\"},{\"source\":\"http://en.wikipedia.org/wiki/File:Teashades.gif\"}]}"

	resp, err = client.Post("/items", authorization, body)

	if err != nil {
		log.Printf("Error %s\n", err.Error())
	}
	userInfo, _= ioutil.ReadAll(resp.Body)
	fmt.Printf("response:%s\n", userInfo)

	/*
	  Example 6)
	  This example shows you how to post a new Item.
	 */


}


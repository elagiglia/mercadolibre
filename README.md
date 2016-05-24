# MercadoLibre's Golang SDK [DRAFT - WORK IN PROGRESS]

This is the official GO SDK for MercadoLibre's Platform.

## How do I install it?

You can download the latest build at:
    https://github.com/mercadolibre/go-sdk/archive/master.zip

How do I install it using go:

Just run the following command within your $GOPATH

```bash
go get github.com/mercadolibre/sdk
```

And that's it!

## How do I start using it?

The first thing to do is to instance a ```Client``` object. You'll need to give a ```clientId``` and a ```clientSecret```. You can obtain both after creating your own application. For more information on this please read: [creating an application](http://developers.mercadolibre.com/application-manager/)

```go
client := sdk.NewClient(123456,"client secret")
```
With this instance you can start working on MercadoLibre's APIs.

There are some design considerations worth to mention.
This SDK is just a thin layer on top of an http client to handle all the OAuth WebServer flow for you.


## How do I redirect users to authorize my application?

This is a 2 step process.

First get the link to redirect the user. This is very easy! Just:

```go
client := sdk.NewClient(191151416211796, "qM66avGpv5rcQxNWF4sno5oH7Cjph0I7")
url := client.GetAuthURL(sdk.MLA,"https://www.example.com")
fmt.Printf("url:%s\n", url)
```

This will give you the url to redirect the user. You need to specify a callback url which will be the one that the user will redirected after a successfull authrization process.

Once the user is redirected to your callback url, you'll receive in the query string, a parameter named ```code```. You'll need this for the second part of the process.

```go
authorization, err := client.Authorize("TG-57445a71e4b0744714824b93-19793657","https://www.example.com")

if err != nil {
    log.Printf("err: %s", err.Error())
    return
}

js, err := json.Marshal(authorization)
fmt.Printf("Token:%s\n", js)
```

This will get an ```accessToken``` and a ```refreshToken``` (is case your application has the ```offline_access```) for your application and your user.

At this stage your are ready to make call to the API on behalf of the user.

## Making GET calls

```GO
resp, err := client.Get("/users/me", authorization)

if err != nil {
	log.Printf("Error %s\n", err.Error())
}
userInfo, _:= ioutil.ReadAll(resp.Body)
fmt.Printf("response:%s\n", userInfo)

```

## Making POST calls

```GO
body :=	"{\"title\":\"Item de test - No Ofertar\",\"category_id\":\"MLA1912\",\"price\":10,\"currency_id\":\"ARS\",\"available_quantity\":1,\"buying_mode\":\"buy_it_now\",\"listing_type_id\":\"bronze\",\"condition\":\"new\",\"description\": \"Item:,  Ray-Ban WAYFARER Gloss Black RB2140 901  Model: RB2140. Size: 50mm. Name: WAYFARER. Color: Gloss Black. Includes Ray-Ban Carrying Case and Cleaning Cloth. New in Box\",\"video_id\": \"YOUTUBE_ID_HERE\",\"warranty\": \"12 months by Ray Ban\",\"pictures\":[{\"source\":\"http://upload.wikimedia.org/wikipedia/commons/f/fd/Ray_Ban_Original_Wayfarer.jpg\"},{\"source\":\"http://en.wikipedia.org/wiki/File:Teashades.gif\"}]}"

resp, err = client.Post("/items", authorization, body)

if err != nil {
    log.Printf("Error %s\n", err.Error())
}
userInfo, _= ioutil.ReadAll(resp.Body)
fmt.Printf("response:%s\n", userInfo)

```
## Making PUT calls

```GO

```
## Making DELETE calls

```GO
client := sdk.NewClient(123456,"client secret")
client.Delete("/items/123", authorization)
```

## Do I always need to include the ```access_token``` as a parameter?
No. Actually most ```GET``` requests don't need an ```access_token``` and it is easier to avoid them and also it is better in terms of caching.
But this decision is left to you. You should decide when it is necessary to include it or not.


## Community

You can contact us if you have questions using the standard communication channels described in the [developer's site](http://developers-forum.mercadolibre.com/)

## I want to contribute!

That is great! Just fork the project in github. Create a topic branch, write some code, and add some tests for your new code.

To run the tests run ```make test```.

Thanks for helping!

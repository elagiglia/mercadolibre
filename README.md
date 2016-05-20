# MercadoLibre's Golang SDK [DRAFT - WORK IN PROGRESS]

This is the official GO SDK for MercadoLibre's Platform.

## How do I install it?

You can download the latest build at:
    https://github.com/mercadolibre/go-sdk/archive/master.zip

How do I install it using go:

Just add to your pom the following repository

```xml
go get github.com/mercadolibre/sdk
```

And that's it!

## How do I start using it?

The first thing to do is to instance a ```Client``` object. You'll need to give a ```clientId``` and a ```clientSecret```. You can obtain both after creating your own application. For more information on this please read: [creating an application](http://developers.mercadolibre.com/application-manager/)

```go
Client{clientId:123456, clientSecret:"client secret", apiUrl:"https://localhost:3000"}
```
With this instance you can start working on MercadoLibre's APIs.

There are some design considerations worth to mention.
1. This SDK is just a thin layer on top of an http client to handle all the OAuth WebServer flow for you.
2. There is no JSON parsing. This is left to you. But this SDK will include [gson](http://code.google.com/p/google-gson/) library for internal usage.

## How do I redirect users to authorize my application?

This is a 2 step process.

First get the link to redirect the user. This is very easy! Just:

```go

```

This will give you the url to redirect the user. You need to specify a callback url which will be the one that the user will redirected after a successfull authrization process.

Once the user is redirected to your callback url, you'll receive in the query string, a parameter named ```code```. You'll need this for the second part of the process.

```go

```

This will get an ```accessToken``` and a ```refreshToken``` (is case your application has the ```offline_access```) for your application and your user.

At this stage your are ready to make call to the API on behalf of the user.

## Making GET calls

```GO
client := sdk.NewClient()
client.Get("/items")

```

## Making POST calls

```GO
client := sdk.NewClient()
client.Post("/items")

```
## Making PUT calls

```GO

```
## Making DELETE calls

```GO
client := sdk.NewClient()
client.Delete("/items/123")
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
# PropelAuth Go SDK

<p align="center">
  <a href="https://www.propelauth.com/?utm_campaign=github-go" target="_blank" align="center">
    <img src="https://propelauth-logos.s3.us-west-2.amazonaws.com/logo-only.png" width="100">
  </a>
</p>


A Go library for managing authentication, backed by [PropelAuth](https://www.propelauth.com/?utm_campaign=github-go). 

[PropelAuth](https://www.propelauth.com?ref=github) makes it easy to add authentication and authorization to your B2B/multi-tenant application.

Your frontend gets a beautiful, safe, and customizable login screen. Your backend gets easy authorization with just a few lines of code. You get an easy-to-use dashboard to config and manage everything.

## Documentation

- [Full reference for this library (coming soon)](https://docs.propelauth.com/reference/backend-apis/go)
- [Getting started guides for PropelAuth](https://docs.propelauth.com/)

## Creating your first program with the PropelAuth Go SDK

### Initial Steps

1. Sign up for a free account at [Propelauth](https://auth.propelauth.com).
2. Create a new project.
3. Go to the **Backend Integrations** section in the dashboard and note the Auth URL and API Key.


### Install the PropelAuth Go SDK

Make sure you have Go version 1.20.

```sh
go get github.com/propelauth/propelauth-go-sdk
```

### Your sample program


```go
package main

import (
    "os"
    "fmt"
    "github.com/propelauth/propelauth-go-sdk"
)

func main() {
    // you can get these variables from the Backend Integrations section on your dashboard
    apiKey := os.Getenv("PROPELAUTH_API_KEY")
    authUrl := os.Getenv("PROPELAUTH_AUTH_URL")
    
	client, err := InitBaseAuth(authUrl, apiKey, nil)
	if err != nil {
		panic(err)
	}

    // create a new user
    
    newUser = CreateUserParams{
        Email: "tanis@solace.com"
    }

    createdUser, err := client.CreateUser(newUser)
    if err != nil {
        panic(err)
    }

    fmt.Println("Create a user with the ID " + createdUser.UserID)

    // fetch the user we just created

    fetchedUser, err = client.FetchUserMetadataByUserId(createdUser.UserID)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found the user we just created %#v\n", fetchedUser)
}
```

## License

The PropelAuth Go SDK is released under the [MIT license](LICENSE).

## Questions?

Feel free to reach out at support@propelauth.com


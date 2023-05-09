[![Go Report Card](https://goreportcard.com/badge/github.com/propelauth/propelauth-go)](https://goreportcard.com/report/github.com/propelauth/propelauth-go)

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

Make a new directory, and initialize your program.

```sh
mkdir propelauth-example
cd propelauth-example
go mod init propelauth-example
go get github.com/propelauth/propelauth-go
```

Create a file that looks like this.

```go
package main

import (
    "os"
    "fmt"
    propelauth "github.com/propelauth/propelauth-go/pkg"
    models "github.com/propelauth/propelauth-go/pkg/models"
)

func main() {
    // initialize the client

    // (you can get these variables from the Backend Integrations section on your dashboard)
    apiKey := os.Getenv("PROPELAUTH_API_KEY")
    authUrl := os.Getenv("PROPELAUTH_AUTH_URL")

    client, err := propelauth.InitBaseAuth(authUrl, apiKey, nil)
    if err != nil {
        panic(err)
    }

    // see how many users we have now

    queryParams := models.UserQueryParams{}

    users, err := client.FetchUsersByQuery(queryParams)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d users\n", users.TotalUsers)

    // create a new user

    newUser := models.CreateUserParams{
        Email: "tanis@solace.com",
    }

    createdUser, err := client.CreateUser(newUser)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created a user with the ID %#v\n", createdUser.UserID)

    // fetch the user we just created

    fetchedUser, err := client.FetchUserMetadataByUserId(createdUser.UserID, false)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found the user we just created %#v\n", fetchedUser)
}
```

Run it.
    
```sh
go run .
```

## License

The PropelAuth Go SDK is released under the [MIT license](LICENSE).

## Questions?

Feel free to reach out at support@propelauth.com


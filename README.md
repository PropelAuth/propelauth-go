# PropelAuth Go SDK

<p align="center">
  <a href="https://www.propelauth.com/?utm_campaign=github-go" target="_blank" align="center">
    <img src="https://propelauth-logos.s3.us-west-2.amazonaws.com/logo-only.png" width="100">
  </a>
</p>


A go library for managing authentication, backed by [PropelAuth](https://www.propelauth.com/?utm_campaign=github-go). 

[PropelAuth](https://www.propelauth.com?ref=github) makes it easy to add authentication and authorization to your B2B/multi-tenant application.

Your frontend gets a beautiful, safe, and customizable login screen. Your backend gets easy authorization with just a few lines of code. You get an easy-to-use dashboard to config and manage everything.

## Documentation

- [Full reference this library.](https://docs.propelauth.com/reference/backend-apis/go)
- [Getting started guides for PropelAuth.](https://docs.propelauth.com/)

## Installation

```sh
go get github.com/propelauth/propelauth-go-sdk
```

## SDk Usage Requirements

1. Sign up for a free account at [Propelauth](https://propelauth.com).
2. Create a new application.
3. Note the Auth URL and API Key on the Backend Integrations section on your dashboard.

## Sample Program

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
    
	client, err := NewClient(authUrl, apiKey, nil)
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

The PropelAuth Go SDK is released under the [MIT license](LICENSE.md).

## Questions?

Feel free to reach out at support@propelauth.com

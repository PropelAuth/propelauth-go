[![Go Report Card](https://goreportcard.com/badge/github.com/propelauth/propelauth-go)](https://goreportcard.com/report/github.com/propelauth/propelauth-go)

[![Go Reference](https://pkg.go.dev/badge/github.com/propelauth/propelauth-go.svg)](https://pkg.go.dev/github.com/propelauth/propelauth-go)

<p align="center">
  <a href="https://www.propelauth.com?ref=github" target="_blank" align="center">
    <img src="https://www.propelauth.com/imgs/lockup.svg" width="200">
  </a>
</p>

# PropelAuth Go SDK


A Go library for managing authentication, backed by [PropelAuth](https://www.propelauth.com/?utm_campaign=github-go). 

[PropelAuth](https://www.propelauth.com/?utm_campaign=github-go) makes it easy to add authentication and authorization to your B2B/multi-tenant application.

Your frontend gets a beautiful, safe, and customizable login screen. Your backend gets easy authorization with just a few lines of code. You get an easy-to-use dashboard to config and manage everything.

## Documentation

- [Full reference for this library](https://docs.propelauth.com/reference/backend-apis/go)
- [Getting started guides for PropelAuth](https://docs.propelauth.com/)

## Installation

```shell
go get github.com/propelauth/propelauth-go
```


## Initialize

To initialize the library, you call `propelauth.InitBaseAuth` with the configuration for your application:

```go
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
    // ...
}
```

This will fetch the information needed to validate access tokens. In a serverless environment, you may want to skip this one-time fetch,
and you can do so by passing in the `TokenVerificationMetadataInput` object:

```go
client, err := propelauth.InitBaseAuth(authUrl, apiKey, &propelauth.TokenVerificationMetadataInput{
    // (you can get these variables from the Backend Integrations section on your dashboard)
    VerifierKey: os.Getenv("PROPELAUTH_VERIFIER_KEY"),
    Issuer: os.Getenv("PROPELAUTH_ISSUER"),
})
```


## Protect API Routes

After initializing auth, you can verify access tokens by passing in the Authorization header (formatted `Bearer TOKEN`), see [User](https://docs.propelauth.com/reference/backend-apis/go#user) for more information:

```go
user, err := client.GetUser(r.Header.Get("Authorization"))
if err != nil {
    w.WriteHeader(401)
    return
}
```

Here’s an example where we create an auth middleware that will protect a route and set the user on the request context:

```go
func requireUser(client *propelauth.Client, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := client.GetUser(r.Header.Get("Authorization"))
		if err != nil {
				w.WriteHeader(401)
				return
		}
		requestContext := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(requestContext))
	})
}
```

which can then be used like this:

```go
func whoami(w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*models.UserFromToken)
	json.NewEncoder(w).Encode(user)
}
// ...
http.Handle("/api/whoami", requireUser(client, whoami))
```

## Authorization / Organizations

You can also verify which organizations the user is in, and which roles and permissions they have, with the `GetOrgMemberInfo` function on the [user](https://docs.propelauth.com/reference/backend-apis/go#user) object.

### Check Org Membership

Verify that the request was made by a valid user **and** that the user is a member of the specified organization. This can be done using the [User](https://docs.propelauth.com/reference/backend-apis/go#user) object.

```go
orgMemberInfo := user.GetOrgMemberInfo(orgId)
if orgMemberInfo == nil {
        w.WriteHeader(403)
        return
}
```

### Check Org Membership and Role

Similar to checking org membership, but will also verify that the user has a specific Role in the organization. This can be done using the [OrgMemberInfo](https://docs.propelauth.com/reference/backend-apis/go#org-member-info) object.

A user has a Role within an organization. By default, the available roles are Owner, Admin, or Member, but these can be configured. These roles are also hierarchical, so Owner > Admin > Member.

```go
// Assuming a Role structure of Owner => Admin => Member

orgMemberInfo := user.GetOrgMemberInfo(orgId)
if orgMemberInfo == nil {
        w.WriteHeader(403)
        return
}
if !orgMemberInfo.IsRole("Admin") {
        w.WriteHeader(403)
        return
}
```

### Check Org Membership and Permission

Similar to checking org membership, but will also verify that the user has the specified permission in the organization. This can be done using the [OrgMemberInfo](https://docs.propelauth.com/reference/backend-apis/go#org-member-info) object.

Permissions are arbitrary strings associated with a role. For example, `can_view_billing`, `ProductA::CanCreate`, and `ReadOnly` are all valid permissions.
You can create these permissions in the PropelAuth dashboard.

```go
orgMemberInfo := user.GetOrgMemberInfo(orgId)
if orgMemberInfo == nil {
        w.WriteHeader(403)
        return
}
if !orgMemberInfo.HasPermission("can_view_billing") {
        w.WriteHeader(403)
        return
}
```

## Calling Backend APIs

You can also use the library to call the PropelAuth APIs directly, allowing you to fetch users, create orgs, and a lot more.

```go
client, err := propelauth.InitBaseAuth(authUrl, apiKey, nil)

response, err := client.CreateMagicLink(models.CreateMagicLinkParams{
    Email: "test@example.com",
})
```

See the [API Reference](https://docs.propelauth.com/reference) for more information.

## License

The PropelAuth Go SDK is released under the [MIT license](LICENSE).

## Questions?

Feel free to reach out at support@propelauth.com. We like answering questions!

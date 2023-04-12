package client

import (
	"time"
	//jwt "github.com/golang-jwt/jwt/v5"
)

type AccessToken struct {
	AuthId    string
	Key       string
	Username  string
	ValidFrom time.Time
	Lifetime  time.Duration
	Uid       string
	Grants    Grants `json:"grants"`
}

type Grants struct {
}

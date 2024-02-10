# Nordnet

[![Build Status](https://travis-ci.org/0dayfall/nordnet.svg?branch=master)](https://travis-ci.org/0dayfall/nordnet)
[![GoDoc](https://godoc.org/github.com/0dayfall/nordnet?status.svg)](http://godoc.org/github.com/0dayfall/nordnet)
[![Go Report Card](https://goreportcard.com/badge/github.com/0dayfall/nordnet)](https://goreportcard.com/report/github.com/0dayfall/nordnet)

Go implementation of the Nordnet External API.

https://api.test.nordnet.se/api-docs/index.html


## Installation

`go get github.com/0dayfall/nordnet`

## Usage


### REST API Client

```go
package main

import (
	"fmt"
	"github.com/0dayfall/nordnet/api"
	"github.com/0dayfall/nordnet/util"
)

var (
	pemData = []byte(`-----BEGIN PUBLIC KEY-----`)
	user    = []byte(`...`)
	pass    = []byte(`...`)
)

func main() {
	cred, _ := util.GenerateCredentials(user, pass, pemData)
	client := api.NewAPIClient(cred)
	client.Login()

	fmt.Println(client.Accounts())
}
```

To use Nordnet test credentials, try `client := api.NewAPITestClient(cred)`.

### Feed Client

```go
package main

import (
	"fmt"
	"github.com/0dayfall/nordnet/feed"
)

var (
	sessionKey = "..."
	address    = "..."
)

func main() {
	feed, _ := feed.NewPrivateFeed(address)
	feed.Login(sessionKey, nil)

	msgChan := make(chan *PrivateMsg)
	errChan := make(chan error)
	feed.Dispatch(msgChan, errChan)

	for _, msg := range msgChan {
		fmt.Println(msg)
	}
}
```

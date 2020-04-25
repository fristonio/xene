package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/fristonio/xene/pkg/apiserver/client"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
)

func prettyPrintJSON(jsonStr string) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println(jsonStr)
		return
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(jsonStr)
	} else {
		fmt.Println(string(b))
	}
}

func getClientAndAuth() (*client.XeneAPIServer, runtime.ClientAuthInfoWriter) {
	addr, err := url.Parse(config.APIServerAddr)
	if err != nil {
		log.Fatalf("invalid api address URL: %s", err)
	}
	client := client.New(httptransport.New(addr.Host, "", nil), strfmt.Default)
	bearerTokenAuth := httptransport.BearerToken(config.AuthToken)

	return client, bearerTokenAuth
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type permission_url struct {
	client_id     string
	redirect_uri  string
	response_type string
	scope         string
}

type token_authorization_response struct {
	refresh_token                   string
	access_token                    string
	access_token_expiration_seconds string
}

func Generate_permission_url(client_id string) string {
	params := permission_url{
		client_id:     url.QueryEscape(client_id),
		redirect_uri:  url.QueryEscape("urn:ietf:wg:oauth:2.0:oob"),
		response_type: url.QueryEscape("code"),
		scope:         url.QueryEscape("https://mail.google.com/"),
	}

	url := "https://accounts.google.com/o/oauth2/auth?" + "client_id=" + params.client_id + "&redirect_uri=" + params.redirect_uri + "&response_type=" + params.response_type + "&scope=" + params.scope

	return (url)
}

func Generate_token_authorization(client_id string, client_secret string, authorization_code string) token_authorization_response {
	values := map[string]string{"client_id": client_id, "client_secret": client_secret, "code": authorization_code, "redirect_uri": "urn:ietf:wg:oauth:2.0:oob", "grant_type": "authorization_code"}

	json_data, err := json.Marshal(values)
	if err != nil {
		fmt.Println(err)
	}

	respo, err := http.Post("https://accounts.google.com/o/oauth2/token", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err)
	}

	resp := map[string]string{}
	json.NewDecoder(respo.Body).Decode(&resp)

	res := token_authorization_response{
		refresh_token:                   "",
		access_token:                    "",
		access_token_expiration_seconds: "",
	}

	if resp["expires_in"] == "" {
		res = token_authorization_response{
			refresh_token:                   resp["refresh_token"],
			access_token:                    resp["access_token"],
			access_token_expiration_seconds: "3600",
		}
	} else {
		res = token_authorization_response{
			refresh_token:                   resp["refresh_token"],
			access_token:                    resp["access_token"],
			access_token_expiration_seconds: resp["expires_in"],
		}
	}

	return (res)
}

package main

import (
	"net/url"
)

type permission_url struct {
	client_id     string
	redirect_uri  string
	response_type string
	scope         string
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

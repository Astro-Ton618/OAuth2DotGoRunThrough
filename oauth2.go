package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-sasl"
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

func Generate_token_authorization(client_id string, client_secret string, authorization_code string) (token_authorization_response, error) {
	res := token_authorization_response{
		refresh_token:                   "",
		access_token:                    "",
		access_token_expiration_seconds: "",
	}

	values := map[string]string{"client_id": client_id, "client_secret": client_secret, "code": authorization_code, "redirect_uri": "urn:ietf:wg:oauth:2.0:oob", "grant_type": "authorization_code"}

	json_data, err := json.Marshal(values)
	if err != nil {
		return res, err
	}

	respo, err := http.Post("https://accounts.google.com/o/oauth2/token", "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return res, err
	}

	resp := map[string]string{}
	json.NewDecoder(respo.Body).Decode(&resp)

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

	return res, nil
}

func Imap_authentication(email string, access_token string) error {
	c, err_tls := client.DialTLS("imap.gmail.com:993", nil)
	if err_tls != nil {
		return (err_tls)
	}

	c.Authenticate(sasl.NewOAuthBearerClient(&sasl.OAuthBearerOptions{
		Username: email,
		Token:    access_token,
		Host:     "imap.gmail.com",
		Port:     993,
	}))

	mbox, err_select := c.Select("INBOX", false)
	if err_select != nil {
		return (err_select)
	}

	if mbox.Messages == 0 {
		return errors.New("no message in mailbox")
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.Messages)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() error {
		if err_fetch := c.Fetch(seqSet, items, messages); err_fetch != nil {
			return (err_fetch)
		}
		return nil
	}()

	msg := <-messages
	if msg == nil {
		return errors.New("server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		return errors.New("server didn't returned message body")
	}

	mr, err_create_reader := mail.CreateReader(r)
	if err_create_reader != nil {
		return (err_create_reader)
	}

	header := mr.Header
	if date, err_date := header.Date(); err_date == nil {
		fmt.Println("Date:", date)
	}
	if from, err_from := header.AddressList("From"); err_from == nil {
		fmt.Println("From:", from)
	}
	if to, err_to := header.AddressList("To"); err_to == nil {
		fmt.Println("To:", to)
	}
	if subject, err_subject := header.Subject(); err_subject == nil {
		fmt.Println("Subject:", subject)
	}

	for {
		p, err_next_part := mr.NextPart()
		if err_next_part == io.EOF {
			break
		} else if err_next_part != nil {
			return (err_next_part)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			b, _ := ioutil.ReadAll(p.Body)
			fmt.Println("Got text: " + string(b))
		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			fmt.Println("Got attachment: " + filename)
		}
	}

	c.Logout()
	return (nil)
}

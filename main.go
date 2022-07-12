package main

import "fmt"

func main() {
	client_id := ""
	client_secret := ""

	fmt.Println(Generate_permission_url(client_id))

	res_auth, err_token_auth := Generate_token_authorization(client_id, client_secret, "browser_authorization_code")
	if err_token_auth != nil {
		fmt.Println(err_token_auth)
	}
	fmt.Println(res_auth)

	res_refresh, err_token_refresh := Refresh_token(client_id, client_secret, "refresh_token")
	if err_token_refresh != nil {
		fmt.Println(err_token_refresh)
	}
	fmt.Println(res_refresh)

	err_imap := Get_all_email_in_inbox("email", "access_token")
	if err_imap != nil {
		fmt.Println(err_imap)
	}
}

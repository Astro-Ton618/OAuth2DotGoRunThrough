package main

import "fmt"

func main() {
	client_id := ""
	client_secret := ""
	fmt.Println(Generate_permission_url(client_id))
	res, err_token := Generate_token_authorization(client_id, client_secret, "browser_authorization_code")
	if err_token != nil {
		fmt.Println(err_token)
	}
	fmt.Println(res)
	err_imap := Get_all_email_in_inbox("", "")
	if err_imap != nil {
		fmt.Println(err_imap)
	}
}

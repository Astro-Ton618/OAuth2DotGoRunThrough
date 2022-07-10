package main

import "fmt"

func main() {
	client_id := ""
	client_secret := ""
	fmt.Println(Generate_permission_url(client_id))
	fmt.Println(Generate_token_authorization(client_id, client_secret, "browser_authorization_code"))
}

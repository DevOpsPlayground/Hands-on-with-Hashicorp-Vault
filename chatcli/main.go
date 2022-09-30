package main

import(
	"net/http"
	"fmt"
	"log"
	"net/url"
	"io"
	"encoding/json"
)
func failOnError(err error, msg string){
	if err != nil{
		log.Panicf("%s: %s", msg, err)
	}
}

type userpass struct{
	Auth auth `json:"auth"`
}
type auth struct{
	ClientToken string `json:"client_token"`
}
type rabbit struct{
	Data data `json:"data"`
}
type data struct{
	Username string `json:"username"`
	Password string `json:"password"`
}
func main(){
	vaulturl := "http://localhost:8200"
	username := "panda"
	resp, err := http.PostForm(
		fmt.Sprintf("%s/v1/auth/userpass/login/%s",vaulturl, username),
	 	url.Values{"password":{"pass"}})
	failOnError(err, "failed to log in to vault")
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	failOnError(err, "failed to read body")
	fmt.Print(string(body))
	var up userpass
	err = json.Unmarshal(body,&up)
	failOnError(err, "failed to Unmarshal body")
	fmt.Println(up)

	client := &http.Client{}
	req, err := http.NewRequest("GET",fmt.Sprintf("%s/v1/rabbitmq/creds/chat", vaulturl),nil)
	failOnError(err, "failed make req")
	req.Header.Set("X-Vault-Token",up.Auth.ClientToken)
	resp, err = client.Do(req)
	failOnError(err, "failed to send request")
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	failOnError(err, "failed to read body")
	fmt.Print(string(body))
	var rab rabbit
	err = json.Unmarshal(body,&rab)
	failOnError(err, "failed to Unmarshal body")
	fmt.Println(rab)
	runRabbit(rab.Data.Username,rab.Data.Password)
}
//{"request_id":"f086b767-3df9-6203-821c-83dbf71e692e","lease_id":"rabbitmq/creds/chat/uKWmQlixd7ZM63oNAOGW0mJf","renewable":true,"lease_duration":2764800,"data":{"password":"ru7ZLce5kDf7pu6mvEeXVE3SLCxsircwqvWe","username":"userpass-byteford-d0964d28-7606-60ec-75a2-4a962d1f59d2"},"wrap_info":null,"warnings":null,"auth":null}
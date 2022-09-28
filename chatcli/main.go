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
func main(){
	vaulturl := "http://localhost:8200"
	username := "byteford"
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
	fmt.Print(up)
	runRabbit()
}
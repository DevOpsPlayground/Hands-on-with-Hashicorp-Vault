package main

import(
	"log"
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
	
	//code to get go token

	//code to log in to rabbit

	
}

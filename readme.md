

# Start vault and rabbitmq containers

The first thing we need to do it start vault and RabbitMQ, to do this we are just going to use a containerised versions of the apps running in dev mode.

> :warning: **Dont do this in PROD**: This is a quick lab to walk though how vault works, non of this is prod ready

To start the containers run docker compose:
`docker compose up`

# Go login to vault

Now that vault is up time to get familiar with the UI, even though we can (and will) be doing everything with API calls the UI is a nice way of double checking everything works.

UI time
# Set up auth method - userpass

We have seen the UI but its not got a lot going on at the moment, lets set up an auth method so our users can log in.

## Create auth method
Run the following curl request. This will make a userpass authmethod at the path `/userpass`
``` bash
curl --request POST  http://127.0.0.1:8200/v1/sys/auth/userpass \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "type": "userpass"
}'

```

> :exclamation: depending on the auth method used it refers to how you are login diffrently. As we are using userpass I will be talking about user, but could be replces with role if talking about a diffrent auth method

## Create policy
Vault uses policies to let user do things. Before we create a user let create the policy that will let it read the rabbitMQ login details

``` bash
curl --request PUT http://127.0.0.1:8200/v1/sys/policies/acl/rabbitmq \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "policy": "path \"rabbitmq/creds/chat\" {\n capabilities = [\"read\"]\n}"
}'
```

## Create entity
While we could just create a user, we want to alow our app to be easly extended in the future so we are going to create a identity, this will allow us to change (or add) the auth method in the future with out having to change the core functionality. 

The below curl request makes a identity called panda with the policy we attached earlyer. If you wanted to make another user you would just have to rerun the folowing 3 curl commands changeing where you see panda to something else.
``` bash
curl --request POST http://127.0.0.1:8200/v1/identity/entity \
   --header "X-Vault-Token: root" \
   --data-raw '{
  "name": "panda",
  "metadata": {
    "organization": "Playground"
  },
  "policies": ["rabbitmq"]
}'
```
## Create user
Now we have all the ground work set up, time to make the user we can log in with.

Below is the curl command to make a user with the name panda and the password pass. 


``` bash
curl --request POST http://127.0.0.1:8200/v1/auth/userpass/users/panda \
   --header "X-Vault-Token: root" \
   --data-raw '{
    "password": "pass"
}'
```
## Link user to identity
If you now try and log in with the created user, you will find it cant do anything. This is because the identity has the permissions and we havent linked them together yet.  

In the below command you will need to get the canonical_id from the user you just created, and the Mount_accessor we created at the start, if you scrol up you might be able to see them or you can find from the UI:
<pictures showing stuff>
``` bash
curl --request POST http://127.0.0.1:8200/v1/identity/entity-alias \
   --header "X-Vault-Token: root" \
   --data-raw '{
  "name": "panda"
  "canonical_id": "<User canonical id>",
  "mount_accessor": "<Auth method accessor>"
}'
```

# Add code to go to get the auth token

Now that we have created a user, let add the code to allow us to log in.

Open up `./chatcli/main.go`

Replace the import block at the top with
``` go
	import(
	"log"
	"net/http"
	"fmt"
	"net/url"
	"io"
	"encoding/json"
)
```
Then under the comment `//code to get go token` add the following
``` go 
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
```


# Create secrets engine
## create secret mount
``` bash
curl --request POST http://127.0.0.1:8200/v1/sys/mounts/rabbitmq \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "type": "rabbitmq"
}'
```
## connect vault and rabbit
``` bash
curl --request POST 'http://127.0.0.1:8200/v1/rabbitmq/config/connection' \
--header 'X-Vault-Token: root' \
--data-raw '{
    "connection_uri": "http://rabbit:15672",
    "username": "guest",
    "password": "guest"
}'
```
## create rabbit role
``` bash
curl --request POST http://127.0.0.1:8200/v1/rabbitmq/roles/chat \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "vhosts": "{\"/\": {\"configure\":\"^(chat|amq.gen.*)$\", \"write\":\"^(chat|amq.gen.*)$\", \"read\": \"^(chat|amq.gen.*)$\"}}",
  "vhost_topics": "{\"/\": {\"chat\": {\"write\":\".*\", \"read\": \".*\"}}}"
}
'
```
# add code to get rabbitMQ details
``` go
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
```
# start vault and rabbitmq containers
`docker compose up`

# go login to vault
UI time
# set up auth method - userpass
## create auth method
``` bash
curl --request POST  http://127.0.0.1:8200/v1/sys/auth/userpass \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "type": "userpass"
}'

```
## create policy
``` bash
curl --request PUT http://127.0.0.1:8200/v1/sys/policies/acl/rabbitmq \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "policy": "path \"rabbitmq/creds/chat\" {\n capabilities = [\"read\"]\n}"
}'
```

## create entity
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
## create user
``` bash
curl --request POST http://127.0.0.1:8200/v1/auth/userpass/users/panda \
   --header "X-Vault-Token: root" \
   --data-raw '{
    "password": "pass"
}'
```
## link user to identity
``` bash
curl --request POST http://127.0.0.1:8200/v1/identity/entity-alias \
   --header "X-Vault-Token: root" \
   --data-raw '{
  "name": "panda",
  "canonical_id": "2032a31d-84d3-cd64-0a6e-fc58e5919dda",
  "mount_accessor": "auth_userpass_1d260e0c"
}'
```

# add code to go to get the auth token

Add imports
``` go
	"net/http"
	"fmt"
	"net/url"
	"io"
	"encoding/json"
```

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
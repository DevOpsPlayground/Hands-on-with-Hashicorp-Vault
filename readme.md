# start vault and rabbitmq containers
`docker compose up`

# go login to vault
UI time
# set up auth method - userpass
## create auth method
```
curl --request POST  http://127.0.0.1:8200/v1/sys/auth/userpass \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "type": "userpass"
}'

```
## create policy
```
curl --request PUT http://127.0.0.1:8200/v1/sys/policies/acl/rabbitmq \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "policy": "path \"rabbitmq/creds/chat\" {\n capabilities = [\"read\"]\n}"
}'
```

## create entity
```
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
```
curl --request POST http://127.0.0.1:8200/v1/auth/userpass/users/panda \
   --header "X-Vault-Token: root" \
   --data-raw '{
    "password": "pass"
}'
```
## link user to identity
```
curl --request POST http://127.0.0.1:8200/v1/identity/entity-alias \
   --header "X-Vault-Token: root" \
   --data-raw '{
  "name": "panda",
  "canonical_id": "2032a31d-84d3-cd64-0a6e-fc58e5919dda",
  "mount_accessor": "auth_userpass_1d260e0c"
}'
```

# Create secrets engine
## create secret mount
```
curl --request POST http://127.0.0.1:8200/v1/sys/mounts/rabbitmq \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "type": "rabbitmq"
}'
```
## connect vault and rabbit
```
curl --request POST 'http://127.0.0.1:8200/v1/rabbitmq/config/connection' \
--header 'X-Vault-Token: root' \
--data-raw '{
    "connection_uri": "http://rabbit:15672",
    "username": "guest",
    "password": "guest"
}'
```
## create rabbit role
```
curl --request POST http://127.0.0.1:8200/v1/rabbitmq/roles/chat \
    --header "X-Vault-Token: root" \
    --data-raw '{
  "vhosts": "{\"/\": {\"configure\":\"^(chat|amq.gen.*)$\", \"write\":\"^(chat|amq.gen.*)$\", \"read\": \"^(chat|amq.gen.*)$\"}}",
  "vhost_topics": "{\"/\": {\"chat\": {\"write\":\".*\", \"read\": \".*\"}}}"
}
'
```
# go use rabbit
```
cd chatcli
go run .
```
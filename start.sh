docker compose up

cd terraform

terraform init 
terraform apply

export VAULT_ADDR="http://localhost:8200"

curl --request POST 'http://127.0.0.1:8200/v1/rabbitmq/config/connection' \
--header 'X-Vault-Token: root' \
--data-raw '{
    "connection_uri": "http://rabbit:15672",
    "username": "guest",
    "password": "guest"
}'
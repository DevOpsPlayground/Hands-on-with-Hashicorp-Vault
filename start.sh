vault server --dev -dev-root-token-id="root" &

cd terraform

terraform init 
terraform apply

docker run -d --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.10-management

export VAULT_ADDR="http://localhost:8200"
vault write rabbitmq/config/connection \
    connection_uri="http://localhost:15672" \
    username="guest" \
    password="guest"

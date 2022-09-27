vault server --dev -dev-root-token-id="root" &

cd terraform

terraform init 
terraform apply

docker run -d --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.10-management

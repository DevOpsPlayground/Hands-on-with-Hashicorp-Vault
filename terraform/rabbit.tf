resource "vault_rabbitmq_secret_backend" "rabbitmq" {
  connection_uri  = "http://localhost:5672"
  username        = "guest"
  password        = "guest"
  verify_connection = false
}

resource "vault_rabbitmq_secret_backend_role" "role" {
  backend = vault_rabbitmq_secret_backend.rabbitmq.path
  name    = "chat"

  tags = "tag1,tag2"

  vhost {
    host = "/"
    configure = "test"
    read = "^(chat|test)$"
    write = "^(chat|test)$"
  }


  vhost_topic {
    vhost {
      topic = "chat"
      read = ".*"
      write = ".*"
    }

    host = "/"
  }
}
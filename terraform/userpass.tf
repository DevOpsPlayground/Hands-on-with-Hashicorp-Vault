resource "vault_auth_backend" "userpass" {
  type = "userpass"
}


resource "vault_generic_endpoint" "u1" {
  for_each = {byteford:{password:"pass"}}
  depends_on           = [vault_auth_backend.userpass]
  path                 = "auth/userpass/users/${each.key}"
  ignore_absent_fields = true

  data_json = <<EOT
{
  "password": "${each.value.password}"
}
EOT
}

resource "vault_identity_entity" "u1" {
  name      = "byteford"
  policies  = ["rabbitmq"]
  metadata  = {
  }
}
resource "vault_identity_entity_alias" "test" {
  name            = "byteford"
  mount_accessor  = vault_auth_backend.userpass.accessor
  canonical_id    = vault_identity_entity.u1.id
}

resource "vault_policy" "rabbitmq" {
  name = "rabbitmq"

  policy = <<EOT
path "rabbitmq/creds/${vault_rabbitmq_secret_backend_role.role.name}" {
  capabilities = ["read"]
}
EOT
}
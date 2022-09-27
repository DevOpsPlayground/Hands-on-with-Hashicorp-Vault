terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
      version = "3.8.2"
    }
  }
}

provider "vault" {
  # Configuration options
  address = "http://127.0.0.1:8200"
}
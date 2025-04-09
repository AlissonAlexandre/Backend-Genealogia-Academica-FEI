terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.90.0"
    }
  }

  backend "azurerm" {
    resource_group_name  = "genealogia-academica"
    storage_account_name = "genealogiastorage"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"
  }
}

provider "azurerm" {
  features {}

  skip_provider_registration = true
}    

resource "azurerm_resource_group" "main" {
  name     = "genealogia-academica"
  location = "East US"
}

resource "azurerm_container_app_environment" "env" {
  name                = "PROD"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
}

resource "azurerm_container_app" "app" {
  name                         = "genealogia-backend"
  container_app_environment_id = azurerm_container_app_environment.env.id
  resource_group_name          = azurerm_resource_group.main.name
  revision_mode                = "Single"

  template {
    container {
      name   = "genealogia"
      image  = "${var.container_registry_name}.azurecr.io/genealogiaacademicabackend:latest"
      cpu    = 0.5
      memory = "1Gi"

      env {
        name  = "ENV"
        value = "PROD"
      }
      env {
        name  = "NEO4J_URI"
        value = var.NEO4J_URI
      }
      env {
        name  = "NEO4J_USER"
        value = var.NEO4J_USER
      }
      env {
        name  = "NEO4J_PASSWORD"
        value = var.NEO4J_PASSWORD
      }
      env {
        name  = "PORT"
        value = var.PORT
      }
    }

    min_replicas = 0
    max_replicas = 3
    
  }

  ingress {
    external_enabled = true
    target_port      = 8093

    traffic_weight {
      percentage      = 100
      latest_revision = true
    }
  }

  registry {
    server   = "${var.container_registry_name}.azurecr.io"
    username = var.acr_username
    password_secret_name  = var.acr_password
  }
}

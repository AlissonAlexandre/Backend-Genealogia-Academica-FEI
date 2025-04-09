# main.tf
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = var.resource_group
  location = var.location
}

resource "azurerm_container_app_environment" "env" {
  name                = var.environment_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_container_app" "app" {
  name                         = var.container_app_name
  container_app_environment_id = azurerm_container_app_environment.env.id
  resource_group_name          = azurerm_resource_group.rg.name
  revision_mode                = "Single"

  template {
    container {
      name   = "${var.container_app_name}"
      image  = "${var.acr_login_server}/${var.image_name}:latest"

      resources {
        cpu    = 0.5
        memory = "1.0Gi"
      }
      env {
        name  = "ENV"
        value = "prod"
      }
    }
    scale {
      min_replicas = 0
      max_replicas = 3
      rules {
        name = "http-scaling"
        custom {
          type = "http"
          metadata = {
            concurrentRequests = "50"
          }
        }
      }
    }
  }

  ingress {
    external_enabled = true
    target_port      = 8093
    transport        = "auto"
  }

  registry {
    server               = var.acr_login_server
    username             = var.acr_username
    password             = var.acr_password
  }
}

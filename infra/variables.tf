variable "resource_group" {
  description = "Nome do Resource Group"
  type        = string
}

variable "location" {
  description = "Região do Azure"
  type        = string
}

variable "environment_name" {
  description = "Nome do ambiente do Container App"
  type        = string
}

variable "container_app_name" {
  description = "Nome da aplicação"
  type        = string
}

variable "acr_login_server" {
  description = "Servidor do ACR"
  type        = string
}

variable "image_name" {
  description = "Nome da imagem no ACR"
  type        = string
}

variable "acr_username" {
  description = "Usuário do ACR"
  type        = string
  sensitive   = true
}

variable "acr_password" {
  description = "Senha do ACR"
  type        = string
  sensitive   = true
}

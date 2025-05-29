# Переменные для базового примера создания сервера Selectel

variable "project_uuid" {
  description = "UUID проекта Selectel"
  type        = string

  validation {
    condition     = can(regex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", var.project_uuid))
    error_message = "Project UUID должен быть в формате UUID."
  }
}

variable "ssh_public_keys" {
  description = "Список SSH публичных ключей для доступа к серверу"
  type        = list(string)
  default     = []

  validation {
    condition     = length(var.ssh_public_keys) > 0
    error_message = "Необходимо указать хотя бы один SSH ключ."
  }
}

variable "root_password" {
  description = "Пароль root пользователя (опционально, если используются SSH ключи)"
  type        = string
  default     = ""
  sensitive   = true

  validation {
    condition     = var.root_password == "" || length(var.root_password) >= 8
    error_message = "Пароль должен содержать минимум 8 символов."
  }
}

variable "server_name" {
  description = "Имя создаваемого сервера"
  type        = string
  default     = "web-server-basic"

  validation {
    condition     = length(var.server_name) > 0 && length(var.server_name) <= 255
    error_message = "Имя сервера должно содержать от 1 до 255 символов."
  }
}

variable "environment" {
  description = "Окружение (development, staging, production)"
  type        = string
  default     = "production"

  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "Environment должен быть одним из: development, staging, production."
  }
} 
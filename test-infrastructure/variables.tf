# Основные переменные для тестирования

variable "server_name" {
  description = "Имя тестового сервера"
  type        = string
  default     = "terraform-test-server"
  
  validation {
    condition     = length(var.server_name) > 3 && length(var.server_name) <= 50
    error_message = "Имя сервера должно быть от 4 до 50 символов."
  }
}

variable "project_uuid" {
  description = "UUID проекта Selectel"
  type        = string
  sensitive   = true
  
  validation {
    condition     = can(regex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", var.project_uuid))
    error_message = "project_uuid должен быть в формате UUID."
  }
}

variable "service_uuid" {
  description = "UUID услуги (конфигурации сервера)"
  type        = string
  sensitive   = true
  
  validation {
    condition     = can(regex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", var.service_uuid))
    error_message = "service_uuid должен быть в формате UUID."
  }
}

variable "price_plan_uuid" {
  description = "UUID тарифного плана"
  type        = string
  sensitive   = true
  
  validation {
    condition     = can(regex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", var.price_plan_uuid))
    error_message = "price_plan_uuid должен быть в формате UUID."
  }
}

variable "os_template_uuid" {
  description = "UUID шаблона операционной системы"
  type        = string
  sensitive   = true
  
  validation {
    condition     = can(regex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$", var.os_template_uuid))
    error_message = "os_template_uuid должен быть в формате UUID."
  }
}

variable "root_password" {
  description = "Пароль root пользователя"
  type        = string
  sensitive   = true
  default     = null
  
  validation {
    condition = var.root_password == null || (
      length(var.root_password) >= 8 &&
      can(regex("[A-Z]", var.root_password)) &&
      can(regex("[a-z]", var.root_password)) &&
      can(regex("[0-9]", var.root_password))
    )
    error_message = "Пароль должен содержать минимум 8 символов, включая заглавные и строчные буквы, цифры."
  }
}

variable "ssh_keys" {
  description = "Список SSH публичных ключей"
  type        = list(string)
  sensitive   = true
  
  validation {
    condition = alltrue([
      for key in var.ssh_keys : can(regex("^(ssh-rsa|ssh-ed25519|ssh-dss|ecdsa-sha2-nistp256|ecdsa-sha2-nistp384|ecdsa-sha2-nistp521)", key))
    ])
    error_message = "Все SSH ключи должны быть в правильном формате."
  }
}

variable "create_second_server" {
  description = "Создавать ли второй сервер для тестирования множественных ресурсов"
  type        = bool
  default     = false
}

# Переменные для тестирования различных сценариев

variable "test_scenario" {
  description = "Сценарий тестирования"
  type        = string
  default     = "basic"
  
  validation {
    condition     = contains(["basic", "multiple", "advanced"], var.test_scenario)
    error_message = "test_scenario должен быть одним из: basic, multiple, advanced."
  }
}

variable "auto_destroy" {
  description = "Автоматически удалять ресурсы после тестирования"
  type        = bool
  default     = false
} 
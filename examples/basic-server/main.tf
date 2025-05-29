# Базовый пример создания выделенного сервера Selectel
terraform {
  required_version = ">= 1.5"
  required_providers {
    selectel-baremetal = {
      source  = "selectel/selectel-baremetal"
      version = "~> 0.1"
    }
  }
}

# Конфигурация провайдера
provider "selectel-baremetal" {
  # Токен и Project ID будут взяты из переменных окружения:
  # export SELECTEL_TOKEN="your-iam-token"
  # export SELECTEL_PROJECT_ID="your-project-uuid"

  # Опционально можно указать endpoint:
  # endpoint = "https://api.selectel.ru/servers/v2"
}

# Получение списка доступных локаций
data "selectel_baremetal_locations" "all" {}

# Получение списка доступных услуг
data "selectel_baremetal_services" "all" {}

# Получение тарифных планов
data "selectel_baremetal_price_plans" "all" {}

# Получение Ubuntu шаблонов
data "selectel_baremetal_os_templates" "ubuntu" {}

# Локальные переменные для упрощения конфигурации
locals {
  # Выбираем первую доступную локацию (обычно Moscow)
  location_uuid = data.selectel_baremetal_locations.all.locations[0].uuid

  # Выбираем первую доступную услугу
  service_uuid = data.selectel_baremetal_services.all.services[0].uuid

  # Выбираем первый доступный тарифный план
  price_plan_uuid = data.selectel_baremetal_price_plans.all.price_plans[0].uuid

  # Выбираем Ubuntu шаблон
  ubuntu_template = [
    for template in data.selectel_baremetal_os_templates.ubuntu.templates :
    template if can(regex("ubuntu", lower(template.name)))
  ][0]
}

# Создание базового сервера
resource "selectel_baremetal_server" "web" {
  name            = "web-server-basic"
  service_uuid    = local.service_uuid
  location_uuid   = local.location_uuid
  price_plan_uuid = local.price_plan_uuid
  project_uuid    = var.project_uuid

  # Сетевая конфигурация
  network {
    type      = "public"
    bandwidth = 1000
  }

  # Операционная система
  os {
    template_uuid = local.ubuntu_template.uuid
    ssh_keys      = var.ssh_public_keys
    password      = var.root_password
  }

  # Теги для организации ресурсов
  tags = {
    Environment = "production"
    Application = "web"
    Owner       = "terraform"
    Project     = "basic-example"
  }
}

# Переменные
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

# Выводы
output "server_info" {
  description = "Информация о созданном сервере"
  value = {
    uuid         = selectel_baremetal_server.web.uuid
    name         = selectel_baremetal_server.web.name
    status       = selectel_baremetal_server.web.status
    power_status = selectel_baremetal_server.web.power_status
    ip_addresses = selectel_baremetal_server.web.ip_addresses
  }
}

output "location_info" {
  description = "Информация о выбранной локации"
  value = {
    name = data.selectel_baremetal_locations.all.locations[0].name
    uuid = local.location_uuid
  }
}

output "service_info" {
  description = "Информация о выбранной услуге"
  value = {
    name = data.selectel_baremetal_services.all.services[0].name
    uuid = local.service_uuid
  }
}

output "os_template_info" {
  description = "Информация об используемом шаблоне ОС"
  value = {
    name         = local.ubuntu_template.name
    version      = local.ubuntu_template.version
    architecture = local.ubuntu_template.architecture
    uuid         = local.ubuntu_template.uuid
  }
} 
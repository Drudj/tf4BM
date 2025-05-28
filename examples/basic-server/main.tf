terraform {
  required_providers {
    selectel = {
      source  = "selectel/selectel-baremetal"
      version = "~> 0.1"
    }
  }
}

provider "selectel" {
  # Конфигурация будет взята из переменных окружения:
  # SELECTEL_TOKEN
  # SELECTEL_PROJECT_ID
  # SELECTEL_ENDPOINT (опционально)
}

# Получение списка всех локаций
data "selectel_baremetal_locations" "all" {}

# Получение локации по имени
data "selectel_baremetal_location" "msk" {
  name = "Moscow"
}

# Получение списка услуг
data "selectel_baremetal_services" "all" {}

# Получение OS шаблонов
data "selectel_baremetal_os_templates" "ubuntu" {
  type   = "linux"
  family = "ubuntu"
}

# Получение тарифных планов
data "selectel_baremetal_price_plans" "all" {}

# Создание базового сервера
resource "selectel_baremetal_server" "web" {
  name            = "web-server-01"
  service_uuid    = data.selectel_baremetal_services.all.services[0].uuid
  location_uuid   = data.selectel_baremetal_location.msk.uuid
  price_plan_uuid = data.selectel_baremetal_price_plans.all.price_plans[0].uuid
  project_uuid    = var.project_uuid

  network {
    type      = "public"
    bandwidth = 1000
  }

  os {
    template_uuid = data.selectel_baremetal_os_templates.ubuntu.templates[0].uuid
    ssh_keys      = [var.ssh_public_key]
  }

  tags = {
    Environment = "production"
    Service     = "web"
    Terraform   = "true"
  }
}

# Переменные
variable "project_uuid" {
  description = "UUID проекта Selectel"
  type        = string
}

variable "ssh_public_key" {
  description = "SSH публичный ключ для доступа к серверу"
  type        = string
}

# Выводы
output "server_uuid" {
  description = "UUID созданного сервера"
  value       = selectel_baremetal_server.web.uuid
}

output "server_status" {
  description = "Статус сервера"
  value       = selectel_baremetal_server.web.status
}

output "server_ip_addresses" {
  description = "IP адреса сервера"
  value       = selectel_baremetal_server.web.ip_addresses
}

output "all_locations" {
  description = "Список всех доступных локаций"
  value       = data.selectel_baremetal_locations.all.locations
}

output "moscow_location" {
  description = "Информация о локации Moscow"
  value       = data.selectel_baremetal_location.msk
} 
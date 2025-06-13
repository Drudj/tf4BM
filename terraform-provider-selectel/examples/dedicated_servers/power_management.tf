# Пример управления питанием выделенного сервера

terraform {
  required_providers {
    selectel = {
      source = "selectel/selectel"
    }
  }
}

# Настройка провайдера Selectel
provider "selectel" {
  # Основные настройки (через переменные окружения)
  # OS_AUTH_URL
  # OS_REGION_NAME  
  # OS_DOMAIN_NAME
  # OS_USERNAME
  # OS_PASSWORD
  
  # Токен для выделенных серверов
  # SEL_SERVERS_TOKEN
}

# Получение информации о существующем сервере
data "selectel_dedicated_server_v1" "example" {
  id = 12345  # ID вашего сервера
}

# Управление питанием сервера
resource "selectel_dedicated_server_power_v1" "example" {
  server_id = data.selectel_dedicated_server_v1.example.id
  action    = "restart"  # start, stop, restart, power_cycle
  force     = false      # принудительное выполнение для stop/restart
  
  timeouts {
    create = "10m"
    update = "10m"
  }
}

# Переустановка ОС на сервере
resource "selectel_dedicated_server_reinstall_v1" "example" {
  server_id = data.selectel_dedicated_server_v1.example.id
  os_id     = 5  # ID операционной системы для установки
  
  ssh_keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... user@example.com",
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5... admin@example.com"
  ]
  
  preserve_data = false  # сохранить пользовательские данные
  
  timeouts {
    create = "60m"
    update = "60m"
  }
}

# Получение информации о задачах сервера
data "selectel_dedicated_server_tasks_v1" "task_info" {
  task_id = selectel_dedicated_server_power_v1.example.task_id
}

# Вывод информации
output "server_status" {
  description = "Текущий статус сервера"
  value = {
    name   = data.selectel_dedicated_server_v1.example.name
    status = selectel_dedicated_server_power_v1.example.status
    last_action = selectel_dedicated_server_power_v1.example.last_action_at
  }
}

output "reinstall_info" {
  description = "Информация о переустановке ОС"
  value = {
    status = selectel_dedicated_server_reinstall_v1.example.status
    reinstalled_at = selectel_dedicated_server_reinstall_v1.example.reinstalled_at
    os_info = selectel_dedicated_server_reinstall_v1.example.os_info
  }
}

output "task_info" {
  description = "Информация о выполненных задачах"
  value = data.selectel_dedicated_server_tasks_v1.task_info.tasks
} 
# Полный пример использования провайдера выделенных серверов Selectel
# Демонстрирует все возможности провайдера

terraform {
  required_version = ">= 0.13"
  required_providers {
    selectel = {
      source  = "selectel/selectel"
      version = "~> 5.0"
    }
  }
}

# Конфигурация провайдера
provider "selectel" {
  servers_token = var.selectel_servers_token
}

# Переменные
variable "selectel_servers_token" {
  description = "Selectel servers API token"
  type        = string
  sensitive   = true
}

variable "ssh_public_keys" {
  description = "List of SSH public keys for server access"
  type        = list(string)
  default = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7QA... admin@company.com"
  ]
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

# Data sources для получения доступных ресурсов

# Получение списка доступных локаций
data "selectel_dedicated_server_locations_v1" "available_locations" {}

# Получение списка доступных операционных систем
data "selectel_dedicated_server_os_v1" "available_os" {}

# Получение конфигураций для московской локации
data "selectel_dedicated_server_configurations_v1" "moscow_configs" {
  location_id = 1  # Moscow
}

# Локальные переменные для удобства
locals {
  # Выбираем Ubuntu 20.04 LTS
  ubuntu_os_id = [
    for os in data.selectel_dedicated_server_os_v1.available_os.os_list :
    os.id if contains(lower(os.name), "ubuntu") && contains(os.version, "20.04")
  ][0]

  # Выбираем среднюю конфигурацию
  medium_config_id = data.selectel_dedicated_server_configurations_v1.moscow_configs.configurations[1].id

  # Общие теги
  common_tags = [
    var.environment,
    "terraform-managed",
    "selectel-dedicated"
  ]
}

# Основные серверы приложения
resource "selectel_dedicated_server_v1" "app_servers" {
  count = 2

  name        = "app-server-${count.index + 1}"
  location_id = 1  # Moscow
  config_id   = local.medium_config_id
  os_id       = local.ubuntu_os_id

  comment = "Application server ${count.index + 1} for ${var.environment} environment"

  tags = concat(local.common_tags, [
    "app-server",
    "web-tier"
  ])

  ssh_keys = var.ssh_public_keys

  enable_ipmi   = true
  enable_backup = true

  network_config {
    additional_ips  = 1
    private_network = true
  }

  timeouts {
    create = "60m"
    update = "30m"
    delete = "30m"
  }
}

# Сервер базы данных
resource "selectel_dedicated_server_v1" "database_server" {
  name        = "db-server-primary"
  location_id = 1  # Moscow
  config_id   = data.selectel_dedicated_server_configurations_v1.moscow_configs.configurations[2].id  # Более мощная конфигурация
  os_id       = local.ubuntu_os_id

  comment = "Primary database server for ${var.environment}"

  tags = concat(local.common_tags, [
    "database",
    "primary",
    "critical"
  ])

  ssh_keys = var.ssh_public_keys

  enable_ipmi   = true
  enable_backup = true

  network_config {
    additional_ips  = 0
    private_network = true
  }

  timeouts {
    create = "60m"
    delete = "30m"
  }

  lifecycle {
    prevent_destroy = true
  }
}

# Балансировщик нагрузки
resource "selectel_dedicated_server_v1" "load_balancer" {
  name        = "lb-server"
  location_id = 1  # Moscow
  config_id   = data.selectel_dedicated_server_configurations_v1.moscow_configs.configurations[0].id  # Минимальная конфигурация
  os_id       = local.ubuntu_os_id

  comment = "Load balancer for ${var.environment} environment"

  tags = concat(local.common_tags, [
    "load-balancer",
    "frontend"
  ])

  ssh_keys = var.ssh_public_keys

  enable_ipmi   = false
  enable_backup = false

  network_config {
    additional_ips  = 2  # Для failover
    private_network = true
  }

  timeouts {
    create = "60m"
    delete = "30m"
  }
}

# Управление питанием серверов

# Запуск всех серверов приложений
resource "selectel_dedicated_server_power_v1" "start_app_servers" {
  count = length(selectel_dedicated_server_v1.app_servers)

  server_id = selectel_dedicated_server_v1.app_servers[count.index].id
  action    = "start"

  timeouts {
    create = "10m"
  }

  depends_on = [selectel_dedicated_server_v1.app_servers]
}

# Запуск сервера базы данных
resource "selectel_dedicated_server_power_v1" "start_database" {
  server_id = selectel_dedicated_server_v1.database_server.id
  action    = "start"

  timeouts {
    create = "10m"
  }

  depends_on = [selectel_dedicated_server_v1.database_server]
}

# Запуск балансировщика (после серверов приложений)
resource "selectel_dedicated_server_power_v1" "start_load_balancer" {
  server_id = selectel_dedicated_server_v1.load_balancer.id
  action    = "start"

  timeouts {
    create = "10m"
  }

  depends_on = [
    selectel_dedicated_server_power_v1.start_app_servers,
    selectel_dedicated_server_power_v1.start_database
  ]
}

# Мониторинг задач серверов
data "selectel_dedicated_server_tasks_v1" "app_server_tasks" {
  count = length(selectel_dedicated_server_v1.app_servers)

  server_id = selectel_dedicated_server_v1.app_servers[count.index].id

  depends_on = [selectel_dedicated_server_power_v1.start_app_servers]
}

data "selectel_dedicated_server_tasks_v1" "database_tasks" {
  server_id = selectel_dedicated_server_v1.database_server.id

  depends_on = [selectel_dedicated_server_power_v1.start_database]
}

# Получение информации о созданных серверах
data "selectel_dedicated_server_v1" "app_server_info" {
  count = length(selectel_dedicated_server_v1.app_servers)

  id = selectel_dedicated_server_v1.app_servers[count.index].id

  depends_on = [selectel_dedicated_server_power_v1.start_app_servers]
}

# Outputs для использования в других модулях или для отображения информации

output "server_summary" {
  description = "Summary of all created servers"
  value = {
    app_servers = {
      for i, server in selectel_dedicated_server_v1.app_servers :
      server.name => {
        id         = server.id
        status     = server.status
        primary_ip = server.network.primary_ip
        location   = server.location.name
      }
    }
    database_server = {
      id         = selectel_dedicated_server_v1.database_server.id
      name       = selectel_dedicated_server_v1.database_server.name
      status     = selectel_dedicated_server_v1.database_server.status
      primary_ip = selectel_dedicated_server_v1.database_server.network.primary_ip
      location   = selectel_dedicated_server_v1.database_server.location.name
    }
    load_balancer = {
      id         = selectel_dedicated_server_v1.load_balancer.id
      name       = selectel_dedicated_server_v1.load_balancer.name
      status     = selectel_dedicated_server_v1.load_balancer.status
      primary_ip = selectel_dedicated_server_v1.load_balancer.network.primary_ip
      location   = selectel_dedicated_server_v1.load_balancer.location.name
    }
  }
}

output "server_ips" {
  description = "IP addresses of all servers"
  value = {
    app_servers    = [for server in selectel_dedicated_server_v1.app_servers : server.network.primary_ip]
    database       = selectel_dedicated_server_v1.database_server.network.primary_ip
    load_balancer  = selectel_dedicated_server_v1.load_balancer.network.primary_ip
  }
}

output "server_costs" {
  description = "Monthly costs for all servers"
  value = {
    app_servers_total = sum([
      for server in selectel_dedicated_server_v1.app_servers :
      server.price.amount if server.price.period == "monthly"
    ])
    database_cost = selectel_dedicated_server_v1.database_server.price.amount
    load_balancer_cost = selectel_dedicated_server_v1.load_balancer.price.amount
    total_monthly_cost = sum([
      for server in concat(
        selectel_dedicated_server_v1.app_servers,
        [selectel_dedicated_server_v1.database_server],
        [selectel_dedicated_server_v1.load_balancer]
      ) : server.price.amount if server.price.period == "monthly"
    ])
    currency = selectel_dedicated_server_v1.database_server.price.currency
  }
}

output "available_resources" {
  description = "Available locations, OS, and configurations"
  value = {
    locations = {
      for location in data.selectel_dedicated_server_locations_v1.available_locations.locations :
      location.code => {
        id   = location.id
        name = location.name
        city = location.city
      }
    }
    operating_systems = {
      for os in data.selectel_dedicated_server_os_v1.available_os.os_list :
      "${os.name} ${os.version}" => {
        id           = os.id
        architecture = os.architecture
        type         = os.type
      }
    }
    configurations_count = length(data.selectel_dedicated_server_configurations_v1.moscow_configs.configurations)
  }
}

# Пример условного создания ресурсов
variable "create_backup_server" {
  description = "Whether to create a backup database server"
  type        = bool
  default     = false
}

resource "selectel_dedicated_server_v1" "backup_database_server" {
  count = var.create_backup_server ? 1 : 0

  name        = "db-server-backup"
  location_id = 2  # SPB для географического разнесения
  config_id   = selectel_dedicated_server_v1.database_server.config_id
  os_id       = selectel_dedicated_server_v1.database_server.os_id

  comment = "Backup database server for ${var.environment}"

  tags = concat(local.common_tags, [
    "database",
    "backup",
    "disaster-recovery"
  ])

  ssh_keys = var.ssh_public_keys

  enable_ipmi   = true
  enable_backup = true

  network_config {
    additional_ips  = 0
    private_network = true
  }

  timeouts {
    create = "60m"
    delete = "30m"
  }
}

# Переустановка ОС на сервере (пример для maintenance)
variable "reinstall_server_id" {
  description = "ID of server to reinstall (0 to skip)"
  type        = number
  default     = 0
}

variable "new_os_id" {
  description = "New OS ID for reinstallation"
  type        = number
  default     = 0
}

resource "selectel_dedicated_server_reinstall_v1" "server_maintenance" {
  count = var.reinstall_server_id > 0 ? 1 : 0

  server_id     = var.reinstall_server_id
  os_id         = var.new_os_id > 0 ? var.new_os_id : local.ubuntu_os_id
  preserve_data = false

  ssh_keys = var.ssh_public_keys

  timeouts {
    create = "60m"
  }
} 
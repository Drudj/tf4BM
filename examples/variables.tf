variable "selectel_auth_url" {
  description = "URL для аутентификации Selectel"
  type        = string
  default     = "https://cloud.api.selcloud.ru/identity/v3/"
}

variable "selectel_auth_region" {
  description = "Регион для аутентификации Selectel"
  type        = string
  default     = "pool"
}

variable "selectel_domain_name" {
  description = "Имя домена Selectel"
  type        = string
}

variable "selectel_username" {
  description = "Имя пользователя Selectel"
  type        = string
}

variable "selectel_password" {
  description = "Пароль пользователя Selectel"
  type        = string
  sensitive   = true
}

variable "selectel_servers_token" {
  description = "Токен для работы с API выделенных серверов Selectel"
  type        = string
  sensitive   = true
}

variable "server_name" {
  description = "Имя создаваемого сервера"
  type        = string
  default     = "terraform-server"
}

variable "server_config_id" {
  description = "ID конфигурации сервера"
  type        = number
  default     = 23
}

variable "server_location_id" {
  description = "ID локации сервера"
  type        = number
  default     = 2
}

variable "server_os_id" {
  description = "ID операционной системы"
  type        = number
  default     = 1
}

variable "server_root_size" {
  description = "Размер корневого раздела в GB"
  type        = number
  default     = 50
}

variable "server_swap_size" {
  description = "Размер swap раздела в GB"
  type        = number
  default     = 5
}

variable "server_raid_type" {
  description = "Тип RAID"
  type        = string
  default     = "RAID1"
}

variable "server_ssh_keys" {
  description = "Список SSH ключей"
  type        = list(string)
  default     = []
} 
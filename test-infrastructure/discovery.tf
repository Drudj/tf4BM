terraform {
  required_version = ">= 1.0"
  
  required_providers {
    selectel = {
      source  = "selectel/selectel-baremetal"
      version = "~> 0.1"
    }
  }
}

provider "selectel" {
  # Credentials из переменных окружения
}

# Получение всех локаций
data "selectel_baremetal_locations" "discovery" {}

# Получение всех услуг
data "selectel_baremetal_services" "discovery" {}

# Получение всех OS шаблонов
data "selectel_baremetal_os_templates" "discovery" {}

# Получение всех тарифных планов
data "selectel_baremetal_price_plans" "discovery" {}

# Outputs для получения информации
output "all_locations" {
  description = "Все доступные локации"
  value = {
    for loc in data.selectel_baremetal_locations.discovery.locations : loc.name => {
      uuid        = loc.uuid
      city        = loc.city
      country     = loc.country
      description = loc.description
    }
  }
}

output "ubuntu_templates" {
  description = "Доступные Ubuntu шаблоны"
  value = {
    for template in data.selectel_baremetal_os_templates.discovery.templates : 
    template.name => {
      uuid   = template.uuid
      type   = template.type
      family = template.family
    } if template.family == "ubuntu"
  }
}

output "services_by_location" {
  description = "Услуги по локациям"
  value = {
    for service in data.selectel_baremetal_services.discovery.services : 
    service.name => {
      uuid         = service.uuid
      location     = service.location_uuid
      category     = service.category
      description  = service.description
    }
  }
}

# Найти локацию где доступен наш сервис
locals {
  target_service_uuid = "3af31318-f0aa-4341-bf64-df88e2e3d887"
  
  service_location = try([
    for service in data.selectel_baremetal_services.discovery.services : 
    service.location_uuid if service.uuid == local.target_service_uuid
  ][0], null)
  
  location_name = try([
    for loc in data.selectel_baremetal_locations.discovery.locations : 
    loc.name if loc.uuid == local.service_location
  ][0], "Unknown")
}

output "target_service_location" {
  description = "Локация для выбранного сервиса"
  value = {
    service_uuid    = local.target_service_uuid
    location_uuid   = local.service_location
    location_name   = local.location_name
  }
}

# Рекомендуемый Ubuntu шаблон (самый новый)
output "recommended_ubuntu" {
  description = "Рекомендуемый Ubuntu шаблон"
  value = try([
    for template in data.selectel_baremetal_os_templates.discovery.templates : 
    {
      uuid = template.uuid
      name = template.name
    } if template.family == "ubuntu"
  ][0], null)
} 
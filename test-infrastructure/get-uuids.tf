# Файл для получения UUID из data sources
# Используйте этот файл для получения реальных UUID перед основным тестированием

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

# Outputs для копирования UUID в terraform.tfvars
output "discovery_info" {
  description = "Информация для заполнения terraform.tfvars"
  value = {
    
    # Локации
    locations = {
      for loc in data.selectel_baremetal_locations.discovery.locations : loc.name => {
        uuid        = loc.uuid
        city        = loc.city
        country     = loc.country
        description = loc.description
      }
    }
    
    # Услуги (первые 5 для примера)
    services = {
      for idx, service in slice(data.selectel_baremetal_services.discovery.services, 0, min(5, length(data.selectel_baremetal_services.discovery.services))) : 
      service.name => {
        uuid         = service.uuid
        category     = service.category
        location     = service.location_uuid
        description  = service.description
      }
    }
    
    # OS шаблоны Ubuntu
    ubuntu_templates = {
      for template in data.selectel_baremetal_os_templates.discovery.templates : 
      template.name => {
        uuid   = template.uuid
        type   = template.type
        family = template.family
      } if template.family == "ubuntu"
    }
    
    # Тарифные планы (первые 3)
    price_plans = {
      for idx, plan in slice(data.selectel_baremetal_price_plans.discovery.price_plans, 0, min(3, length(data.selectel_baremetal_price_plans.discovery.price_plans))) : 
      plan.name => {
        uuid = plan.uuid
        type = plan.type
      }
    }
  }
}

# Рекомендуемые значения для быстрого старта
output "recommended_values" {
  description = "Рекомендуемые значения для terraform.tfvars"
  value = {
    # Московская локация (если доступна)
    location_uuid = try([
      for loc in data.selectel_baremetal_locations.discovery.locations : loc.uuid 
      if contains(["Moscow", "MSK", "Москва"], loc.name)
    ][0], data.selectel_baremetal_locations.discovery.locations[0].uuid)
    
    # Первая доступная услуга
    service_uuid = length(data.selectel_baremetal_services.discovery.services) > 0 ? 
                   data.selectel_baremetal_services.discovery.services[0].uuid : 
                   "no-services-available"
    
    # Первый Ubuntu шаблон
    os_template_uuid = try([
      for template in data.selectel_baremetal_os_templates.discovery.templates : template.uuid 
      if template.family == "ubuntu"
    ][0], "no-ubuntu-templates-available")
    
    # Первый тарифный план
    price_plan_uuid = length(data.selectel_baremetal_price_plans.discovery.price_plans) > 0 ? 
                      data.selectel_baremetal_price_plans.discovery.price_plans[0].uuid : 
                      "no-price-plans-available"
  }
} 
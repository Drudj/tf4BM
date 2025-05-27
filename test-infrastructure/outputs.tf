# Outputs для мониторинга и проверки результатов

# Информация о data sources
output "available_locations" {
  description = "Список всех доступных локаций"
  value = {
    count     = length(data.selectel_baremetal_locations.all.locations)
    locations = [for loc in data.selectel_baremetal_locations.all.locations : {
      uuid = loc.uuid
      name = loc.name
      city = loc.city
    }]
  }
}

output "moscow_location_info" {
  description = "Информация о локации Moscow"
  value = {
    uuid        = data.selectel_baremetal_location.moscow.uuid
    name        = data.selectel_baremetal_location.moscow.name
    city        = data.selectel_baremetal_location.moscow.city
    country     = data.selectel_baremetal_location.moscow.country
    description = data.selectel_baremetal_location.moscow.description
  }
}

output "available_services_count" {
  description = "Количество доступных услуг"
  value = {
    total_services        = length(data.selectel_baremetal_services.all.services)
    moscow_services_count = length(data.selectel_baremetal_services.moscow_services.services)
  }
}

output "ubuntu_templates_info" {
  description = "Информация о доступных Ubuntu шаблонах"
  value = {
    count     = length(data.selectel_baremetal_os_templates.ubuntu.templates)
    templates = [for template in data.selectel_baremetal_os_templates.ubuntu.templates : {
      uuid = template.uuid
      name = template.name
      type = template.type
    }]
  }
}

output "price_plans_info" {
  description = "Информация о тарифных планах"
  value = {
    count = length(data.selectel_baremetal_price_plans.monthly.price_plans)
    plans = [for plan in data.selectel_baremetal_price_plans.monthly.price_plans : {
      uuid = plan.uuid
      name = plan.name
      type = plan.type
    }]
  }
}

# Информация о созданных серверах
output "test_server_info" {
  description = "Информация о тестовом сервере"
  value = {
    uuid         = selectel_baremetal_server.test_server.uuid
    name         = selectel_baremetal_server.test_server.name
    status       = selectel_baremetal_server.test_server.status
    power_status = selectel_baremetal_server.test_server.power_status
    created_at   = selectel_baremetal_server.test_server.created_at
    updated_at   = selectel_baremetal_server.test_server.updated_at
    tags         = selectel_baremetal_server.test_server.tags
  }
}

output "test_server_network" {
  description = "Сетевая конфигурация тестового сервера"
  value = {
    network_type      = selectel_baremetal_server.test_server.network.type
    network_bandwidth = selectel_baremetal_server.test_server.network.bandwidth
    ip_addresses      = selectel_baremetal_server.test_server.ip_addresses
  }
}

output "test_server_2_info" {
  description = "Информация о втором тестовом сервере (если создан)"
  value = var.create_second_server ? {
    uuid         = selectel_baremetal_server.test_server_2[0].uuid
    name         = selectel_baremetal_server.test_server_2[0].name
    status       = selectel_baremetal_server.test_server_2[0].status
    power_status = selectel_baremetal_server.test_server_2[0].power_status
    ip_addresses = selectel_baremetal_server.test_server_2[0].ip_addresses
  } : null
}

# Сводная информация о тестировании
output "test_summary" {
  description = "Сводка результатов тестирования"
  value = {
    test_scenario        = var.test_scenario
    servers_created      = var.create_second_server ? 2 : 1
    data_sources_tested  = 5
    provider_version     = "0.1.0"
    test_timestamp       = timestamp()
    
    # Проверки работоспособности
    data_sources_working = {
      locations    = length(data.selectel_baremetal_locations.all.locations) > 0
      services     = length(data.selectel_baremetal_services.all.services) > 0
      os_templates = length(data.selectel_baremetal_os_templates.ubuntu.templates) > 0
      price_plans  = length(data.selectel_baremetal_price_plans.monthly.price_plans) > 0
    }
    
    server_creation_successful = selectel_baremetal_server.test_server.uuid != ""
  }
}

# Информация для подключения к серверам
output "connection_info" {
  description = "Информация для подключения к созданным серверам"
  value = {
    primary_server = {
      name         = selectel_baremetal_server.test_server.name
      uuid         = selectel_baremetal_server.test_server.uuid
      ip_addresses = [for ip in selectel_baremetal_server.test_server.ip_addresses : ip.address if ip.type == "public"]
      ssh_command  = length([for ip in selectel_baremetal_server.test_server.ip_addresses : ip.address if ip.type == "public"]) > 0 ? 
                     "ssh root@${[for ip in selectel_baremetal_server.test_server.ip_addresses : ip.address if ip.type == "public"][0]}" : 
                     "No public IP available"
    }
    
    secondary_server = var.create_second_server ? {
      name         = selectel_baremetal_server.test_server_2[0].name
      uuid         = selectel_baremetal_server.test_server_2[0].uuid
      ip_addresses = [for ip in selectel_baremetal_server.test_server_2[0].ip_addresses : ip.address if ip.type == "public"]
      ssh_command  = length([for ip in selectel_baremetal_server.test_server_2[0].ip_addresses : ip.address if ip.type == "public"]) > 0 ? 
                     "ssh root@${[for ip in selectel_baremetal_server.test_server_2[0].ip_addresses : ip.address if ip.type == "public"][0]}" : 
                     "No public IP available"
    } : null
  }
  
  sensitive = false
} 
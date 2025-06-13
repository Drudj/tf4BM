package selectel

// flattenServerCPU преобразует ServerCPU в формат для Terraform
func flattenServerCPU(cpu *ServerCPU) []interface{} {
	if cpu == nil {
		return []interface{}{}
	}

	cpuMap := map[string]interface{}{
		"model":     cpu.Model,
		"cores":     cpu.Cores,
		"threads":   cpu.Threads,
		"frequency": cpu.Frequency,
		"cache":     cpu.Cache,
	}

	return []interface{}{cpuMap}
}

// flattenServerRAM преобразует ServerRAM в формат для Terraform
func flattenServerRAM(ram *ServerRAM) []interface{} {
	if ram == nil {
		return []interface{}{}
	}

	ramMap := map[string]interface{}{
		"size": ram.Size,
		"type": ram.Type,
		"ecc":  ram.ECC,
	}

	return []interface{}{ramMap}
}

// flattenServerStorage преобразует массив ServerStorage в формат для Terraform
func flattenServerStorage(storage []*ServerStorage) []interface{} {
	if storage == nil {
		return []interface{}{}
	}

	storageList := make([]interface{}, len(storage))
	for i, s := range storage {
		storageMap := map[string]interface{}{
			"type":  s.Type,
			"size":  s.Size,
			"count": s.Count,
			"raid":  s.RAID,
		}
		storageList[i] = storageMap
	}

	return storageList
}

// flattenServerNetwork преобразует ServerNetwork в формат для Terraform
func flattenServerNetwork(network *ServerNetwork) []interface{} {
	if network == nil {
		return []interface{}{}
	}

	networkMap := map[string]interface{}{
		"primary_ip":     network.PrimaryIP,
		"gateway":        network.Gateway,
		"netmask":        network.Netmask,
		"additional_ips": network.AdditionalIPs,
		"bandwidth":      network.Bandwidth,
	}

	return []interface{}{networkMap}
}

// flattenServerLocation преобразует ServerLocation в формат для Terraform
func flattenServerLocation(location *ServerLocation) []interface{} {
	if location == nil {
		return []interface{}{}
	}

	locationMap := map[string]interface{}{
		"id":         location.ID,
		"name":       location.Name,
		"code":       location.Code,
		"country":    location.Country,
		"city":       location.City,
		"datacenter": location.Datacenter,
	}

	return []interface{}{locationMap}
}

// flattenServerOS преобразует ServerOS в формат для Terraform
func flattenServerOS(os *ServerOS) []interface{} {
	if os == nil {
		return []interface{}{}
	}

	osMap := map[string]interface{}{
		"id":           os.ID,
		"name":         os.Name,
		"version":      os.Version,
		"architecture": os.Architecture,
		"type":         os.Type,
		"distribution": os.Distribution,
	}

	return []interface{}{osMap}
}

// flattenServerIPMI преобразует ServerIPMI в формат для Terraform
func flattenServerIPMI(ipmi *ServerIPMI) []interface{} {
	if ipmi == nil {
		return []interface{}{}
	}

	ipmiMap := map[string]interface{}{
		"enabled": ipmi.Enabled,
		"ip":      ipmi.IP,
		"login":   ipmi.Login,
		// Пароль не возвращаем по соображениям безопасности
	}

	return []interface{}{ipmiMap}
}

// flattenServerBackup преобразует ServerBackup в формат для Terraform
func flattenServerBackup(backup *ServerBackup) []interface{} {
	if backup == nil {
		return []interface{}{}
	}

	backupMap := map[string]interface{}{
		"enabled":   backup.Enabled,
		"schedule":  backup.Schedule,
		"retention": backup.Retention,
	}

	return []interface{}{backupMap}
}

// flattenServerPrice преобразует ServerPrice в формат для Terraform
func flattenServerPrice(price *ServerPrice) []interface{} {
	if price == nil {
		return []interface{}{}
	}

	priceMap := map[string]interface{}{
		"amount":   price.Amount,
		"currency": price.Currency,
		"period":   price.Period,
	}

	return []interface{}{priceMap}
}

// flattenServerConfigurations преобразует массив ServerConfiguration в формат для Terraform
func flattenServerConfigurations(configs []*ServerConfiguration) []interface{} {
	if configs == nil {
		return []interface{}{}
	}

	configList := make([]interface{}, len(configs))
	for i, config := range configs {
		configMap := map[string]interface{}{
			"id":           config.ID,
			"name":         config.Name,
			"cpu":          flattenServerCPU(config.CPU),
			"ram":          flattenServerRAM(config.RAM),
			"storage":      flattenServerStorage(config.Storage),
			"price":        flattenServerPrice(config.Price),
			"location_ids": config.LocationIDs,
			"available":    config.Available,
		}
		configList[i] = configMap
	}

	return configList
}

// flattenServerLocations преобразует массив ServerLocation в формат для Terraform
func flattenServerLocations(locations []*ServerLocation) []interface{} {
	if locations == nil {
		return []interface{}{}
	}

	locationList := make([]interface{}, len(locations))
	for i, location := range locations {
		locationMap := map[string]interface{}{
			"id":         location.ID,
			"uuid":       location.UUID,
			"name":       location.Name,
			"code":       location.Code,
			"country":    location.Country,
			"city":       location.City,
			"datacenter": location.Datacenter,
		}
		locationList[i] = locationMap
	}

	return locationList
}

// flattenServerOSList преобразует массив ServerOS в формат для Terraform
func flattenServerOSList(osList []*ServerOS) []interface{} {
	if osList == nil {
		return []interface{}{}
	}

	osListFlattened := make([]interface{}, len(osList))
	for i, os := range osList {
		osMap := map[string]interface{}{
			"id":           os.ID,
			"name":         os.Name,
			"version":      os.Version,
			"architecture": os.Architecture,
			"type":         os.Type,
			"distribution": os.Distribution,
		}
		osListFlattened[i] = osMap
	}

	return osListFlattened
}

// flattenDedicatedServers преобразует массив DedicatedServer в формат для Terraform
func flattenDedicatedServers(servers []*DedicatedServer) []interface{} {
	if servers == nil {
		return []interface{}{}
	}

	serverList := make([]interface{}, len(servers))
	for i, server := range servers {
		serverMap := map[string]interface{}{
			"id":        server.ID,
			"name":      server.Name,
			"status":    server.Status,
			"status_hd": server.StatusHD,
			"cpu":       flattenServerCPU(server.CPU),
			"ram":       flattenServerRAM(server.RAM),
			"storage":   flattenServerStorage(server.Storage),
			"network":   flattenServerNetwork(server.Network),
			"location":  flattenServerLocation(server.Location),
			"os":        flattenServerOS(server.OS),
			"ipmi":      flattenServerIPMI(server.IPMI),
			"backup":    flattenServerBackup(server.Backup),
			"price":     flattenServerPrice(server.Price),
			"comment":   server.Comment,
			"tags":      server.Tags,
		}

		if server.CreatedAt != nil {
			serverMap["created_at"] = server.CreatedAt.Format("2006-01-02T15:04:05Z")
		}

		if server.UpdatedAt != nil {
			serverMap["updated_at"] = server.UpdatedAt.Format("2006-01-02T15:04:05Z")
		}

		serverList[i] = serverMap
	}

	return serverList
}

// flattenServerServicesList преобразует массив ServerService в формат для Terraform
func flattenServerServicesList(services []*ServerService) []interface{} {
	if services == nil {
		return []interface{}{}
	}

	serviceList := make([]interface{}, len(services))
	for i, service := range services {
		serviceMap := map[string]interface{}{
			"uuid":        service.UUID,
			"name":        service.Name,
			"description": service.Description,
			"type":        service.Type,
			"state":       service.State,
		}
		serviceList[i] = serviceMap
	}

	return serviceList
}

package models

// ServiceType представляет тип услуги
type ServiceType string

const (
	ServiceTypeServer     ServiceType = "server"
	ServiceTypeColocation ServiceType = "colocation"
	ServiceTypeNetwork    ServiceType = "network"
	ServiceTypeBackup     ServiceType = "backup"
)

// Service представляет услугу в каталоге
type Service struct {
	UUID         string      `json:"uuid"`
	Name         string      `json:"name"`
	Type         ServiceType `json:"type"`
	Category     string      `json:"category"`
	Description  string      `json:"description,omitempty"`
	Available    bool        `json:"available"`
	LocationUUID string      `json:"location_uuid,omitempty"`

	// Характеристики сервера (если тип = server)
	CPU     *CPUSpec      `json:"cpu,omitempty"`
	Memory  *MemorySpec   `json:"memory,omitempty"`
	Storage []StorageSpec `json:"storage,omitempty"`
	Network *NetworkSpec  `json:"network,omitempty"`

	// Цены
	PricePlans []PricePlan `json:"price_plans,omitempty"`
}

// CPUSpec представляет характеристики процессора
type CPUSpec struct {
	Model        string `json:"model"`
	Cores        int    `json:"cores"`
	Threads      int    `json:"threads"`
	Frequency    string `json:"frequency"`
	Cache        string `json:"cache,omitempty"`
	Architecture string `json:"architecture"`
}

// MemorySpec представляет характеристики памяти
type MemorySpec struct {
	Size     int    `json:"size"`     // в GB
	Type     string `json:"type"`     // DDR4, DDR5, etc.
	Speed    string `json:"speed"`    // частота
	Channels int    `json:"channels"` // количество каналов
}

// StorageType представляет тип накопителя
type StorageType string

const (
	StorageTypeHDD  StorageType = "hdd"
	StorageTypeSSD  StorageType = "ssd"
	StorageTypeNVMe StorageType = "nvme"
)

// StorageSpec представляет характеристики накопителя
type StorageSpec struct {
	Type      StorageType `json:"type"`
	Size      int         `json:"size"`      // в GB
	Count     int         `json:"count"`     // количество дисков
	Interface string      `json:"interface"` // SATA, SAS, NVMe
	Model     string      `json:"model,omitempty"`
}

// NetworkSpec представляет сетевые характеристики
type NetworkSpec struct {
	Bandwidth int      `json:"bandwidth"` // в Mbps
	Ports     []string `json:"ports"`     // типы портов
	IPv4      bool     `json:"ipv4"`
	IPv6      bool     `json:"ipv6"`
}

// PricePlan представляет тарифный план
type PricePlan struct {
	UUID        string  `json:"uuid"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`      // monthly, hourly, etc.
	Price       float64 `json:"price"`     // цена
	Currency    string  `json:"currency"`  // валюта
	Period      string  `json:"period"`    // период
	SetupFee    float64 `json:"setup_fee"` // разовая плата
	Available   bool    `json:"available"`
	Description string  `json:"description,omitempty"`
}

// ServicesResponse представляет ответ со списком услуг
type ServicesResponse struct {
	Services []Service `json:"services"`
}

// PricePlansResponse представляет ответ со списком тарифных планов
type PricePlansResponse struct {
	PricePlans []PricePlan `json:"price_plans"`
}

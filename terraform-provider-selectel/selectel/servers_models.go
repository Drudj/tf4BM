package selectel

import (
	"time"
)

// DedicatedServer представляет выделенный сервер Selectel
type DedicatedServer struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	StatusHD string `json:"status_hd"`

	// Конфигурация
	CPU      *ServerCPU       `json:"cpu,omitempty"`
	RAM      *ServerRAM       `json:"ram,omitempty"`
	Storage  []*ServerStorage `json:"storage,omitempty"`
	Network  *ServerNetwork   `json:"network,omitempty"`
	Location *ServerLocation  `json:"location,omitempty"`

	// Операционная система
	OS *ServerOS `json:"os,omitempty"`

	// Дополнительные услуги
	IPMI   *ServerIPMI   `json:"ipmi,omitempty"`
	Backup *ServerBackup `json:"backup,omitempty"`

	// Временные метки
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// Биллинг
	Price *ServerPrice `json:"price,omitempty"`

	// Дополнительная информация
	Comment string   `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// ServerCPU представляет информацию о процессоре
type ServerCPU struct {
	Model     string `json:"model"`
	Cores     int    `json:"cores"`
	Threads   int    `json:"threads"`
	Frequency string `json:"frequency"`
	Cache     string `json:"cache,omitempty"`
}

// ServerRAM представляет информацию о оперативной памяти
type ServerRAM struct {
	Size string `json:"size"`
	Type string `json:"type"`
	ECC  bool   `json:"ecc,omitempty"`
}

// ServerStorage представляет информацию о хранилище
type ServerStorage struct {
	Type  string `json:"type"`           // "HDD", "SSD", "NVMe"
	Size  string `json:"size"`           // "1TB", "500GB"
	Count int    `json:"count"`          // количество дисков
	RAID  string `json:"raid,omitempty"` // уровень RAID
}

// ServerNetwork представляет сетевую конфигурацию
type ServerNetwork struct {
	// Основная сеть
	PrimaryIP string `json:"primary_ip,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Netmask   string `json:"netmask,omitempty"`

	// Дополнительные IP
	AdditionalIPs []string `json:"additional_ips,omitempty"`

	// Публичная сеть
	PublicNetwork *ServerNetworkConfig `json:"public_network,omitempty"`

	// Приватная сеть
	PrivateNetwork *ServerNetworkConfig `json:"private_network,omitempty"`

	// Пропускная способность
	Bandwidth string `json:"bandwidth,omitempty"`
}

// ServerNetworkConfig представляет конфигурацию сети
type ServerNetworkConfig struct {
	IP      string `json:"ip"`
	Gateway string `json:"gateway"`
	Netmask string `json:"netmask"`
	VLAN    int    `json:"vlan,omitempty"`
}

// ServerLocation представляет местоположение сервера
type ServerLocation struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	LocationID  int    `json:"location_id"`
	Description string `json:"description"`
	DCCount     *int   `json:"dc_count"`
	Visibility  string `json:"visibility"`

	// Для совместимости с UI добавляем маппинг на старые поля
	ID         int    `json:"-"`
	Code       string `json:"-"`
	Country    string `json:"-"`
	City       string `json:"-"`
	Datacenter string `json:"-"`
}

// ServerOS представляет операционную систему
type ServerOS struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`           // "x86_64"
	Type         string `json:"type"`                   // "linux", "windows"
	Distribution string `json:"distribution,omitempty"` // "Ubuntu", "CentOS"
}

// ServerIPMI представляет настройки IPMI
type ServerIPMI struct {
	Enabled  bool   `json:"enabled"`
	IP       string `json:"ip,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

// ServerBackup представляет настройки резервного копирования
type ServerBackup struct {
	Enabled   bool   `json:"enabled"`
	Schedule  string `json:"schedule,omitempty"`  // "daily", "weekly"
	Retention int    `json:"retention,omitempty"` // дни хранения
}

// ServerPrice представляет информацию о стоимости
type ServerPrice struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Period   string  `json:"period"` // "monthly", "hourly"
}

// DedicatedServerCreate содержит данные для создания нового сервера
type DedicatedServerCreate struct {
	Name       string `json:"name"`
	LocationID int    `json:"location_id"`
	ConfigID   int    `json:"config_id,omitempty"`
	OSID       int    `json:"os_id,omitempty"`

	// Дополнительные опции
	Comment string   `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`

	// Сетевые настройки
	NetworkConfig *DedicatedServerNetworkCreate `json:"network_config,omitempty"`

	// Дополнительные услуги
	EnableIPMI   bool `json:"enable_ipmi,omitempty"`
	EnableBackup bool `json:"enable_backup,omitempty"`

	// SSH ключи
	SSHKeys []string `json:"ssh_keys,omitempty"`
}

// DedicatedServerNetworkCreate содержит сетевые настройки для создания сервера
type DedicatedServerNetworkCreate struct {
	AdditionalIPs  int  `json:"additional_ips,omitempty"`
	PrivateNetwork bool `json:"private_network,omitempty"`
}

// DedicatedServerUpdate содержит данные для обновления сервера
type DedicatedServerUpdate struct {
	Name    *string  `json:"name,omitempty"`
	Comment *string  `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// DedicatedServerAction представляет действие над сервером
type DedicatedServerAction struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// ServerConfiguration представляет доступную конфигурацию сервера
type ServerConfiguration struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	CPU         *ServerCPU       `json:"cpu"`
	RAM         *ServerRAM       `json:"ram"`
	Storage     []*ServerStorage `json:"storage"`
	Price       *ServerPrice     `json:"price"`
	LocationIDs []int            `json:"location_ids"`
	Available   bool             `json:"available"`
}

// ServerTaskStatus представляет статус выполнения задачи
type ServerTaskStatus struct {
	ID          int        `json:"id"`
	Status      string     `json:"status"`
	Progress    int        `json:"progress"`
	Message     string     `json:"message,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Error       string     `json:"error,omitempty"`
}

// Константы статусов серверов
var (
	ServerStatusActive      = "active"
	ServerStatusInstalling  = "installing"
	ServerStatusRebooting   = "rebooting"
	ServerStatusMaintenance = "maintenance"
	ServerStatusError       = "error"
	ServerStatusStopped     = "stopped"
)

// Константы действий над серверами
var (
	ServerActionStart      = "start"
	ServerActionStop       = "stop"
	ServerActionRestart    = "restart"
	ServerActionReinstall  = "reinstall"
	ServerActionRescue     = "rescue"
	ServerActionPowerCycle = "power_cycle"
)

// Константы статусов задач
var (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
	TaskStatusCancelled = "cancelled"
)

// ServerService представляет сервис для получения service_uuid
type ServerService struct {
	ID          string `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	State       string `json:"state,omitempty"`
	Description string `json:"description,omitempty"`
	Region      string `json:"region,omitempty"`
}

// DedicatedServerCreateBilling содержит данные для создания сервера через биллинг API
type DedicatedServerCreateBilling struct {
	Name         string `json:"name"`
	LocationUUID string `json:"location_uuid"`
	ServiceUUID  string `json:"service_uuid"`
	ConfigID     int    `json:"config_id,omitempty"`
	OSID         int    `json:"os_id,omitempty"`

	// Дополнительные опции
	Comment string   `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`

	// SSH ключи
	SSHKeys []string `json:"ssh_keys,omitempty"`

	// Биллинговые опции
	Period      string `json:"period,omitempty"` // "monthly", "hourly"
	AutoRenewal bool   `json:"auto_renewal,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`

	// Расширенные поля для биллинг API
	PricePlanUUID    string      `json:"price_plan_uuid,omitempty"`
	OSTemplate       string      `json:"os_template,omitempty"`
	Arch             string      `json:"arch,omitempty"`
	Version          string      `json:"version,omitempty"`
	UserHostname     string      `json:"userhostname,omitempty"`
	PayCurrency      string      `json:"pay_currency,omitempty"`
	UserDesc         string      `json:"user_desc,omitempty"`
	PartitionsConfig interface{} `json:"partitions_config,omitempty"`
}

// DedicatedServerCreateResponse представляет ответ на создание сервера
type DedicatedServerCreateResponse struct {
	ID        string `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	OrderID   string `json:"order_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	TaskID    string `json:"task_id,omitempty"`

	// Дополнительная информация
	Message   string     `json:"message,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// Результат создания (массив серверов)
	Result []DedicatedServerCreateResult `json:"result,omitempty"`
}

// DedicatedServerCreateResult представляет элемент результата создания сервера
type DedicatedServerCreateResult struct {
	UUID    string `json:"uuid"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	OrderID string `json:"order_id,omitempty"`
}

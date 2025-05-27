package models

import (
	"time"
)

// ServerStatus представляет статус сервера
type ServerStatus string

const (
	ServerStatusOrdering     ServerStatus = "ordering"
	ServerStatusProvisioning ServerStatus = "provisioning"
	ServerStatusActive       ServerStatus = "active"
	ServerStatusStopped      ServerStatus = "stopped"
	ServerStatusMaintenance  ServerStatus = "maintenance"
	ServerStatusError        ServerStatus = "error"
	ServerStatusCancelled    ServerStatus = "cancelled"
)

// PowerStatus представляет статус питания
type PowerStatus string

const (
	PowerStatusOn  PowerStatus = "on"
	PowerStatusOff PowerStatus = "off"
)

// Server представляет выделенный сервер
type Server struct {
	UUID          string       `json:"uuid"`
	Name          string       `json:"name"`
	Status        ServerStatus `json:"status"`
	PowerStatus   PowerStatus  `json:"power_status"`
	ServiceUUID   string       `json:"service_uuid"`
	LocationUUID  string       `json:"location_uuid"`
	ProjectUUID   string       `json:"project_uuid"`
	PricePlanUUID string       `json:"price_plan_uuid"`

	// Сетевые настройки
	Networks []ServerNetwork `json:"networks,omitempty"`

	// Операционная система
	OS *ServerOS `json:"os,omitempty"`

	// Дополнительные услуги
	AdditionalServices []string `json:"additional_services,omitempty"`

	// Теги
	Tags map[string]string `json:"tags,omitempty"`

	// Временные метки
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`

	// Информация о конфигурации (только для чтения)
	Service  *Service  `json:"service,omitempty"`
	Location *Location `json:"location,omitempty"`
}

// NetworkType представляет тип сети
type NetworkType string

const (
	NetworkTypePublic  NetworkType = "public"
	NetworkTypePrivate NetworkType = "private"
)

// ServerNetwork представляет сетевую конфигурацию сервера
type ServerNetwork struct {
	Type      NetworkType `json:"type"`
	Bandwidth int         `json:"bandwidth"` // в Mbps
	IPv4      []string    `json:"ipv4,omitempty"`
	IPv6      []string    `json:"ipv6,omitempty"`
	Subnet    string      `json:"subnet,omitempty"`
	VLANID    int         `json:"vlan_id,omitempty"`
	Gateway   string      `json:"gateway,omitempty"`
}

// ServerOS представляет операционную систему сервера
type ServerOS struct {
	TemplateUUID string   `json:"template_uuid"`
	Name         string   `json:"name,omitempty"`
	Version      string   `json:"version,omitempty"`
	Architecture string   `json:"architecture,omitempty"`
	SSHKeys      []string `json:"ssh_keys,omitempty"`
	Password     string   `json:"password,omitempty"`
	UserData     string   `json:"user_data,omitempty"`
}

// OSTemplate представляет шаблон операционной системы
type OSTemplate struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Type         string `json:"type"`   // linux, windows, etc.
	Family       string `json:"family"` // ubuntu, centos, windows, etc.
	Available    bool   `json:"available"`
	Description  string `json:"description,omitempty"`
	MinDiskSize  int    `json:"min_disk_size,omitempty"` // в GB
}

// CreateServerRequest представляет запрос на создание сервера
type CreateServerRequest struct {
	Name               string            `json:"name"`
	ServiceUUID        string            `json:"service_uuid"`
	LocationUUID       string            `json:"location_uuid"`
	ProjectUUID        string            `json:"project_uuid"`
	PricePlanUUID      string            `json:"price_plan_uuid"`
	Networks           []ServerNetwork   `json:"networks,omitempty"`
	OS                 *ServerOS         `json:"os,omitempty"`
	AdditionalServices []string          `json:"additional_services,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
}

// UpdateServerRequest представляет запрос на обновление сервера
type UpdateServerRequest struct {
	Name               string            `json:"name,omitempty"`
	AdditionalServices []string          `json:"additional_services,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
}

// ServersResponse представляет ответ со списком серверов
type ServersResponse struct {
	Servers []Server `json:"servers"`
}

// OSTemplatesResponse представляет ответ со списком шаблонов ОС
type OSTemplatesResponse struct {
	Templates []OSTemplate `json:"templates"`
}

// IsActive проверяет, активен ли сервер
func (s *Server) IsActive() bool {
	return s.Status == ServerStatusActive
}

// IsProvisioning проверяет, находится ли сервер в процессе создания
func (s *Server) IsProvisioning() bool {
	return s.Status == ServerStatusOrdering || s.Status == ServerStatusProvisioning
}

// HasError проверяет, есть ли ошибка у сервера
func (s *Server) HasError() bool {
	return s.Status == ServerStatusError
}

// IsPoweredOn проверяет, включен ли сервер
func (s *Server) IsPoweredOn() bool {
	return s.PowerStatus == PowerStatusOn
}

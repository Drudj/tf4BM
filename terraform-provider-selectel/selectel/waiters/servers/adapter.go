package servers

import (
	"context"
	"time"
)

// ServiceAdapter адаптирует ServersService для работы с waiters
type ServiceAdapter struct {
	service ServersServiceInterface
}

// NewServiceAdapter создает новый adapter
func NewServiceAdapter(service ServersServiceInterface) *ServiceAdapter {
	return &ServiceAdapter{service: service}
}

// GetServer получает сервер через adapter
func (a *ServiceAdapter) GetServer(ctx context.Context, serverID int) (*DedicatedServer, error) {
	server, err := a.service.GetServer(ctx, serverID)
	if err != nil {
		return nil, err
	}

	// Конвертируем основной тип в waiter тип
	return &DedicatedServer{
		ID:     server.ID,
		Status: server.Status,
	}, nil
}

// GetTask получает задачу через adapter
func (a *ServiceAdapter) GetTask(ctx context.Context, taskID int) (*ServerTaskStatus, error) {
	task, err := a.service.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Конвертируем основной тип в waiter тип
	waiterTask := &ServerTaskStatus{
		ID:       task.ID,
		Status:   task.Status,
		Progress: task.Progress,
		Error:    task.Error,
		Message:  task.Message,
	}

	return waiterTask, nil
}

// ServersServiceInterface интерфейс для основного сервиса серверов
type ServersServiceInterface interface {
	GetServer(ctx context.Context, serverID int) (*MainDedicatedServer, error)
	GetTask(ctx context.Context, taskID int) (*MainServerTaskStatus, error)
}

// MainDedicatedServer основной тип сервера
type MainDedicatedServer struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Status    string          `json:"status"`
	StatusHD  string          `json:"status_hd"`
	Comment   string          `json:"comment"`
	Tags      []string        `json:"tags"`
	CPU       *ServerCPU      `json:"cpu"`
	RAM       *ServerRAM      `json:"ram"`
	Storage   []ServerStorage `json:"storage"`
	Network   *ServerNetwork  `json:"network"`
	Location  *ServerLocation `json:"location"`
	OS        *ServerOS       `json:"os"`
	IPMI      *ServerIPMI     `json:"ipmi"`
	Backup    *ServerBackup   `json:"backup"`
	Price     *ServerPrice    `json:"price"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
}

// MainServerTaskStatus основной тип задачи
type MainServerTaskStatus struct {
	ID          int        `json:"id"`
	Status      string     `json:"status"`
	Progress    int        `json:"progress"`
	Message     string     `json:"message,omitempty"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Вспомогательные типы для совместимости
type ServerCPU struct {
	Model     string `json:"model"`
	Cores     int    `json:"cores"`
	Threads   int    `json:"threads"`
	Frequency string `json:"frequency"`
	Cache     string `json:"cache"`
}

type ServerRAM struct {
	Size string `json:"size"`
	Type string `json:"type"`
	ECC  bool   `json:"ecc"`
}

type ServerStorage struct {
	Type  string `json:"type"`
	Size  string `json:"size"`
	Count int    `json:"count"`
	RAID  string `json:"raid"`
}

type ServerNetwork struct {
	PrimaryIP     string   `json:"primary_ip"`
	Gateway       string   `json:"gateway"`
	Netmask       string   `json:"netmask"`
	AdditionalIPs []string `json:"additional_ips"`
	Bandwidth     string   `json:"bandwidth"`
}

type ServerLocation struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Country    string `json:"country"`
	City       string `json:"city"`
	Datacenter string `json:"datacenter"`
}

type ServerOS struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Type         string `json:"type"`
	Distribution string `json:"distribution"`
}

type ServerIPMI struct {
	Enabled bool   `json:"enabled"`
	IP      string `json:"ip"`
	Login   string `json:"login"`
}

type ServerBackup struct {
	Enabled   bool   `json:"enabled"`
	Schedule  string `json:"schedule"`
	Retention int    `json:"retention"`
}

type ServerPrice struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Period   string  `json:"period"`
}

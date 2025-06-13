package selectel

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ServersService предоставляет методы для работы с выделенными серверами
type ServersService struct {
	client *ServersClient
}

// NewServersService создает новый сервис для работы с серверами
func NewServersService(client *ServersClient) *ServersService {
	return &ServersService{
		client: client,
	}
}

// ListServers возвращает список серверов
func (s *ServersService) ListServers(ctx context.Context, opts *ServersListOptions) ([]*DedicatedServer, error) {
	path := "server"

	if opts != nil {
		path += opts.BuildQueryString()
	}

	resp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []*DedicatedServer `json:"data"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetServer возвращает информацию о конкретном сервере
func (s *ServersService) GetServer(ctx context.Context, serverID int) (*DedicatedServer, error) {
	path := fmt.Sprintf("server/%d", serverID)

	resp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data *DedicatedServer `json:"data"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// CreateServer создает новый сервер
func (s *ServersService) CreateServer(ctx context.Context, createOpts *DedicatedServerCreate) (*DedicatedServer, error) {
	path := "server"

	resp, err := s.client.DoRequest(ctx, http.MethodPost, path, createOpts)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data *DedicatedServer `json:"data"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// UpdateServer обновляет сервер
func (s *ServersService) UpdateServer(ctx context.Context, serverID int, updateOpts *DedicatedServerUpdate) (*DedicatedServer, error) {
	path := fmt.Sprintf("server/%d", serverID)

	resp, err := s.client.DoRequest(ctx, http.MethodPatch, path, updateOpts)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data *DedicatedServer `json:"data"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// DeleteServer удаляет сервер
func (s *ServersService) DeleteServer(ctx context.Context, serverID int) error {
	path := fmt.Sprintf("server/%d", serverID)

	resp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return s.client.ParseResponse(resp, nil)
}

// ServerAction выполняет действие над сервером (start, stop, restart и т.д.)
func (s *ServersService) ServerAction(ctx context.Context, serverID int, action *DedicatedServerAction) (*ServerTaskStatus, error) {
	path := fmt.Sprintf("server/%d/action", serverID)

	resp, err := s.client.DoRequest(ctx, http.MethodPost, path, action)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data *ServerTaskStatus `json:"data"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// StartServer запускает сервер
func (s *ServersService) StartServer(ctx context.Context, serverID int) (*ServerTaskStatus, error) {
	action := &DedicatedServerAction{
		Action: ServerActionStart,
	}
	return s.ServerAction(ctx, serverID, action)
}

// StopServer останавливает сервер
func (s *ServersService) StopServer(ctx context.Context, serverID int) (*ServerTaskStatus, error) {
	action := &DedicatedServerAction{
		Action: ServerActionStop,
	}
	return s.ServerAction(ctx, serverID, action)
}

// RestartServer перезапускает сервер
func (s *ServersService) RestartServer(ctx context.Context, serverID int) (*ServerTaskStatus, error) {
	action := &DedicatedServerAction{
		Action: ServerActionRestart,
	}
	return s.ServerAction(ctx, serverID, action)
}

// ReinstallServer переустанавливает ОС на сервере
func (s *ServersService) ReinstallServer(ctx context.Context, serverID int, osID int, sshKeys []string) (*ServerTaskStatus, error) {
	params := map[string]interface{}{
		"os_id": osID,
	}

	if len(sshKeys) > 0 {
		params["ssh_keys"] = sshKeys
	}

	action := &DedicatedServerAction{
		Action: ServerActionReinstall,
		Params: params,
	}

	return s.ServerAction(ctx, serverID, action)
}

// PowerCycleServer выполняет жесткий перезапуск сервера
func (s *ServersService) PowerCycleServer(ctx context.Context, serverID int) (*ServerTaskStatus, error) {
	action := &DedicatedServerAction{
		Action: ServerActionPowerCycle,
	}
	return s.ServerAction(ctx, serverID, action)
}

// GetTask возвращает статус выполнения задачи
func (s *ServersService) GetTask(ctx context.Context, taskID int) (*ServerTaskStatus, error) {
	path := fmt.Sprintf("task/%d", taskID)

	resp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data *ServerTaskStatus `json:"data"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// WaitForTask ожидает завершения задачи с таймаутом
func (s *ServersService) WaitForTask(ctx context.Context, taskID int) (*ServerTaskStatus, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			task, err := s.GetTask(ctx, taskID)
			if err != nil {
				return nil, err
			}

			switch task.Status {
			case TaskStatusCompleted:
				return task, nil
			case TaskStatusFailed:
				return task, fmt.Errorf("task failed: %s", task.Error)
			case TaskStatusCancelled:
				return task, fmt.Errorf("task was cancelled")
			default:
				// Задача еще выполняется, ждем
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(5 * time.Second):
					continue
				}
			}
		}
	}
}

// ListConfigurations возвращает список доступных конфигураций серверов
func (s *ServersService) ListConfigurations(ctx context.Context) ([]*ServerConfiguration, error) {
	path := "configuration"

	resp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []*ServerConfiguration `json:"data"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetConfiguration возвращает информацию о конкретной конфигурации
func (s *ServersService) GetConfiguration(ctx context.Context, configID int) (*ServerConfiguration, error) {
	path := fmt.Sprintf("configuration/%d", configID)

	resp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data *ServerConfiguration `json:"data"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// ListLocations возвращает список доступных локаций
func (s *ServersService) ListLocations(ctx context.Context) ([]*ServerLocation, error) {
	path := "location"

	resp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		TaskID    string            `json:"task_id"`
		Status    string            `json:"status"`
		Progress  int               `json:"progress"`
		Page      int               `json:"page"`
		Limit     int               `json:"limit"`
		ItemCount int               `json:"item_count"`
		Result    []*ServerLocation `json:"result"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	// Заполняем совместимые поля для UI
	for _, location := range result.Result {
		location.ID = location.LocationID
		// Пытаемся извлечь код локации из имени (SPB-4 -> SPB)
		if len(location.Name) >= 3 {
			location.Code = location.Name[:3]
		}
		location.Datacenter = location.Description
	}

	return result.Result, nil
}

// ListOperatingSystems возвращает список доступных операционных систем
func (s *ServersService) ListOperatingSystems(ctx context.Context) ([]*ServerOS, error) {
	path := "os"

	resp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Result []*ServerOS `json:"result"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

// GetServices возвращает список доступных сервисов для получения service_uuid
func (s *ServersService) GetServices(ctx context.Context) ([]*ServerService, error) {
	path := "service"

	resp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Result []*ServerService `json:"result"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Result, nil
}

// ListOperatingSystemsNew возвращает список доступных операционных систем через новый эндпоинт
func (s *ServersService) ListOperatingSystemsNew(ctx context.Context, locationUUID, serviceUUID string) ([]*ServerOS, error) {
	path := fmt.Sprintf("boot/template/os/new?location_uuid=%s&service_uuid=%s", locationUUID, serviceUUID)

	resp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []*ServerOS `json:"data"`
	}

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// CreateServerResource создает новый выделенный сервер через эндпоинт биллинга
func (s *ServersService) CreateServerResource(ctx context.Context, createOpts *DedicatedServerCreateBilling) (*DedicatedServerCreateResponse, error) {
	path := "resource/serverchip/billing"

	// Логируем структуру перед отправкой
	log.Printf("[DEBUG] CreateServerResource: UserHostname='%s', UserDesc='%s'", createOpts.UserHostname, createOpts.UserDesc)
	log.Printf("[DEBUG] CreateServerResource: Full struct: %+v", createOpts)

	resp, err := s.client.DoRequest(ctx, http.MethodPost, path, createOpts)
	if err != nil {
		return nil, err
	}

	var result DedicatedServerCreateResponse

	if err := s.client.ParseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

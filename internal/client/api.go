package client

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Drudj/tf_for_BareMetal/internal/models"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// APIClient представляет клиент для работы с API Selectel
type APIClient struct {
	*Client
}

// NewAPIClient создает новый API клиент
func NewAPIClient(config Config) *APIClient {
	return &APIClient{
		Client: NewClient(config),
	}
}

// Locations API

// GetLocations получает список всех локаций
func (c *APIClient) GetLocations(ctx context.Context) ([]models.Location, error) {
	var response models.LocationsResponse
	err := c.Get(ctx, "location", &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}
	return response.Locations, nil
}

// GetLocation получает локацию по UUID
func (c *APIClient) GetLocation(ctx context.Context, uuid string) (*models.Location, error) {
	var location models.Location
	err := c.Get(ctx, fmt.Sprintf("location/%s", uuid), &location)
	if err != nil {
		return nil, fmt.Errorf("failed to get location %s: %w", uuid, err)
	}
	return &location, nil
}

// Services API

// GetServices получает список всех услуг
func (c *APIClient) GetServices(ctx context.Context) ([]models.Service, error) {
	var response models.ServicesResponse
	err := c.Get(ctx, "service", &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}
	return response.Services, nil
}

// GetServicesWithFilters получает список услуг с фильтрами
func (c *APIClient) GetServicesWithFilters(ctx context.Context, filters map[string]string) ([]models.Service, error) {
	params := url.Values{}
	for key, value := range filters {
		params.Add(key, value)
	}

	var response models.ServicesResponse
	err := c.GetWithQuery(ctx, "service", params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get services with filters: %w", err)
	}
	return response.Services, nil
}

// GetService получает услугу по UUID
func (c *APIClient) GetService(ctx context.Context, uuid string) (*models.Service, error) {
	var service models.Service
	err := c.Get(ctx, fmt.Sprintf("service/%s", uuid), &service)
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %w", uuid, err)
	}
	return &service, nil
}

// Price Plans API

// GetPricePlans получает список тарифных планов
func (c *APIClient) GetPricePlans(ctx context.Context) ([]models.PricePlan, error) {
	var response models.PricePlansResponse
	err := c.Get(ctx, "plan", &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get price plans: %w", err)
	}
	return response.PricePlans, nil
}

// GetPricePlansForService получает тарифные планы для конкретной услуги
func (c *APIClient) GetPricePlansForService(ctx context.Context, serviceUUID string) ([]models.PricePlan, error) {
	params := url.Values{}
	params.Add("service_uuid", serviceUUID)

	var response models.PricePlansResponse
	err := c.GetWithQuery(ctx, "plan", params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get price plans for service %s: %w", serviceUUID, err)
	}
	return response.PricePlans, nil
}

// OS Templates API

// GetOSTemplates получает список шаблонов операционных систем
func (c *APIClient) GetOSTemplates(ctx context.Context) ([]models.OSTemplate, error) {
	var response models.OSTemplatesResponse
	err := c.Get(ctx, "boot/template/os/new", &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get OS templates: %w", err)
	}
	return response.Templates, nil
}

// GetOSTemplate получает шаблон ОС по UUID
func (c *APIClient) GetOSTemplate(ctx context.Context, uuid string) (*models.OSTemplate, error) {
	var template models.OSTemplate
	err := c.Get(ctx, fmt.Sprintf("boot/template/os/new/%s", uuid), &template)
	if err != nil {
		return nil, fmt.Errorf("failed to get OS template %s: %w", uuid, err)
	}
	return &template, nil
}

// Servers API

// GetServers получает список серверов
func (c *APIClient) GetServers(ctx context.Context) ([]models.Server, error) {
	var response models.ServersResponse
	err := c.Get(ctx, "servers", &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get servers: %w", err)
	}
	return response.Servers, nil
}

// GetServer получает сервер по UUID
func (c *APIClient) GetServer(ctx context.Context, uuid string) (*models.Server, error) {
	var server models.Server
	err := c.Get(ctx, fmt.Sprintf("servers/%s", uuid), &server)
	if err != nil {
		return nil, fmt.Errorf("failed to get server %s: %w", uuid, err)
	}
	return &server, nil
}

// CreateServer создает новый сервер
func (c *APIClient) CreateServer(ctx context.Context, request models.CreateServerRequest) (*models.Server, error) {
	var server models.Server
	err := c.Post(ctx, "servers", request, &server)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}
	return &server, nil
}

// UpdateServer обновляет сервер
func (c *APIClient) UpdateServer(ctx context.Context, uuid string, request models.UpdateServerRequest) (*models.Server, error) {
	var server models.Server
	err := c.Patch(ctx, fmt.Sprintf("servers/%s", uuid), request, &server)
	if err != nil {
		return nil, fmt.Errorf("failed to update server %s: %w", uuid, err)
	}
	return &server, nil
}

// DeleteServer удаляет сервер
func (c *APIClient) DeleteServer(ctx context.Context, uuid string) error {
	err := c.Delete(ctx, fmt.Sprintf("servers/%s", uuid))
	if err != nil {
		return fmt.Errorf("failed to delete server %s: %w", uuid, err)
	}
	return nil
}

// Server Power Management

// PowerOnServer включает сервер
func (c *APIClient) PowerOnServer(ctx context.Context, uuid string) (*models.Task, error) {
	var task models.Task
	err := c.Post(ctx, fmt.Sprintf("servers/%s/power/on", uuid), nil, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to power on server %s: %w", uuid, err)
	}
	return &task, nil
}

// PowerOffServer выключает сервер
func (c *APIClient) PowerOffServer(ctx context.Context, uuid string) (*models.Task, error) {
	var task models.Task
	err := c.Post(ctx, fmt.Sprintf("servers/%s/power/off", uuid), nil, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to power off server %s: %w", uuid, err)
	}
	return &task, nil
}

// RebootServer перезагружает сервер
func (c *APIClient) RebootServer(ctx context.Context, uuid string) (*models.Task, error) {
	var task models.Task
	err := c.Post(ctx, fmt.Sprintf("servers/%s/power/reboot", uuid), nil, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to reboot server %s: %w", uuid, err)
	}
	return &task, nil
}

// Tasks API

// GetTask получает задачу по UUID
func (c *APIClient) GetTask(ctx context.Context, uuid string) (*models.Task, error) {
	var task models.Task
	err := c.Get(ctx, fmt.Sprintf("tasks/%s", uuid), &task)
	if err != nil {
		return nil, fmt.Errorf("failed to get task %s: %w", uuid, err)
	}
	return &task, nil
}

// WaitForTask ожидает завершения задачи
func (c *APIClient) WaitForTask(ctx context.Context, taskUUID string, timeout time.Duration) (*models.Task, error) {
	tflog.Debug(ctx, "Waiting for task completion", map[string]interface{}{
		"task_uuid": taskUUID,
		"timeout":   timeout,
	})

	deadline := time.Now().Add(timeout)
	checkInterval := 5 * time.Second

	for time.Now().Before(deadline) {
		task, err := c.GetTask(ctx, taskUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to check task status: %w", err)
		}

		tflog.Debug(ctx, "Task status check", map[string]interface{}{
			"task_uuid": taskUUID,
			"status":    task.Status,
			"progress":  task.Progress,
		})

		if task.IsCompleted() {
			tflog.Info(ctx, "Task completed successfully", map[string]interface{}{
				"task_uuid": taskUUID,
			})
			return task, nil
		}

		if task.IsFailed() {
			return task, fmt.Errorf("task failed: %s", task.Error)
		}

		select {
		case <-ctx.Done():
			return task, ctx.Err()
		case <-time.After(checkInterval):
			// Продолжаем ожидание
		}
	}

	// Получаем финальный статус задачи
	task, err := c.GetTask(ctx, taskUUID)
	if err != nil {
		return nil, fmt.Errorf("timeout waiting for task completion, failed to get final status: %w", err)
	}

	return task, fmt.Errorf("timeout waiting for task completion after %v", timeout)
}

// Helper methods

// FindLocationByName находит локацию по имени
func (c *APIClient) FindLocationByName(ctx context.Context, name string) (*models.Location, error) {
	locations, err := c.GetLocations(ctx)
	if err != nil {
		return nil, err
	}

	for _, location := range locations {
		if location.Name == name {
			return &location, nil
		}
	}

	return nil, fmt.Errorf("location with name '%s' not found", name)
}

// FindServiceByName находит услугу по имени
func (c *APIClient) FindServiceByName(ctx context.Context, name string) (*models.Service, error) {
	services, err := c.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		if service.Name == name {
			return &service, nil
		}
	}

	return nil, fmt.Errorf("service with name '%s' not found", name)
}

// FindOSTemplateByName находит шаблон ОС по имени
func (c *APIClient) FindOSTemplateByName(ctx context.Context, name string) (*models.OSTemplate, error) {
	templates, err := c.GetOSTemplates(ctx)
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		if template.Name == name {
			return &template, nil
		}
	}

	return nil, fmt.Errorf("OS template with name '%s' not found", name)
}

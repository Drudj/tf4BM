package servers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// WaiterConfig содержит настройки для ожидания состояний
type WaiterConfig struct {
	Pending      []string
	Target       []string
	Timeout      time.Duration
	Delay        time.Duration
	MinTimeout   time.Duration
	PollInterval time.Duration
}

// DefaultServerWaiterConfig возвращает стандартную конфигурацию ожидания для серверов
func DefaultServerWaiterConfig() *WaiterConfig {
	return &WaiterConfig{
		Timeout:      60 * time.Minute,
		Delay:        10 * time.Second,
		MinTimeout:   5 * time.Second,
		PollInterval: 15 * time.Second,
	}
}

// DefaultTaskWaiterConfig возвращает стандартную конфигурацию ожидания для задач
func DefaultTaskWaiterConfig() *WaiterConfig {
	return &WaiterConfig{
		Timeout:      30 * time.Minute,
		Delay:        5 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 10 * time.Second,
	}
}

// ServerStateRefreshFunc создает StateRefreshFunc для ожидания состояния сервера
func ServerStateRefreshFunc(ctx context.Context, serversService ServerService, serverID int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refreshing server %d state", serverID)

		server, err := serversService.GetServer(ctx, serverID)
		if err != nil {
			if isServerNotFoundError(err) {
				log.Printf("[DEBUG] Server %d not found, considering as deleted", serverID)
				return nil, "DELETED", nil
			}
			log.Printf("[ERROR] Error getting server %d: %s", serverID, err)
			return nil, "", err
		}

		log.Printf("[DEBUG] Server %d current status: %s", serverID, server.Status)
		return server, server.Status, nil
	}
}

// TaskStateRefreshFunc создает StateRefreshFunc для ожидания завершения задачи
func TaskStateRefreshFunc(ctx context.Context, serversService ServerService, taskID int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refreshing task %d state", taskID)

		task, err := serversService.GetTask(ctx, taskID)
		if err != nil {
			if isTaskNotFoundError(err) {
				log.Printf("[DEBUG] Task %d not found, considering as completed", taskID)
				return nil, "completed", nil
			}
			log.Printf("[ERROR] Error getting task %d: %s", taskID, err)
			return nil, "", err
		}

		log.Printf("[DEBUG] Task %d current status: %s, progress: %d%%", taskID, task.Status, task.Progress)

		// Если задача завершилась с ошибкой, возвращаем ошибку
		if task.Status == "failed" {
			return task, task.Status, fmt.Errorf("task %d failed: %s", taskID, task.Error)
		}

		return task, task.Status, nil
	}
}

// WaitForServerState ожидает достижения сервером определенного состояния
func WaitForServerState(ctx context.Context, serversService ServerService, serverID int, config *WaiterConfig) (interface{}, error) {
	log.Printf("[DEBUG] Waiting for server %d to reach target state %v", serverID, config.Target)

	stateConf := &resource.StateChangeConf{
		Pending:      config.Pending,
		Target:       config.Target,
		Refresh:      ServerStateRefreshFunc(ctx, serversService, serverID),
		Timeout:      config.Timeout,
		Delay:        config.Delay,
		MinTimeout:   config.MinTimeout,
		PollInterval: config.PollInterval,
	}

	return stateConf.WaitForStateContext(ctx)
}

// WaitForTaskCompletion ожидает завершения задачи
func WaitForTaskCompletion(ctx context.Context, serversService ServerService, taskID int, config *WaiterConfig) (interface{}, error) {
	log.Printf("[DEBUG] Waiting for task %d to complete", taskID)

	if config == nil {
		config = DefaultTaskWaiterConfig()
	}

	// Устанавливаем стандартные состояния для задач
	config.Pending = []string{"pending", "running", "in_progress"}
	config.Target = []string{"completed", "success"}

	stateConf := &resource.StateChangeConf{
		Pending:      config.Pending,
		Target:       config.Target,
		Refresh:      TaskStateRefreshFunc(ctx, serversService, taskID),
		Timeout:      config.Timeout,
		Delay:        config.Delay,
		MinTimeout:   config.MinTimeout,
		PollInterval: config.PollInterval,
	}

	return stateConf.WaitForStateContext(ctx)
}

// WaitForServerDeletion ожидает удаления сервера
func WaitForServerDeletion(ctx context.Context, serversService ServerService, serverID int, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for server %d to be deleted", serverID)

	config := &WaiterConfig{
		Pending:      []string{"active", "stopped", "rebooting", "maintenance"},
		Target:       []string{"DELETED"},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		MinTimeout:   5 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err := WaitForServerState(ctx, serversService, serverID, config)
	return err
}

// WaitForServerToBeActive ожидает активации сервера
func WaitForServerToBeActive(ctx context.Context, serversService ServerService, serverID int, timeout time.Duration) (interface{}, error) {
	log.Printf("[DEBUG] Waiting for server %d to become active", serverID)

	config := DefaultServerWaiterConfig()
	config.Pending = []string{"installing", "rebooting", "maintenance", "stopped"}
	config.Target = []string{"active"}
	config.Timeout = timeout

	return WaitForServerState(ctx, serversService, serverID, config)
}

// WaitForServerToBeStopped ожидает остановки сервера
func WaitForServerToBeStopped(ctx context.Context, serversService ServerService, serverID int, timeout time.Duration) (interface{}, error) {
	log.Printf("[DEBUG] Waiting for server %d to be stopped", serverID)

	config := DefaultServerWaiterConfig()
	config.Pending = []string{"active", "rebooting", "stopping"}
	config.Target = []string{"stopped"}
	config.Timeout = timeout

	return WaitForServerState(ctx, serversService, serverID, config)
}

// WaitForPowerAction ожидает завершения операции управления питанием
func WaitForPowerAction(ctx context.Context, serversService ServerService, serverID int, taskID int, action string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for power action %s on server %d (task %d) to complete", action, serverID, taskID)

	// Сначала ждем завершения задачи
	taskConfig := DefaultTaskWaiterConfig()
	taskConfig.Timeout = timeout

	_, err := WaitForTaskCompletion(ctx, serversService, taskID, taskConfig)
	if err != nil {
		return fmt.Errorf("error waiting for power action task %d: %w", taskID, err)
	}

	// Затем ждем достижения сервером нужного состояния
	var targetStatus []string
	switch action {
	case "start":
		targetStatus = []string{"active"}
	case "stop":
		targetStatus = []string{"stopped"}
	case "restart", "power_cycle":
		targetStatus = []string{"active"}
	default:
		// Для неизвестных действий просто ждем завершения задачи
		return nil
	}

	if len(targetStatus) > 0 {
		_, err = WaitForServerToBeInStates(ctx, serversService, serverID, targetStatus, timeout)
		if err != nil {
			return fmt.Errorf("error waiting for server %d to reach target status after %s: %w", serverID, action, err)
		}
	}

	return nil
}

// WaitForServerToBeInStates ожидает достижения сервером одного из указанных состояний
func WaitForServerToBeInStates(ctx context.Context, serversService ServerService, serverID int, targetStates []string, timeout time.Duration) (interface{}, error) {
	config := DefaultServerWaiterConfig()
	config.Pending = []string{"installing", "rebooting", "maintenance", "stopping", "starting"}
	config.Target = targetStates
	config.Timeout = timeout

	return WaitForServerState(ctx, serversService, serverID, config)
}

// Вспомогательные функции для определения типа ошибки

func isServerNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errorStr := strings.ToLower(err.Error())
	return strings.Contains(errorStr, "not found") ||
		strings.Contains(errorStr, "404") ||
		strings.Contains(errorStr, "server not found")
}

func isTaskNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errorStr := strings.ToLower(err.Error())
	return strings.Contains(errorStr, "not found") ||
		strings.Contains(errorStr, "404") ||
		strings.Contains(errorStr, "task not found")
}

// ServerService интерфейс для работы с серверами (для тестирования)
type ServerService interface {
	GetServer(ctx context.Context, serverID int) (*DedicatedServer, error)
	GetTask(ctx context.Context, taskID int) (*ServerTaskStatus, error)
}

// Типы для совместимости с основным кодом
type DedicatedServer struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	// Другие поля...
}

type ServerTaskStatus struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Error    string `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
	// Другие поля...
}

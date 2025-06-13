package selectel

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/selectel/go-selvpcclient/v4/selvpcclient"
)

var (
	cfgSingletone *Config
	once          sync.Once
)

// Config contains all available configuration options.
type Config struct {
	Region    string
	ProjectID string

	Context        context.Context
	AuthURL        string
	AuthRegion     string
	Username       string
	Password       string
	UserDomainName string
	DomainName     string
	clientsCache   map[string]*selvpcclient.Client

	// Dedicated servers configuration
	ServersToken  string
	serversClient *ServersClient
	lock          sync.Mutex
}

func getConfig(d *schema.ResourceData) (*Config, diag.Diagnostics) {
	// Отключаем singleton для тестирования - создаем конфиг каждый раз заново
	cfgSingletone = &Config{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		DomainName: d.Get("domain_name").(string),
		AuthURL:    d.Get("auth_url").(string),
		AuthRegion: d.Get("auth_region").(string),
	}
	if v, ok := d.GetOk("user_domain_name"); ok {
		cfgSingletone.UserDomainName = v.(string)
	}
	if v, ok := d.GetOk("project_id"); ok {
		cfgSingletone.ProjectID = v.(string)
	}
	if v, ok := d.GetOk("region"); ok {
		cfgSingletone.Region = v.(string)
	}
	// Dedicated servers token (optional)
	if v, ok := d.GetOk("servers_token"); ok {
		cfgSingletone.ServersToken = v.(string)
	}

	return cfgSingletone, nil
}

func (c *Config) GetSelVPCClient() (*selvpcclient.Client, error) {
	return c.GetSelVPCClientWithProjectScope("")
}

func (c *Config) GetSelVPCClientWithProjectScope(projectID string) (*selvpcclient.Client, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	clientsCacheKey := fmt.Sprintf("client_%s", projectID)

	if client, ok := c.clientsCache[clientsCacheKey]; ok {
		return client, nil
	}

	opts := &selvpcclient.ClientOptions{
		DomainName:     c.DomainName,
		Username:       c.Username,
		Password:       c.Password,
		ProjectID:      projectID,
		AuthURL:        c.AuthURL,
		AuthRegion:     c.AuthRegion,
		UserDomainName: c.UserDomainName,
	}

	client, err := selvpcclient.NewClient(opts)
	if err != nil {
		return nil, err
	}

	if c.clientsCache == nil {
		c.clientsCache = map[string]*selvpcclient.Client{}
	}

	c.clientsCache[clientsCacheKey] = client

	return client, nil
}

// GetServersClient возвращает клиент для работы с выделенными серверами
func (c *Config) GetServersClient() (*ServersClient, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	fmt.Fprintf(os.Stderr, "*** GetServersClient() CALLED ***\n")
	log.Printf("[INFO] GetServersClient called, current ServersToken: %s", c.ServersToken)

	// Если клиент уже создан, возвращаем его
	if c.serversClient != nil {
		fmt.Fprintf(os.Stderr, "*** REUSING CACHED SERVERS CLIENT ***\n")
		return c.serversClient, nil
	}

	var token string

	// Приоритет: используем servers_token если указан, иначе получаем токен через Keystone
	if c.ServersToken != "" {
		log.Printf("[DEBUG] Using provided servers_token (length: %d)", len(c.ServersToken))
		log.Printf("[DEBUG] Provided token: %s", c.ServersToken)
		token = c.ServersToken
	} else {
		log.Printf("[INFO] No servers_token provided, attempting Keystone authentication")
		log.Printf("[INFO] ATTEMPTING KEYSTONE AUTH WITH USER: %s", c.Username)
		fmt.Fprintf(os.Stderr, "*** KEYSTONE AUTH ATTEMPT FOR USER: %s ***\n", c.Username)
		// Получаем токен через автоматическую аутентификацию Keystone
		selvpcClient, err := c.GetSelVPCClient()
		if err != nil {
			return nil, fmt.Errorf("failed to get selvpc client for servers authentication: %w", err)
		}

		keystoneToken := selvpcClient.GetXAuthToken()
		if keystoneToken == "" {
			return nil, fmt.Errorf("failed to obtain authentication token via Keystone")
		}

		log.Printf("[INFO] Successfully obtained Keystone token (length: %d)", len(keystoneToken))
		log.Printf("[INFO] Keystone token: %s", keystoneToken)
		token = keystoneToken
	}

	// Создаем клиент для выделенных серверов
	opts := &ServersClientOptions{
		Token:   token,
		Context: c.Context,
	}

	client, err := NewServersClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create servers client: %w", err)
	}

	c.serversClient = client
	return client, nil
}

// GetServersService возвращает сервис для работы с выделенными серверами
func (c *Config) GetServersService() (*ServersService, error) {
	log.Printf("[INFO] GetServersService() called")
	client, err := c.GetServersClient()
	if err != nil {
		log.Printf("[ERROR] GetServersClient failed: %v", err)
		return nil, err
	}

	log.Printf("[INFO] GetServersService: client created successfully")
	return NewServersService(client), nil
}

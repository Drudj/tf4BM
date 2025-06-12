# Интеграция автоматической аутентификации с Terraform провайдером Selectel

## Обзор изменений

Реализована интеграция автоматической аутентификации для выделенных серверов Selectel с использованием существующей системы аутентификации Keystone из оригинального Terraform провайдера.

## Что было сделано

### 1. Обновлен механизм получения токена

**Файл:** `terraform-provider-selectel/selectel/config.go`

Модифицирован метод `GetServersClient()` для поддержки автоматической аутентификации:

```go
// Приоритет: используем servers_token если указан, иначе получаем токен через Keystone
if c.ServersToken != "" {
    token = c.ServersToken
} else {
    // Получаем токен через автоматическую аутентификацию Keystone
    selvpcClient, err := c.GetSelVPCClient()
    if err != nil {
        return nil, fmt.Errorf("failed to get selvpc client for servers authentication: %w", err)
    }
    
    keystoneToken := selvpcClient.GetXAuthToken()
    if keystoneToken == "" {
        return nil, fmt.Errorf("failed to obtain authentication token via Keystone")
    }
    
    token = keystoneToken
}
```

### 2. Обновлено описание поля в провайдере

**Файл:** `terraform-provider-selectel/selectel/provider.go`

Поле `servers_token` теперь опциональное:

```go
"servers_token": {
    Type:        schema.TypeString,
    Optional:    true,
    DefaultFunc: schema.EnvDefaultFunc("SEL_SERVERS_TOKEN", nil),
    Description: "Bearer token for dedicated servers API access. If not provided, will use Keystone authentication token.",
    Sensitive:   true,
},
```

### 3. Добавлены недостающие типы данных

**Файл:** `terraform-provider-selectel/selectel/servers_models.go`

Добавлены структуры:
- `ServerService` - с полем `State` для фильтрации активных сервисов
- `DedicatedServerCreateBilling` - с расширенными полями для биллинг API
- `DedicatedServerCreateResponse` и `DedicatedServerCreateResult` - для обработки ответов

### 4. Исправлены ошибки компиляции

Добавлена функция `GenerateAdvancedPartitionsConfig()` для генерации конфигурации партиций.

## Принцип работы

1. **Приоритет токенов:**
   - Если указан `servers_token` - используется он
   - Если `servers_token` не указан - автоматически получается токен через Keystone

2. **Совместимость:**
   - Полная обратная совместимость со старым способом (через `servers_token`)
   - Новый способ позволяет использовать единую аутентификацию для всех сервисов Selectel

3. **Аутентификация через Keystone:**
   - Использует тот же механизм, что и остальные сервисы Selectel (MKS, DBaaS, etc.)
   - Токен получается автоматически при создании `selvpcClient`
   - Токен извлекается через `selvpcClient.GetXAuthToken()`

## Использование

### Способ 1: Через переменные окружения (новый, рекомендуемый)

```bash
export OS_DOMAIN_NAME="123456"
export OS_USERNAME="your_service_user"
export OS_PASSWORD="your_password"
export OS_AUTH_URL="https://cloud.api.selcloud.ru/identity/v3/"
export OS_REGION_NAME="pool"
```

```hcl
provider "selectel" {
  domain_name = var.selectel_account_id
  username    = var.selectel_username
  password    = var.selectel_password
  auth_url    = "https://cloud.api.selcloud.ru/identity/v3/"
  auth_region = "pool"
  
  # servers_token теперь не обязателен!
}
```

### Способ 2: Старый способ (для обратной совместимости)

```bash
export SEL_SERVERS_TOKEN="your_servers_token"
```

```hcl
provider "selectel" {
  # ... остальные параметры ...
  servers_token = var.selectel_servers_token
}
```

## Тестирование

1. **Скопируйте файл переменных:**
   ```bash
   cp test_auth.tfvars.example test_auth.tfvars
   ```

2. **Заполните файл `test_auth.tfvars`:**
   ```hcl
   selectel_account_id = "123456"
   selectel_username   = "your_service_user"
   selectel_password   = "your_password"
   ```

3. **Запустите тест:**
   ```bash
   terraform init
   terraform plan -var-file="test_auth.tfvars"
   ```

## Преимущества интеграции

1. **Единая аутентификация:** Больше не нужно запрашивать отдельный токен для выделенных серверов
2. **Автоматическое обновление токенов:** Keystone токены обновляются автоматически
3. **Согласованность:** Тот же механизм аутентификации, что используется в MKS, DBaaS и других сервисах
4. **Безопасность:** Токены не нужно хранить в переменных окружения длительное время
5. **Простота использования:** Меньше настроек для пользователей

## Документация Selectel

- [Аутентификация в Terraform провайдере](https://docs.selectel.ru/terraform/selectel-provider-reference/authentication/)
- [Быстрый старт](https://docs.selectel.ru/terraform/quickstart/)

## Техническая информация

- **go-selvpcclient версия:** v4
- **Метод получения токена:** `selvpcClient.GetXAuthToken()`
- **Endpoint аутентификации:** `https://cloud.api.selcloud.ru/identity/v3/`
- **Регион аутентификации:** `pool` (по умолчанию) 
# Базовый пример создания выделенного сервера Selectel

Этот пример показывает, как создать базовый выделенный сервер в Selectel с помощью Terraform провайдера.

## Что создается

- **Один выделенный сервер** с операционной системой Ubuntu
- **Публичная сеть** с полосой пропускания 1000 Mbps  
- **SSH доступ** через публичные ключи
- **Теги** для организации ресурсов

## Предварительные требования

1. **Terraform** версии 1.5 или выше
2. **Selectel IAM токен** для доступа к API
3. **UUID проекта** Selectel
4. **SSH публичный ключ** для доступа к серверу

## Настройка

### 1. Установка переменных окружения

```bash
export SELECTEL_TOKEN="your-iam-token"
export SELECTEL_PROJECT_ID="your-project-uuid"
```

### 2. Подготовка конфигурации

Скопируйте файл с примером переменных:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Отредактируйте `terraform.tfvars` и укажите ваши значения:

```hcl
project_uuid = "f75e16974d78419488ff638c19d799a4"
ssh_public_keys = [
  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7... your-key@hostname"
]
server_name = "my-web-server"
environment = "production"
```

### 3. Инициализация и применение

```bash
# Инициализация Terraform
terraform init

# Просмотр плана
terraform plan

# Применение конфигурации
terraform apply
```

## Переменные

| Переменная | Описание | Тип | По умолчанию | Обязательная |
|------------|----------|-----|--------------|-------------|
| `project_uuid` | UUID проекта Selectel | `string` | - | ✅ |
| `ssh_public_keys` | Список SSH публичных ключей | `list(string)` | `[]` | ✅ |
| `root_password` | Пароль root пользователя | `string` | `""` | ❌ |
| `server_name` | Имя сервера | `string` | `"web-server-basic"` | ❌ |
| `environment` | Окружение | `string` | `"production"` | ❌ |

## Выходные данные

После успешного создания сервера вы получите:

```hcl
server_info = {
  uuid         = "server-uuid"
  name         = "web-server-basic"
  status       = "active"
  power_status = "on"
  ip_addresses = [
    {
      address = "192.168.1.100"
      type    = "public"
      version = "ipv4"
    }
  ]
}
```

## Подключение к серверу

После создания сервера вы можете подключиться по SSH:

```bash
ssh root@<server-ip-address>
```

IP адрес можно получить из выходных данных Terraform:

```bash
terraform output server_info
```

## Управление ресурсами

### Просмотр состояния
```bash
terraform show
```

### Обновление конфигурации
Измените переменные в `terraform.tfvars` и выполните:
```bash
terraform apply
```

### Удаление ресурсов
```bash
terraform destroy
```

## Безопасность

- SSH ключи являются основным способом аутентификации
- Пароль root опционален и должен быть сложным (минимум 8 символов)
- Все чувствительные данные помечены как `sensitive`
- Рекомендуется использовать переменные окружения для токенов

## Примечания

- Создание сервера может занять 5-15 минут
- Статус сервера можно отслеживать через `terraform refresh`
- IP адреса назначаются автоматически при создании
- Теги помогают организовать и отслеживать ресурсы 
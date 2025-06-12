# Terraform Provider для Selectel Выделенных Серверов

Terraform провайдер для управления выделенными серверами Selectel через API.

## Возможности

- ✅ Создание выделенных серверов
- ✅ Настройка RAID конфигураций
- ✅ Управление разделами диска
- ✅ Поддержка различных операционных систем
- ✅ Интеграция с Selectel Cloud API

## Требования

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19 (для сборки из исходников)
- Аккаунт Selectel с доступом к API выделенных серверов

## Установка

### Из исходников

1. Клонируйте репозиторий:
```bash
git clone https://github.com/Drudj/tf4BM.git
cd tf4BM
```

2. Соберите провайдер:
```bash
cd terraform-provider-selectel
go build -o ../terraform-provider-selectel_v1.0.0
cd ..
```

3. Настройте development overrides в `~/.terraformrc`:
```hcl
provider_installation {
  dev_overrides {
    "selectel/selectel" = "/path/to/tf4BM"
  }
  direct {}
}
```

## Настройка

### Получение токенов

1. **Токен для Cloud API**: Получите в панели управления Selectel Cloud
2. **Токен для Servers API**: Получите в разделе "Выделенные серверы"

### Переменные окружения

```bash
export OS_AUTH_URL="https://cloud.api.selcloud.ru/identity/v3/"
export OS_DOMAIN_NAME="YOUR_DOMAIN_ID"
export OS_USERNAME="YOUR_USERNAME"
export OS_PASSWORD="YOUR_PASSWORD"
export OS_REGION_NAME="pool"
export SEL_SERVERS_TOKEN="YOUR_SERVERS_TOKEN"
```

## Использование

### Базовый пример

```hcl
terraform {
  required_providers {
    selectel = {
      source = "selectel/selectel"
    }
  }
}

provider "selectel" {
  auth_url      = "https://cloud.api.selcloud.ru/identity/v3/"
  auth_region   = "pool"
  domain_name   = var.selectel_domain_name
  username      = var.selectel_username
  password      = var.selectel_password
  servers_token = var.selectel_servers_token
}

resource "selectel_dedicated_server_v1" "example" {
  name        = "my-server"
  config_id   = 23  # CL23-SSD
  location_id = 2   # MSK-2
  os_id       = 1   # Debian
  root_size   = 50
  swap_size   = 5
  raid_type   = "RAID1"
  ssh_keys    = []
}
```

### Полный пример

Смотрите файлы в папке `examples/` для полного примера с переменными.

## Ресурсы

### selectel_dedicated_server_v1

Создает и управляет выделенным сервером.

#### Аргументы

- `name` (string, обязательный) - Имя сервера
- `config_id` (number, обязательный) - ID конфигурации сервера
- `location_id` (number, обязательный) - ID локации
- `os_id` (number, обязательный) - ID операционной системы
- `root_size` (number, опциональный) - Размер корневого раздела в GB (по умолчанию: 50)
- `swap_size` (number, опциональный) - Размер swap раздела в GB (по умолчанию: 5)
- `raid_type` (string, опциональный) - Тип RAID (по умолчанию: "RAID1")
- `ssh_keys` (list(string), опциональный) - Список SSH ключей
- `enable_backup` (bool, опциональный) - Включить резервное копирование
- `enable_ipmi` (bool, опциональный) - Включить IPMI

#### Атрибуты

- `id` - UUID созданного сервера
- `status` - Статус сервера
- `cpu` - Информация о процессоре
- `ram` - Информация о памяти
- `storage` - Информация о хранилище
- `network` - Сетевая информация
- `location` - Информация о локации
- `os` - Информация об ОС
- `price` - Информация о стоимости
- `created_at` - Время создания
- `updated_at` - Время последнего обновления

## Конфигурации серверов

| Config ID | Название | Описание |
|-----------|----------|----------|
| 23 | CL23-SSD | Intel Xeon, SSD диски |

## Локации

| Location ID | Название | Описание |
|-------------|----------|----------|
| 2 | MSK-2 | Москва, ЦОД 2 |

## Операционные системы

| OS ID | Название |
|-------|----------|
| 1 | Debian |

## Разработка

### Структура проекта

```
tf4BM/
├── terraform-provider-selectel/    # Исходный код провайдера
│   ├── selectel/                   # Основные файлы провайдера
│   ├── main.go                     # Точка входа
│   └── go.mod                      # Go модуль
├── examples/                       # Примеры использования
│   ├── main.tf                     # Основная конфигурация
│   ├── variables.tf                # Переменные
│   └── terraform.tfvars.example    # Пример переменных
├── private/                        # Приватные конфигурации (игнорируется git)
└── README.md                       # Документация
```

### Сборка

```bash
cd terraform-provider-selectel
go build -o ../terraform-provider-selectel_v1.0.0
```

### Тестирование

```bash
# Создайте файл private/test.tf с вашими данными
cd private
terraform init
terraform plan
terraform apply
```

## Устранение неполадок

### Частые проблемы

1. **"Provider development overrides are in effect"** - Это нормальное предупреждение при использовании локальной сборки
2. **"HTTP 401"** - Проверьте правильность токенов и учетных данных
3. **"HTTP 400: partitions_config"** - Убедитесь, что используете правильную структуру конфигурации

### Логирование

Для включения подробного логирования:

```bash
export TF_LOG=DEBUG
terraform apply
```

## Лицензия

MIT License

## Поддержка

Для вопросов и предложений создавайте issues в GitHub репозитории.

## Changelog

### v1.0.0
- Первый рабочий релиз
- Поддержка создания выделенных серверов
- Интеграция с Selectel API
- Поддержка RAID конфигураций 
# Terraform Provider для Selectel Выделенных Серверов

Terraform провайдер для управления выделенными серверами Selectel через API.

## 🚀 Возможности

- ✅ Создание выделенных серверов
- ✅ Настройка RAID конфигураций
- ✅ Управление разделами диска
- ✅ Поддержка различных операционных систем
- ✅ Интеграция с Selectel Cloud API
- ✅ Безопасная работа с переменными окружения

## 📋 Требования

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19 (для сборки из исходников)
- Аккаунт Selectel с доступом к API выделенных серверов

## 🔧 Установка

### Из исходников

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/terraform-provider-selectel-dedicated.git
cd terraform-provider-selectel-dedicated
```

2. Настройте переменные окружения:
```bash
cp test.env.example test.env
# Отредактируйте test.env своими данными
source test.env
```

3. Соберите провайдер:
```bash
cd terraform-provider-selectel
go build -o ../terraform-provider-selectel_v1.0.0
cd ..
```

4. Настройте development overrides в `~/.terraformrc`:
```hcl
provider_installation {
  dev_overrides {
    "selectel/selectel" = "/path/to/your/project"
  }
  direct {}
}
```

## 🔐 Безопасная настройка

⚠️ **ВАЖНО**: Следуйте инструкциям в [SETUP.md](SETUP.md) для безопасной настройки

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

## 📚 Использование

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

## 📖 Ресурсы

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

## 📊 Справочная информация

### Конфигурации серверов

| Config ID | Название | Описание |
|-----------|----------|----------|
| 23 | CL23-SSD | Intel Xeon, SSD диски |

### Локации

| Location ID | Название | Описание |
|-------------|----------|----------|
| 2 | MSK-2 | Москва, ЦОД 2 |

### Операционные системы

| OS ID | Название |
|-------|----------|
| 1 | Debian |

## 🛠️ Разработка

### Структура проекта

```
terraform-provider-selectel-dedicated/
├── terraform-provider-selectel/    # Исходный код провайдера
│   ├── selectel/                   # Основные файлы провайдера
│   │   ├── provider.go             # Конфигурация провайдера
│   │   ├── config.go               # Конфигурация клиентов
│   │   ├── servers_*.go            # Логика для выделенных серверов
│   │   └── resource_selectel_dedicated_server_v1.go
│   ├── main.go                     # Точка входа
│   └── go.mod                      # Go модуль
├── examples/                       # Примеры использования
├── private/                        # Приватные конфигурации (в .gitignore)
├── test.env.example                # Пример переменных окружения
├── SETUP.md                        # Инструкции по безопасной настройке
└── README.md                       # Документация
```

### Сборка

```bash
cd terraform-provider-selectel
go build -o ../terraform-provider-selectel_v1.0.0
```

### Тестирование

```bash
# Настройте переменные окружения
source test.env
cd examples
terraform init
terraform plan
terraform apply
```

## 🐛 Устранение неполадок

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

## 🔒 Безопасность

- Все секретные данные должны храниться в переменных окружения
- Файл `test.env` добавлен в `.gitignore`
- Перед коммитом выполняйте проверку: `grep -r "gAAAAA\|OS_PASSWORD\|SEL_SERVERS_TOKEN" . --exclude-dir=.git --exclude="*.md" --exclude="*.example"`

## 📄 Лицензия

MIT License

## 🆘 Поддержка

Для вопросов и предложений создавайте issues в GitHub репозитории.

## 📝 Changelog

### v1.0.0
- Первый рабочий релиз
- Поддержка создания выделенных серверов
- Интеграция с Selectel API
- Поддержка RAID конфигураций
- Безопасная работа с переменными окружения 
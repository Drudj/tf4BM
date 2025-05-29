# Terraform Provider для Selectel Bare Metal Серверов

[![CI](https://github.com/Drudj/tf_for_BareMetal/actions/workflows/ci.yml/badge.svg)](https://github.com/Drudj/tf_for_BareMetal/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/selectel/terraform-provider-selectel-baremetal)](https://goreportcard.com/report/github.com/selectel/terraform-provider-selectel-baremetal)

Terraform провайдер для управления выделенными серверами (bare metal) Selectel через Infrastructure as Code.

## 🚀 Статус разработки

✅ **Готов к использованию** - Основная функциональность полностью реализована и протестирована.

### ✅ Реализованные возможности

- **Управление выделенными серверами**: Создание, обновление и удаление физических серверов
- **Сетевые настройки**: Публичные и приватные сети, настройка полосы пропускания, VLAN
- **Операционные системы**: Автоматическая установка с SSH ключами, паролями и cloud-init
- **Data Sources**: Получение информации о локациях, услугах, OS шаблонах и тарифах
- **Теги и метаданные**: Полная поддержка пользовательских тегов для организации ресурсов
- **Импорт существующих серверов**: Интеграция с уже созданными ресурсами

## 📋 Требования

- **Terraform** >= 1.5.0
- **Go** >= 1.23 (только для разработки)
- **Selectel аккаунт** с IAM токеном и доступом к API выделенных серверов

## 🔧 Быстрый старт

### 1. Настройка аутентификации

```bash
export SELECTEL_TOKEN="your-iam-token"
export SELECTEL_PROJECT_ID="your-project-uuid"
```

### 2. Базовая конфигурация

```hcl
terraform {
  required_version = ">= 1.5"
  required_providers {
    selectel-baremetal = {
      source  = "selectel/selectel-baremetal"
      version = "~> 0.1"
    }
  }
}

provider "selectel-baremetal" {
  # Конфигурация берется из переменных окружения
}

# Получение доступных ресурсов
data "selectel_baremetal_locations" "all" {}
data "selectel_baremetal_services" "all" {}
data "selectel_baremetal_os_templates" "ubuntu" {}

# Создание базового сервера
resource "selectel_baremetal_server" "web" {
  name            = "my-web-server"
  service_uuid    = data.selectel_baremetal_services.all.services[0].uuid
  location_uuid   = data.selectel_baremetal_locations.all.locations[0].uuid
  price_plan_uuid = data.selectel_baremetal_services.all.services[0].price_plans[0].uuid
  project_uuid    = var.project_uuid

  network {
    type      = "public"
    bandwidth = 1000
  }

  os {
    template_uuid = data.selectel_baremetal_os_templates.ubuntu.templates[0].uuid
    ssh_keys      = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7..."]
  }

  tags = {
    Environment = "production"
    Application = "web"
  }
}
```

### 3. Применение конфигурации

```bash
terraform init
terraform plan
terraform apply
```

## 📖 Примеры использования

В директории [`examples/`](./examples/) доступны готовые примеры:

### 🔥 [Базовый сервер](./examples/basic-server/)
Простой пример создания одного сервера с минимальными настройками.

```bash
cd examples/basic-server
cp terraform.tfvars.example terraform.tfvars
# Отредактируйте terraform.tfvars
terraform init && terraform apply
```

### 🚀 [Множественные серверы](./examples/multiple-servers/)
Создание нескольких серверов с разными ролями (web, database, API).

### 🌐 [Расширенные сетевые настройки](./examples/with-networking/)
Серверы с различными типами сетей, VLAN и cloud-init скриптами.

### ⚙️ [Настраиваемый сервер](./examples/custom-server/)
Полностью настраиваемая конфигурация с выбором локации, ОС и сетевых параметров.

## 📚 Ресурсы провайдера

### Основные ресурсы

| Ресурс | Описание |
|--------|----------|
| `selectel_baremetal_server` | Управление выделенными серверами |

### Data Sources

| Data Source | Описание |
|-------------|----------|
| `selectel_baremetal_locations` | Список доступных локаций |
| `selectel_baremetal_location` | Информация о конкретной локации |
| `selectel_baremetal_services` | Каталог доступных услуг |
| `selectel_baremetal_service` | Информация о конкретной услуге |
| `selectel_baremetal_os_templates` | Доступные шаблоны ОС |
| `selectel_baremetal_os_template` | Информация о конкретном шаблоне |
| `selectel_baremetal_price_plans` | Тарифные планы |

## 🔐 Аутентификация

### Переменные окружения (рекомендуется)

```bash
export SELECTEL_TOKEN="your-iam-token"
export SELECTEL_PROJECT_ID="your-project-uuid"
```

### Конфигурация провайдера

```hcl
provider "selectel-baremetal" {
  token      = "your-iam-token"
  project_id = "your-project-uuid"
  endpoint   = "https://api.selectel.ru/servers/v2"  # опционально
}
```

## 🛠️ Разработка

### Настройка окружения разработки

```bash
git clone https://github.com/Drudj/tf_for_BareMetal.git
cd tf_for_BareMetal

# Установка зависимостей
go mod download

# Установка инструментов разработки
make dev-setup
```

### Основные команды

```bash
# Сборка провайдера
make build

# Локальная установка
make install

# Запуск тестов
make test

# Полная проверка кода
make check

# Форматирование кода
make fmt

# Линтинг кода
make lint
```

### Структура проекта

```
├── cmd/terraform-provider-selectel-baremetal/  # Основной исполняемый файл
├── internal/
│   ├── client/         # HTTP клиент для API
│   ├── datasources/    # Terraform data sources
│   ├── models/         # Модели данных API
│   ├── provider/       # Конфигурация провайдера
│   └── resources/      # Terraform ресурсы
├── examples/           # Примеры использования
├── docs/              # Документация
└── Makefile           # Команды сборки и тестирования
```

## 🧪 Тестирование

### Unit тесты

```bash
make test
```

### Тесты с покрытием

```bash
make test-coverage
```

### Acceptance тесты

```bash
export SELECTEL_TOKEN="your-token"
export SELECTEL_PROJECT_ID="your-project"
make testacc
```

## 📖 API Документация

- [Selectel Dedicated Servers API](https://docs.selectel.ru/api/dedicated/)
- [Авторизация в API](https://docs.selectel.ru/api/authorization/)
- [Управление серверами](https://docs.selectel.ru/servers-and-infrastructure/dedicated/)

## 🆘 Поддержка

- **Issues**: [GitHub Issues](https://github.com/Drudj/tf_for_BareMetal/issues) для багов и запросов функций
- **Документация**: [Selectel API Docs](https://docs.selectel.ru/)
- **Техподдержка**: [Selectel Support](https://selectel.ru/support/)

## 📝 Лицензия

Этот проект распространяется под лицензией Apache 2.0. См. файл [LICENSE](LICENSE) для деталей.

## 🤝 Вклад в развитие

Мы приветствуем участие сообщества! Пожалуйста, ознакомьтесь с [CONTRIBUTING.md](CONTRIBUTING.md) для получения инструкций по разработке.

### Как внести вклад

1. Fork репозитория
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request 
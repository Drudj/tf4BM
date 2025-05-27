# Terraform Provider for Selectel Bare Metal Servers

[![CI](https://github.com/selectel/terraform-provider-selectel-baremetal/actions/workflows/ci.yml/badge.svg)](https://github.com/selectel/terraform-provider-selectel-baremetal/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/selectel/terraform-provider-selectel-baremetal)](https://goreportcard.com/report/github.com/selectel/terraform-provider-selectel-baremetal)
[![codecov](https://codecov.io/gh/selectel/terraform-provider-selectel-baremetal/branch/main/graph/badge.svg)](https://codecov.io/gh/selectel/terraform-provider-selectel-baremetal)

Terraform provider для управления выделенными серверами (bare metal) Selectel через Infrastructure as Code.

## Статус разработки

✅ **Основная функциональность реализована** - Провайдер готов для базового использования.

### Реализованные возможности

- ✅ **Управление выделенными серверами**: Создание, настройка и удаление физических серверов
- ✅ **Сетевые настройки**: Настройка публичных и приватных сетей, VLAN
- ✅ **Управление ОС**: Установка операционных систем с SSH ключами и паролями
- ✅ **Data Sources**: Получение информации о локациях, услугах, OS шаблонах и тарифах
- ✅ **Интеграция**: Полная совместимость с Terraform Framework
- ✅ **Теги и метаданные**: Поддержка пользовательских тегов

### Планируемые возможности

- 🚧 **Управление питанием**: Включение, выключение, перезагрузка серверов  
- 🚧 **Расширенные сетевые настройки**: Дополнительные опции сетевой конфигурации
- 🚧 **Мониторинг**: Отслеживание состояния серверов и задач
- 🚧 **Интеграционные тесты**: Полное тестирование с реальным API

## Требования

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (для разработки)
- Аккаунт Selectel с доступом к API выделенных серверов

## Использование

### Базовая конфигурация

```hcl
terraform {
  required_providers {
    selectel = {
      source  = "selectel/selectel-baremetal"
      version = "~> 0.1"
    }
  }
}

provider "selectel" {
  token      = var.selectel_token
  project_id = var.project_id
}

# Получение доступных локаций
data "selectel_baremetal_location" "msk" {
  name = "Moscow"
}

# Получение доступных конфигураций серверов
data "selectel_baremetal_service" "server" {
  name = "Intel Xeon E-2288G"
}

# Получение OS шаблонов
data "selectel_baremetal_os_template" "ubuntu" {
  name = "Ubuntu 22.04 LTS"
}

# Создание выделенного сервера
resource "selectel_baremetal_server" "example" {
  name              = "my-server"
  service_uuid      = data.selectel_baremetal_service.server.uuid
  location_uuid     = data.selectel_baremetal_location.msk.uuid
  price_plan_uuid   = data.selectel_baremetal_service.server.price_plans[0].uuid
  project_uuid      = var.project_id
  
  network {
    type      = "public"
    bandwidth = 1000
  }
  
  os {
    template_uuid = data.selectel_baremetal_os_template.ubuntu.uuid
    ssh_keys      = [var.ssh_key]
  }
  
  tags = {
    Environment = "production"
    Team        = "infrastructure"
  }
}
```

### Примеры

Больше примеров доступно в директории [examples/](./examples/):

- [Базовый сервер](./examples/basic-server/)
- [Кастомная конфигурация](./examples/custom-server/)
- [Несколько серверов](./examples/multiple-servers/)
- [Сервер с сетевыми настройками](./examples/with-networking/)

## Ресурсы

- `selectel_baremetal_server` - Основной ресурс выделенного сервера
- `selectel_baremetal_server_power` - Управление питанием сервера
- `selectel_baremetal_server_network` - Сетевые настройки сервера
- `selectel_baremetal_server_os` - Управление ОС сервера

## Источники данных

- `selectel_baremetal_locations` - Список доступных локаций
- `selectel_baremetal_services` - Каталог услуг и конфигураций серверов
- `selectel_baremetal_os_templates` - Доступные шаблоны операционных систем
- `selectel_baremetal_price_plans` - Тарифные планы

## Разработка

### Настройка окружения

```bash
# Клонирование репозитория
git clone https://github.com/selectel/terraform-provider-selectel-baremetal.git
cd terraform-provider-selectel-baremetal

# Установка зависимостей
make deps

# Установка инструментов разработки
make dev-setup
```

### Сборка

```bash
# Сборка провайдера
make build

# Локальная установка для тестирования
make install
```

### Тестирование

```bash
# Запуск unit тестов
make test

# Запуск тестов с покрытием
make test-coverage

# Запуск acceptance тестов (требует настройки API токенов)
make testacc
```

### Проверка кода

```bash
# Форматирование кода
make fmt

# Линтинг
make lint

# Полная проверка (форматирование + линтинг + тесты)
make check
```

## Аутентификация

Провайдер поддерживает несколько способов аутентификации:

### Переменные окружения

```bash
export SELECTEL_TOKEN="your-api-token"
export SELECTEL_PROJECT_ID="your-project-id"
```

### Конфигурация провайдера

```hcl
provider "selectel" {
  token      = "your-api-token"
  project_id = "your-project-id"
  endpoint   = "https://api.selectel.ru/dedicated/v2"  # опционально
}
```

## Документация

- [Документация ресурсов](./docs/resources/)
- [Документация источников данных](./docs/data-sources/)
- [Руководства пользователя](./docs/guides/)

## Поддержка

- [Issues](https://github.com/selectel/terraform-provider-selectel-baremetal/issues) - сообщения об ошибках и запросы функций
- [Discussions](https://github.com/selectel/terraform-provider-selectel-baremetal/discussions) - вопросы и обсуждения
- [Selectel Support](https://selectel.ru/support/) - техническая поддержка Selectel

## Лицензия

Этот проект лицензирован под [Mozilla Public License 2.0](LICENSE).

## Участие в разработке

Мы приветствуем участие в разработке! Пожалуйста, ознакомьтесь с [CONTRIBUTING.md](CONTRIBUTING.md) для получения информации о том, как внести свой вклад.

### Авторы

- Команда разработки Selectel
- Сообщество участников

## Связанные проекты

- [terraform-provider-selectel](https://github.com/selectel/terraform-provider-selectel) - Основной провайдер Selectel для облачных ресурсов
- [Selectel API Documentation](https://docs.selectel.ru/api/dedicated/) - Документация API выделенных серверов 
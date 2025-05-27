# Базовый пример использования Terraform Provider для Selectel Bare Metal

Этот пример демонстрирует базовое использование terraform провайдера для получения информации о локациях Selectel.

## Требования

- Terraform >= 1.0
- API токен Selectel
- ID проекта Selectel

## Настройка

1. Установите переменные окружения:

```bash
export SELECTEL_TOKEN="your-api-token"
export SELECTEL_PROJECT_ID="your-project-id"
```

2. Инициализируйте Terraform:

```bash
terraform init
```

3. Запустите план:

```bash
terraform plan
```

4. Примените конфигурацию:

```bash
terraform apply
```

## Что делает этот пример

- Получает список всех доступных локаций Selectel
- Находит конкретную локацию по имени ("Moscow")
- Выводит информацию о локациях

## Ожидаемый результат

После выполнения `terraform apply` вы увидите:

- `all_locations` - список всех доступных локаций с их характеристиками
- `moscow_location` - подробную информацию о локации Moscow

## Структура файлов

- `main.tf` - основная конфигурация Terraform
- `README.md` - этот файл с инструкциями 
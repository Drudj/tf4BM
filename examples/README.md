# Примеры использования Terraform Provider для Selectel

Эта папка содержит примеры конфигураций для использования Terraform провайдера Selectel выделенных серверов.

## Быстрый старт

1. Скопируйте пример файла с переменными:
```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Отредактируйте `terraform.tfvars` и заполните своими данными:
```hcl
selectel_domain_name   = "YOUR_DOMAIN_ID"
selectel_username      = "YOUR_USERNAME"
selectel_password      = "YOUR_PASSWORD"
selectel_servers_token = "YOUR_SERVERS_TOKEN"

server_name        = "my-terraform-server"
server_config_id   = 23  # CL23-SSD
server_location_id = 2   # MSK-2
server_os_id       = 1   # Debian
```

3. Инициализируйте и примените конфигурацию:
```bash
terraform init
terraform plan
terraform apply
```

## Файлы

- `main.tf` - Основная конфигурация с ресурсом сервера
- `variables.tf` - Определения переменных
- `terraform.tfvars.example` - Пример файла с переменными
- `README.md` - Эта документация

## Получение данных для аутентификации

### Domain Name (ID домена)
1. Войдите в панель управления Selectel
2. В правом верхнем углу найдите ваш ID домена

### Username и Password
Используйте ваши учетные данные от панели управления Selectel.

### Servers Token
1. Перейдите в раздел "Выделенные серверы"
2. В настройках API создайте новый токен
3. Скопируйте полученный токен

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

## Безопасность

⚠️ **Важно**: Никогда не коммитьте файл `terraform.tfvars` с реальными данными в Git! 

Файл `terraform.tfvars` автоматически игнорируется через `.gitignore`. 
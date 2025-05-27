# Тестирование Terraform провайдера на реальных серверах

Этот каталог содержит инфраструктуру для тестирования нашего Terraform провайдера на реальных серверах Selectel Bare Metal.

## 🔐 Безопасность credentials

**ВАЖНО:** Все приватные данные передаются через переменные окружения и локальные файлы, которые НЕ коммитятся в git.

## 📋 Предварительные требования

1. **Аккаунт Selectel** с доступом к Bare Metal серверам
2. **API токен** с правами на управление серверами
3. **UUID проекта** из панели управления Selectel
4. **SSH ключи** для доступа к серверам
5. **Установленный провайдер** (выполните `make install` в корне проекта)

## 🚀 Пошаговое тестирование

### Шаг 1: Настройка переменных окружения

```bash
# Экспортируйте ваши credentials (НЕ добавляйте в код!)
export SELECTEL_TOKEN="your-real-api-token"
export SELECTEL_PROJECT_ID="your-project-uuid"
export SELECTEL_ENDPOINT="https://api.selectel.ru/dedicated/v2"  # опционально
```

### Шаг 2: Получение UUID для конфигурации

Сначала получите доступные UUID из API:

```bash
cd test-infrastructure

# Инициализация только для discovery
terraform init

# Получение доступных ресурсов (используется get-uuids.tf)
terraform plan -target=data.selectel_baremetal_locations.discovery
terraform apply -target=data.selectel_baremetal_locations.discovery -auto-approve

# Просмотр доступных ресурсов
terraform output discovery_info
terraform output recommended_values
```

### Шаг 3: Создание terraform.tfvars

```bash
# Скопируйте шаблон
cp terraform.tfvars.template terraform.tfvars

# Отредактируйте terraform.tfvars с реальными значениями
# Используйте UUID из предыдущего шага
nano terraform.tfvars
```

Пример заполнения:
```hcl
project_uuid     = "12345678-1234-1234-1234-123456789abc"
service_uuid     = "87654321-4321-4321-4321-cba987654321"
price_plan_uuid  = "11111111-2222-3333-4444-555555555555"
os_template_uuid = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

ssh_keys = [
  "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC... your-actual-key"
]
```

### Шаг 4: Тестирование data sources

```bash
# Тестирование только data sources (без создания серверов)
terraform plan -target=data.selectel_baremetal_locations.all
terraform plan -target=data.selectel_baremetal_services.all
terraform plan -target=data.selectel_baremetal_os_templates.ubuntu

# Применение data sources
terraform apply -target=data.selectel_baremetal_locations.all -auto-approve
terraform apply -target=data.selectel_baremetal_services.all -auto-approve
terraform apply -target=data.selectel_baremetal_os_templates.ubuntu -auto-approve

# Проверка результатов
terraform output available_locations
terraform output available_services_count
terraform output ubuntu_templates_info
```

### Шаг 5: Создание тестового сервера

⚠️ **ВНИМАНИЕ:** Этот шаг создаст реальный сервер и может повлечь расходы!

```bash
# Планирование создания сервера
terraform plan

# Создание сервера (подтвердите расходы!)
terraform apply

# Проверка результатов
terraform output test_server_info
terraform output test_server_network
terraform output connection_info
```

### Шаг 6: Тестирование подключения

```bash
# Получение IP адреса для подключения
terraform output connection_info

# Подключение к серверу (замените IP)
ssh root@YOUR_SERVER_IP

# Проверка сервера
uname -a
df -h
free -m
```

### Шаг 7: Тестирование обновлений

```bash
# Изменение тегов сервера
# Отредактируйте main.tf, добавьте новые теги

# Планирование обновления
terraform plan

# Применение изменений
terraform apply
```

### Шаг 8: Тестирование множественных серверов

```bash
# Включение второго сервера
echo 'create_second_server = true' >> terraform.tfvars

# Планирование и создание
terraform plan
terraform apply

# Проверка
terraform output test_server_2_info
```

## 🧪 Сценарии тестирования

### Базовое тестирование
```bash
export TF_VAR_test_scenario="basic"
terraform apply
```

### Тестирование множественных ресурсов
```bash
export TF_VAR_test_scenario="multiple"
export TF_VAR_create_second_server="true"
terraform apply
```

### Расширенное тестирование
```bash
export TF_VAR_test_scenario="advanced"
# Дополнительные тесты...
```

## 🔍 Проверка функциональности

### Проверка data sources
- ✅ Получение списка локаций
- ✅ Фильтрация по имени локации
- ✅ Получение услуг с фильтрацией
- ✅ Получение OS шаблонов
- ✅ Получение тарифных планов

### Проверка ресурсов
- ✅ Создание сервера
- ✅ Обновление тегов и имени
- ✅ Получение статуса сервера
- ✅ Получение IP адресов
- ✅ Импорт существующих серверов

### Проверка жизненного цикла
```bash
# Импорт существующего сервера
terraform import selectel_baremetal_server.imported_server "server-uuid"

# Обновление ресурса
terraform apply

# Удаление ресурса (ОСТОРОЖНО!)
terraform destroy -target=selectel_baremetal_server.test_server
```

## 🧹 Очистка ресурсов

⚠️ **ВАЖНО:** Не забудьте удалить созданные серверы!

```bash
# Удаление всех ресурсов
terraform destroy

# Или удаление конкретного сервера
terraform destroy -target=selectel_baremetal_server.test_server
```

## 📊 Мониторинг результатов

```bash
# Сводка тестирования
terraform output test_summary

# Детальная информация
terraform show

# Состояние ресурсов
terraform state list
terraform state show selectel_baremetal_server.test_server
```

## 🐛 Отладка

### Включение debug логов
```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log
terraform apply
```

### Проверка API запросов
```bash
# Логи HTTP запросов
export TF_LOG=TRACE
terraform apply 2>&1 | grep -E "(HTTP|API|Request|Response)"
```

### Проверка состояния провайдера
```bash
terraform providers
terraform version
```

## 📝 Отчет о тестировании

После завершения тестирования создайте отчет:

```bash
# Сохранение результатов
terraform output -json > test-results.json
terraform show -json > terraform-state.json

# Создание отчета
echo "# Отчет о тестировании $(date)" > test-report.md
echo "## Результаты:" >> test-report.md
terraform output test_summary >> test-report.md
```

## 🔒 Безопасность

- ✅ Credentials только в переменных окружения
- ✅ terraform.tfvars в .gitignore
- ✅ Sensitive переменные помечены как sensitive
- ✅ Логи не содержат приватных данных
- ✅ Lifecycle prevent_destroy для защиты

## 📞 Поддержка

При возникновении проблем:

1. Проверьте логи: `cat terraform.log`
2. Проверьте credentials: `echo $SELECTEL_TOKEN`
3. Проверьте API доступность: `curl -H "Authorization: Bearer $SELECTEL_TOKEN" $SELECTEL_ENDPOINT/locations`
4. Создайте issue в репозитории с логами (без credentials!)

---

**Удачного тестирования! 🚀** 
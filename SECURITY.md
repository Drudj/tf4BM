# Безопасность

## Важно! Никогда не коммитьте приватные данные

Этот проект настроен для безопасной работы с чувствительными данными:

### Файлы, исключенные из Git:
- `user.md` - файл с вашими credentials
- `*.tfvars` - файлы с переменными Terraform
- `test_*.tf` - тестовые файлы (кроме examples/)
- Все файлы состояния Terraform

### Как безопасно использовать провайдер:

1. **Через переменные окружения (рекомендуется):**
```bash
export SELECTEL_TOKEN="your-iam-token"
export SELECTEL_PROJECT_ID="your-project-uuid"
```

2. **Через файл terraform.tfvars (локально):**
```hcl
# terraform.tfvars (этот файл в .gitignore)
selectel_token = "your-iam-token"
selectel_project_id = "your-project-uuid"
```

3. **В конфигурации Terraform:**
```hcl
provider "selectel" {
  token      = var.selectel_token
  project_id = var.selectel_project_id
}
```

### Получение IAM токена:

```bash
curl -i -XPOST \
  -H 'Content-Type: application/json' \
  -d '{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"<service_user>","domain":{"name":"<account_id>"},"password":"<password>"}}},"scope":{"project":{"name":"<project_name>","domain":{"name":"<account_id>"}}}}}' \
  'https://cloud.api.selcloud.ru/identity/v3/auth/tokens'
```

Токен будет в заголовке `X-Subject-Token` ответа.

### Проверка безопасности перед коммитом:

```bash
# Проверить что приватные файлы игнорируются
git check-ignore user.md
git status --ignored

# Убедиться что нет токенов в коде
grep -r "gAAAAA" . --exclude-dir=.git
grep -r "token.*=" . --exclude-dir=.git --exclude="*.md"
``` 
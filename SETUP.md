# Безопасная настройка Terraform Provider для Selectel

## Подготовка переменных окружения

### 1. Создайте файл с переменными окружения

```bash
cp test.env.example test.env
```

### 2. Заполните файл test.env своими данными

```bash
# Редактируйте файл test.env
nano test.env
```

Замените следующие значения на ваши реальные:
- `your-region` → ваш регион (например: ru-9)
- `your-domain-id` → ID домена Selectel
- `your-username` → имя пользователя
- `your-password` → пароль пользователя
- `your-servers-api-token` → токен для API выделенных серверов

### 3. Загрузите переменные в терминал

```bash
source test.env
```

### 4. Проверьте настройки

```bash
echo "Domain: $OS_DOMAIN_NAME"
echo "Username: $OS_USERNAME"
echo "Region: $OS_REGION_NAME"
# НЕ выводите пароли и токены в логи!
```

## Сборка провайдера

```bash
cd terraform-provider-selectel
go build -o ../terraform-provider-selectel_v1.0.0
cd ..
```

## Настройка Development Overrides

Создайте файл `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "selectel/selectel" = "/полный/путь/к/проекту"
  }
  direct {}
}
```

## Безопасность

⚠️ **ВАЖНО**: 
- Никогда не коммитьте файл `test.env` в git
- Файл `test.env` уже добавлен в `.gitignore`
- Используйте переменные окружения для всех секретных данных
- Проверяйте `.gitignore` перед каждым коммитом

## Проверка безопасности

Перед коммитом выполните:

```bash
# Проверка на утечку секретов
grep -r "gAAAAA\|OS_PASSWORD\|SEL_SERVERS_TOKEN" . --exclude-dir=.git --exclude="*.md" --exclude="*.example"

# Убедитесь, что test.env игнорируется
git status --ignored
```

Если команда `grep` что-то находит - НЕ коммитьте изменения! 
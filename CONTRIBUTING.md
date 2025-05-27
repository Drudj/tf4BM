# Участие в разработке

Мы приветствуем участие в разработке terraform провайдера для управления выделенными серверами Selectel!

## Требования для разработки

- [Go](https://golang.org/doc/install) >= 1.23
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [golangci-lint](https://golangci-lint.run/usage/install/) для линтинга
- [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs) для генерации документации

## Настройка окружения разработки

1. Клонируйте репозиторий:
```bash
git clone https://github.com/selectel/terraform-provider-selectel-baremetal.git
cd terraform-provider-selectel-baremetal
```

2. Установите зависимости:
```bash
make deps
```

3. Установите инструменты разработки:
```bash
make dev-setup
```

## Процесс разработки

### Сборка

```bash
make build
```

### Тестирование

```bash
# Unit тесты
make test

# Тесты с покрытием
make test-coverage

# Acceptance тесты (требуют API токены)
make testacc
```

### Форматирование и линтинг

```bash
# Форматирование кода
make fmt

# Линтинг
make lint

# Полная проверка
make check
```

### Локальная установка

```bash
make install
```

## Структура проекта

```
.
├── cmd/terraform-provider-selectel-baremetal/  # Точка входа
├── internal/
│   ├── provider/                               # Основной провайдер
│   ├── resources/                              # Terraform ресурсы
│   ├── datasources/                            # Terraform data sources
│   ├── client/                                 # HTTP клиент для API
│   ├── models/                                 # Модели данных
│   └── utils/                                  # Утилиты
├── examples/                                   # Примеры использования
├── docs/                                       # Документация
├── test/                                       # Тесты
└── scripts/                                    # Скрипты автоматизации
```

## Стандарты кодирования

### Go код

- Следуйте стандартам Go (используйте `go fmt`)
- Используйте `golangci-lint` для проверки качества кода
- Пишите тесты для всех новых функций
- Документируйте публичные функции и типы

### Terraform код

- Используйте `terraform fmt` для форматирования
- Следуйте [best practices](https://developer.hashicorp.com/terraform/plugin/best-practices) HashiCorp

## Добавление новых ресурсов

1. Создайте модель данных в `internal/models/`
2. Добавьте методы API в `internal/client/`
3. Создайте ресурс в `internal/resources/` или data source в `internal/datasources/`
4. Зарегистрируйте в провайдере (`internal/provider/provider.go`)
5. Добавьте тесты
6. Создайте документацию в `docs/`
7. Добавьте пример в `examples/`

## Тестирование

### Unit тесты

Пишите unit тесты для всех новых функций:

```go
func TestNewFunction(t *testing.T) {
    // Ваш тест
}
```

### Acceptance тесты

Для acceptance тестов используйте реальное API:

```go
func TestAccResourceServer_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Ваши тестовые шаги
        },
    })
}
```

## Документация

- Документируйте все ресурсы и data sources в `docs/`
- Используйте примеры в документации
- Обновляйте README.md при необходимости

## Процесс отправки изменений

1. Создайте fork репозитория
2. Создайте feature branch: `git checkout -b feature/my-new-feature`
3. Внесите изменения и добавьте тесты
4. Убедитесь, что все тесты проходят: `make check`
5. Зафиксируйте изменения: `git commit -am 'Add some feature'`
6. Отправьте в ваш fork: `git push origin feature/my-new-feature`
7. Создайте Pull Request

## Стиль commit сообщений

Используйте [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add server resource
fix: handle API timeout errors
docs: update installation guide
test: add unit tests for client
```

## Отчеты об ошибках

При создании issue включите:

- Версию провайдера
- Версию Terraform
- Конфигурацию Terraform (без секретов)
- Полный вывод ошибки
- Шаги для воспроизведения

## Вопросы и поддержка

- [GitHub Issues](https://github.com/selectel/terraform-provider-selectel-baremetal/issues) - для багов и feature requests
- [GitHub Discussions](https://github.com/selectel/terraform-provider-selectel-baremetal/discussions) - для вопросов и обсуждений

## Лицензия

Участвуя в разработке, вы соглашаетесь с тем, что ваши изменения будут лицензированы под [Mozilla Public License 2.0](LICENSE). 
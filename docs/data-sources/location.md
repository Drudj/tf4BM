# selectel_baremetal_location

Получает информацию о конкретной локации (дата-центре) Selectel по имени.

## Пример использования

```hcl
data "selectel_baremetal_location" "moscow" {
  name = "Moscow"
}

output "moscow_uuid" {
  value = data.selectel_baremetal_location.moscow.uuid
}
```

## Схема аргументов

- `name` (String, Required) - Название локации для поиска.

## Схема атрибутов

- `id` (String) - Идентификатор data source (равен UUID локации).
- `uuid` (String) - UUID локации.
- `name` (String) - Название локации.
- `code` (String) - Код локации.
- `country` (String) - Страна.
- `city` (String) - Город.
- `description` (String) - Описание локации.
- `available` (Boolean) - Доступность локации.

## Пример вывода

```hcl
id          = "12345678-1234-1234-1234-123456789012"
uuid        = "12345678-1234-1234-1234-123456789012"
name        = "Moscow"
code        = "msk"
country     = "Russia"
city        = "Moscow"
description = "Moscow datacenter"
available   = true
```

## Обработка ошибок

Если локация с указанным именем не найдена, data source вернет ошибку:

```
Error: Client Error
Unable to find location 'NonExistentLocation', got error: location with name 'NonExistentLocation' not found
``` 
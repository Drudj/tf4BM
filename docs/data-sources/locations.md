# selectel_baremetal_locations

Получает список всех доступных локаций (дата-центров) Selectel для выделенных серверов.

## Пример использования

```hcl
data "selectel_baremetal_locations" "all" {}

output "available_locations" {
  value = data.selectel_baremetal_locations.all.locations
}
```

## Схема аргументов

Этот data source не принимает аргументов.

## Схема атрибутов

- `id` (String) - Идентификатор data source.
- `locations` (List of Object) - Список локаций. Каждая локация содержит:
  - `uuid` (String) - UUID локации.
  - `name` (String) - Название локации.
  - `code` (String) - Код локации.
  - `country` (String) - Страна.
  - `city` (String) - Город.
  - `description` (String) - Описание локации.
  - `available` (Boolean) - Доступность локации.

## Пример вывода

```hcl
locations = [
  {
    uuid        = "12345678-1234-1234-1234-123456789012"
    name        = "Moscow"
    code        = "msk"
    country     = "Russia"
    city        = "Moscow"
    description = "Moscow datacenter"
    available   = true
  },
  {
    uuid        = "87654321-4321-4321-4321-210987654321"
    name        = "Saint Petersburg"
    code        = "spb"
    country     = "Russia"
    city        = "Saint Petersburg"
    description = "Saint Petersburg datacenter"
    available   = true
  }
]
``` 
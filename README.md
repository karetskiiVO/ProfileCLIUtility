# ProfileCLIUtility

Небольшая CLI-утилита для управления "профилями", которые хранятся как отдельные YAML-файлы в выбранной директории.

## Требования

- Go 1.25+

## Сборка

Windows:

```bash
go build -o profile-utility.exe .
```

Linux/macOS:

```bash
go build -o profile-utility .
```

## Быстрый старт

Контейнер профилей — это папка, где лежат файлы `*.yaml`.
По умолчанию используется текущая директория, либо укажите её через `--path`.

Создать профиль:

```bash
profile-utility profile create --name dev --user alice --project proj [--path <DIR>]
profile-utility profile create --name=dev --user=alice --project=proj [--path=<DIR>]
```

Получить профиль:

```bash
profile-utility profile get --name dev [--path <DIR>]
profile-utility profile get --name=dev [--path=<DIR>]
```

Показать список:

```bash
profile-utility profile list [--path <DIR>]
profile-utility profile list [--path=<DIR>]
```

Строгий режим (ошибка при невалидных профилях):

```bash
profile-utility profile list [--path <DIR> | --strict]
profile-utility profile list [--path=<DIR> | --strict]
```

Удалить профиль:

```bash
profile-utility profile delete --name dev [--path <DIR>  | -v]
profile-utility profile delete --name=dev [--path=<DIR>  | -v]

```

## Формат профиля

Профиль хранится в файле `<name>.yaml` и содержит поля:

```yaml
user: alice
project: proj
```

Парсинг YAML строгий: неизвестные поля считаются ошибкой (в командах `get` и `delete` всегда, в `list` — при `--strict`).

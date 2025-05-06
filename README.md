# **biathlon_competitions**
[![Tests](https://github.com/PushinMax/biathlon_competitions/actions/workflows/go.yml/badge.svg)](https://github.com/PushinMax/biathlon_competitions/actions/workflows/go.yml)


Прототип системы для биатлонных соревнований.

## *Содержание*
- _Запуск_
- _Принцип работы_


## Запуск
Программа дает возможность выбора как файла конфигурации, так и файла с командами. По умолчанию в качестве таких файлов выступают: [configs/config.json](configs/config.json), [configs/events](configs/events)

Запуск с файлами по умолчанию:
```bash
make run
```

Запуск с кастомными путями:
```bash
make run-custom EVENT=path/to/events CONFIG=path/to/config.json
```


## Покрытие тестами

Тестами была покрыт функционал парсинга events из файла, а также реал



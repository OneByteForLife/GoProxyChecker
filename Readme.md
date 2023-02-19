# GoProxyChecker

## Дорожная карта

- [GoProxyChecker](#goproxychecker)
  - [Дорожная карта](#дорожная-карта)
  - [Описание](#описание)
  - [Структура проекта](#структура-проекта)
  - [Разработчики](#разработчики)
  - [Лицензия](#лицензия)

## Описание

Проект служит для проверки прокси адресов для дальнейшего их использования.

Он это часть другого проекта связанного с прокси адресами.

## Структура проекта
``` 
.
├── app
│   └── cmd
│       └── main.go
├── Dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── run.go
│   ├── models
│   │   └── checker.go
│   └── proxy
│       └── proxy.go
├── pkg
│   ├── config
│   │   ├── config.go
│   │   └── config.yaml
│   ├── database
│   │   └── connect.go
│   ├── http_check
│   │   └── http_check.go
│   └── log
│       └── log.go
└── Readme.md
```
## Разработчики

- [OneByteForLife](https://github.com/OneByteForLife)
  
## Лицензия

- Это программное обеспечение защищено лицензией MIT!
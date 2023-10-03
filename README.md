# Driver Installer

Этот проект представляет собой инструмент для автоматической установки драйверов из текущей директории.

## Зависимости

- Go
- [progressbar](https://github.com/schollz/progressbar)

## Установка

Чтобы установить проект, выполните следующие команды:

```bash
go get github.com/schollz/progressbar/v3
```

## Использование

Для запуска установщика драйверов выполните следующую команду:

```bash
go run main.go
```

Чтобы включить логирование в файл `log.txt`, выполните следующую команду:

```bash
go run main.go log
```

## Описание

Программа сканирует директории и ищет файлы с расширением .inf, которые представляют собой драйверы устройств. Для каждого найденного файла .inf программа запускает pnputil.exe с аргументами для установки драйвера. Прогресс установки отображается с использованием текстового индикатора прогресса.

Если включено логирование, все сообщения и ошибки будут записаны в файл log.txt.

По окончании установки программа выводит количество успешно установленных и неудачно установленных драйверов.

## Лицензия

Этот проект лицензирован под лицензией Apache 2.0

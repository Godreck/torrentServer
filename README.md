# ВКР-2025-ВеремьёвАрсенийНиколаевич-TorrentServer

Выпускная квалификационная работа "Разработка серверной части сервиса трансляции видео с помощью протокола распределённой передачи данных". Автор: Веремьёв Арсений Николаевич, группа 4013.

## Run
Для запуска клонируйте репозиторий, затем используйте один из способов ниже.

### Первый способ
```
export UID=$(id -u)
export GID=$(id -g)
```

затем

```
docker-compose up -d
```

### Второй способ
```
UID=$(id -u) GID=$(id -g) docker-compose up -d
```

### Остановка сервиса

Чтобы остановить работу сервиса выполните команду:

```
docker compose down -v
```

### Флаги в docker compose
Используйте флаги по необходимости, подробнее о флагах можно узнать по ссылке: https://www.geeksforgeeks.org/flag-do-in-docker-compose/

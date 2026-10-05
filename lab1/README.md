# Практическая работа №1. Docker image, volume, container

## Цель
Получить навык работы с Docker image, Docker volume и Docker container, а также научиться подключать контейнеры к сетям и организовывать обмен данными через volumes.

## 1. Подготовка

```bash
docker pull nginx:alpine
docker pull alpine:latest
docker network ls
docker volume ls
```

## 2. Запуск двух контейнеров

```bash
docker run -d --name web1 -p 8080:80 nginx:alpine
docker run -d --name tools1 alpine:latest sleep 1d
docker ps
```

Порт 8080 хоста проброшен на порт 80 контейнера web1.

## 3. Сеть и взаимодействие контейнеров

```bash
docker network create lab-net
docker network connect lab-net web1
docker network connect lab-net tools1
docker exec tools1 ping -c 4 web1
docker exec tools1 wget -qO- http://web1
```

Контейнер `tools1` находит `web1` по имени через встроенный DNS Docker.

## 4. Изолированная сеть (без внешнего доступа)

```bash
docker network create --internal isolated-net
docker network connect isolated-net tools1
docker exec tools1 sh -c "apk update"
```

Ожидаемый результат: ошибка обращения к внешним репозиториям (нет выхода в интернет).

## 5. tmpfs

```bash
docker run -d --name tmpfs-demo --tmpfs /app/tmp:rw,size=64m alpine:latest sleep 1d
docker inspect tmpfs-demo
docker exec tmpfs-demo sh -c "mount | grep /app/tmp"
```

tmpfs хранится в оперативной памяти и не создаёт постоянный Docker volume.

## 6. Общий volume для двух контейнеров

```bash
docker volume create shared-data

docker run -d --name writer -v shared-data:/shared alpine:latest sleep 1d
docker run -d --name reader -v shared-data:/shared alpine:latest sleep 1d

docker exec writer sh -c 'echo "Docker shared volume test" > /shared/info.txt'
docker exec writer sh -c 'date > /shared/date.txt'

docker exec reader ls -la /shared
docker exec reader cat /shared/info.txt

docker volume inspect shared-data
```

Файлы, созданные в `writer`, видны в `reader` — оба используют один volume.

## 7. Очистка

```bash
docker rm -f web1 tools1 tmpfs-demo writer reader
docker network rm lab-net isolated-net
docker volume rm shared-data
```

## Вывод
Освоены базовые операции с образами, контейнерами, сетями и volumes. Подтверждено, что пользовательская сеть даёт DNS-разрешение имён контейнеров, isolated-сеть блокирует внешний доступ, tmpfs хранит данные в RAM, а общий volume обеспечивает обмен файлами между контейнерами.
# Практическая работа №3. Docker Hub, Docker Registry и DCT

## Цель
Получить навык работы с локальным Docker Registry, Docker Hub и механизмом подписывания образов.

## 1. Развёртывание локального Registry

```bash
docker run -d --name registry -p 5000:5000 --restart=always registry:2
docker ps
```

## 2. Получение образа с Docker Hub

```bash
docker pull nginx:alpine
docker images
```

## 3. Тег и push в локальный Registry

```bash
docker tag nginx:alpine localhost:5000/my-nginx:1.0
docker push localhost:5000/my-nginx:1.0

curl http://localhost:5000/v2/_catalog
curl http://localhost:5000/v2/my-nginx/tags/list
```

## 4. Push в Docker Hub

```bash
docker login
docker tag nginx:alpine YOUR_DOCKERHUB_USERNAME/my-nginx:1.0
docker push YOUR_DOCKERHUB_USERNAME/my-nginx:1.0
```

`YOUR_DOCKERHUB_USERNAME` заменить на собственный логин Docker Hub. Репозиторий должен быть публичным.

## 5. DCT (опционально)

```bash
export DOCKER_CONTENT_TRUST=1
docker push YOUR_DOCKERHUB_USERNAME/my-nginx:1.0
```

Docker Content Trust требует поддержки Notary. Пункт опциональный.

## Вывод
Освоены локальный Registry (localhost:5000) и публикация образа в Docker Hub. Проверены catalog и tags через Registry HTTP API.
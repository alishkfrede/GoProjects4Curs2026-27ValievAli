# Практическая работа №2. Dockerfile и Docker Builder

## Цель
Научиться составлять Dockerfile и выполнять многоэтапную (multi-stage) сборку.

## 1. Создание React-приложения

```bash
node --version
npm --version
npx create-react-app my-app
cd my-app
```

## 2. Структура проекта

```
my-app/
├── public/
├── src/
├── Dockerfile
├── .dockerignore
├── myapp_nginx.conf
├── package.json
└── package-lock.json
```

Файлы `Dockerfile`, `.dockerignore`, `myapp_nginx.conf` лежат в этом же каталоге (lab2).

## 3. Сборка образа

```bash
docker build -t my-react-app:1.0 .
docker images
```

Ключ `-t` задаёт имя и тег. Точка в конце — build context (текущий каталог).

## 4. Запуск и проверка

```bash
docker run -d --name react-app -p 8080:80 my-react-app:1.0
docker ps
```

Открыть в браузере: http://localhost:8080 — должна отобразиться стартовая страница React.

## Пояснение по multi-stage
Первый этап (`builder`) ставит зависимости и собирает production-сборку. Второй этап (`nginx:alpine`) содержит только nginx и статические файлы — инструменты Node.js в итоговый образ не попадают, образ получается меньше.

## Вывод
Получен навык написания Dockerfile и многоэтапной сборки. Итоговый образ содержит только runtime-часть.
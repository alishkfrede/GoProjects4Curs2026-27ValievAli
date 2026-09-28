# Практическая работа №3 — Docker Hub, Registry и DCT

Локальный Registry выполняется полностью без внешнего аккаунта. Для Docker Hub используется интерактивный вход `docker login`, после чего скрипт запросит имя Docker Hub пользователя.

## Локальная часть
```bash
chmod +x run_registry.sh
./run_registry.sh
```

## Docker Hub
```bash
chmod +x push_dockerhub.sh
./push_dockerhub.sh
```

## DCT
DCT указан в методичке как дополнительный пункт. Его поддержка зависит от версии Docker/учебной инфраструктуры. Скрипт запускается отдельно.

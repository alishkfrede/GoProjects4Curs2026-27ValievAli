# Практическая работа №5. Docker daemon и Storage Driver

Windows + Docker Desktop, PowerShell. Часть команд из методички (Linux) не работает — показал где.

## 1. Проверка

```bash
docker info | grep -i "Storage Driver"
```
→ ошибка: `grep` не распознан.

```bash
sudo systemctl status docker
```
→ ошибка: нет `sudo` и `systemctl`.

## 2. daemon.json

Путь в Windows: `C:\ProgramData\Docker\config\daemon.json`

```json
{
  "log-driver": "json-file",
  "log-opts": { "max-size": "10m", "max-file": "3" },
  "live-restore": true
}
```

## 3. Перезапуск

```bash
sudo systemctl restart docker
```
→ ошибка: нет `systemctl`.

```powershell
Restart-Service com.docker.service
```
→ ошибка: нужны права администратора.

Перезапустил Docker Desktop через трей.

Проверка:

```powershell
docker info
docker inspect logtest
```

Видно: `Logging Driver: json-file`, `Live Restore Enabled: true`, `max-size: 10m`, `max-file: 3`.

## 4. Флаги dockerd

```bash
dockerd --help | grep -E "log-driver|log-opt|live-restore"
```
→ ошибка: нет `dockerd` и `grep`.

## 5. Storage Driver

```bash
docker info | grep -i "Storage Driver"
sudo ls -la /var/lib/docker/overlay2
```
→ ошибки: нет `grep`, нет пути `/var/lib/docker` в Windows.

```powershell
docker info
```
→ `Storage Driver: overlay2`, `Docker Root Dir: /var/lib/docker`.

Каталог — внутри WSL2:

```powershell
wsl
```
```bash
sudo ls -la /var/lib/docker/overlay2
```

## 6. Смена драйвера

Добавил в `daemon.json` `"storage-driver": "vfs"`, перезапустил — драйвер остался `overlay2`.

Причина: в Docker Desktop на Windows драйвер фиксированный, смена невозможна.

## 7. Приложил

- `daemon.json`

## Вывод

Настроил `daemon.json` на Windows, проверил параметры. Команды Linux (`grep`, `sudo`, `systemctl`) не работают. Смена Storage Driver в Docker Desktop невозможна.
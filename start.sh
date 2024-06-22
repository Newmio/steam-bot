#!/bin/sh

# Изменяем MAC-адрес на случайный, если сетевой интерфейс существует
interface=$(ip -o link show | awk -F': ' '{print $2}' | grep -v "lo" | head -n 1)
if [ -n "$interface" ]; then
    macchanger -r $interface
    echo "Changed MAC address of $interface"
fi

# Изменяем UUID файловой системы на случайный, если файловая система поддерживает изменение UUID
fs=$(df / | tail -1 | awk '{print $1}')
if tune2fs -l $fs >/dev/null 2>&1; then
    tune2fs $fs -U random
    echo "Changed UUID of filesystem $fs"
fi

# Запускаем chromedriver в фоне
nohup chromedriver --port=9515 --verbose > /var/log/chromedriver.log 2>&1 &
chromedriver_pid=$!
echo "Started chromedriver with PID $chromedriver_pid"

# Ждем, пока chromedriver не будет готов
for i in $(seq 1 10); do
    if nc -z localhost 9515; then
        echo "Chromedriver is ready"
        break
    fi
    echo "Waiting for chromedriver to be ready..."
    sleep 1
done

if ! nc -z localhost 9515; then
    echo "Chromedriver failed to start"
    kill $chromedriver_pid
    exit 1
fi

# Запускаем Go-приложение
echo "Starting Go application"
if go run cmd/main.go; then
    echo "Go application started successfully"
else
    echo "Go application failed to start"
    kill $chromedriver_pid
    exit 1
fi
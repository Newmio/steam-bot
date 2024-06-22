#!/bin/sh

# Изменяем MAC-адрес на случайный
interface=$(ip -o link show | awk -F': ' '{print $2}' | grep -v "lo" | head -n 1)
macchanger -r $interface

# Изменяем UUID файловой системы на случайный
fs=$(df / | tail -1 | awk '{print $1}')
tune2fs $fs -U random

# Генерируем случайный серийный номер и сохраняем его в файл
serial=$(tr -cd 'a-f0-9' < /dev/urandom | head -c 32)
echo "Serial Number: $serial" > /etc/serial_number

# Запускаем chromedriver в фоне
nohup chromedriver --port=9515 --verbose &

# Ждем 5 секунд, чтобы убедиться, что chromedriver запущен
sleep 5

# Запускаем Go-приложение
go run cmd/main.go

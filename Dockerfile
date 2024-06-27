# Use an official Golang runtime as a parent image
FROM golang:1.22.0

# Install dependencies
RUN apt-get update && apt-get install -y \
    wget \
    unzip \
    gnupg \
    && rm -rf /var/lib/apt/lists/*

# Install Google Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list \
    && apt-get update \
    && apt-get install -y google-chrome-stable

# Install ChromeDriver
RUN wget -O /tmp/chromedriver.zip https://storage.googleapis.com/chrome-for-testing-public/126.0.6478.61/linux64/chromedriver-linux64.zip \
    && unzip /tmp/chromedriver.zip -d /usr/local/bin/ \
    && rm /tmp/chromedriver.zip \
    && mv /usr/local/bin/chromedriver-linux64/chromedriver /usr/local/bin/chromedriver \
    && chmod +x /usr/local/bin/chromedriver

# Set environment variables for Go
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# Create working directory
WORKDIR /go/src/app

# Copy the source code into the container
COPY . .

# Install Go modules
RUN go mod tidy

# Command to run the application
CMD sh -c "chromedriver --port=9515 & sleep 10 && go run cmd/main.go"




# # Используем базовый образ Ubuntu
# FROM ubuntu:latest

# # Обновляем пакеты и устанавливаем необходимые утилиты
# RUN apt-get update && apt-get install -y \
#     wget \
#     unzip \
#     curl \
#     gnupg \
#     software-properties-common 

# # Устанавливаем Google Chrome
# RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
#     && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list \
#     && apt-get update \
#     && apt-get install -y google-chrome-stable

# # Устанавливаем ChromeDriver нужной версии
# RUN wget -O /tmp/chromedriver.zip https://storage.googleapis.com/chrome-for-testing-public/126.0.6478.61/linux64/chromedriver-linux64.zip \
#     && unzip /tmp/chromedriver.zip -d /usr/local/bin/ \
#     && rm /tmp/chromedriver.zip \
#     && mv /usr/local/bin/chromedriver-linux64/chromedriver /usr/local/bin/chromedriver \
#     && chmod +x /usr/local/bin/chromedriver

# # Устанавливаем Go нужной версии
# RUN wget -q https://go.dev/dl/go1.22.0.linux-amd64.tar.gz \
#     && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz \
#     && rm go1.22.0.linux-amd64.tar.gz

# # Устанавливаем переменные окружения для Go
# ENV PATH=$PATH:/usr/local/go/bin
# ENV GOPATH=/go

# # Создаем рабочую директорию
# WORKDIR /go/src/app

# # Копируем исходный код приложения в контейнер
# COPY . .

# # Устанавливаем Selenium и другие пакеты Go
# RUN go mod tidy

# # Команда для запуска сценария
# CMD sh -c "chromedriver --port=9515 & sleep 10 && go run cmd/main.go"

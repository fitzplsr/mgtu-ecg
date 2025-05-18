# Используем официальный образ Python
FROM python:3.11-slim

# Обновляем pip и ставим системные зависимости
RUN apt-get update && apt-get install -y \
    gcc \
    libglib2.0-0 \
    libffi-dev \
    build-essential \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*

# Установка Python-зависимостей
RUN pip install --upgrade pip
RUN pip install fastapi uvicorn pyedflib

# Создаем рабочую директорию
WORKDIR /app

# Копируем исходный код
COPY main.py vecg.py converter.py ./

# Указываем команду запуска
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8081"]

# Используем официальный минимальный образ Python
FROM python:3.11-slim

# Обновляем pip и ставим системные и библиотечные зависимости
RUN apt-get update && apt-get install -y \
    gcc \
    libglib2.0-0 \
    libffi-dev \
    libstdc++6 \
    libpulse0 \
    libxrender1 \
    libsm6 \
    libxext6 \
    libasound2 \
    libz-dev \
    libsndfile1 \
    build-essential \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем зависимости Python
RUN pip install --upgrade pip

# Установка всех нужных библиотек (модель, API, обработка)
RUN pip install \
    fastapi \
    uvicorn \
    pyedflib \
    numpy \
    matplotlib \
    scikit-learn \
    imbalanced-learn \
    tqdm \
    torch \
    torchvision \
    torchaudio

# Создаем рабочую директорию
WORKDIR /app

# Копируем весь код проекта в контейнер
COPY . .

# Указываем команду запуска FastAPI-приложения
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8082"]

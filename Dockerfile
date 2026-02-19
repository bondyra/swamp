FROM node:22-slim AS frontend

ARG version
ENV REACT_APP_VERSION=${version}

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci --no-audit --no-fund

COPY frontend/ ./
RUN npm run build


FROM python:3.12-slim
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    POETRY_HOME=/usr/local \
    POETRY_VIRTUALENVS_CREATE=false
EXPOSE 80 8000
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        nginx \
        supervisor \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /backend
COPY backend/requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
COPY backend/ ./
COPY --from=frontend /app/build /var/www/html
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
WORKDIR /

CMD ["supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

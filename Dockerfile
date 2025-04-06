FROM node:22 AS frontend
ARG version
ENV REACT_APP_VERSION=${version}
COPY frontend /app
WORKDIR /app
RUN npm install && npm run build

FROM python:3.12.3-slim AS server
EXPOSE 80
EXPOSE 8000
ENV POETRY_HOME='/usr/local' POETRY_VIRTUALENVS_CREATE='false'
RUN apt-get update && apt-get install -y supervisor nginx && apt-get clean
COPY --from=frontend /app/build /var/www/html
COPY backend /backend
WORKDIR /backend
RUN pip install -r requirements.txt
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
WORKDIR /
CMD ["supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

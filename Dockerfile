FROM node:22 AS frontend
ARG version
ENV REACT_APP_VERSION=${version}
COPY frontend /app
WORKDIR /app
RUN npm run build

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
COPY nginx.conf /etc/nginx/nginx.conf
WORKDIR /
CMD ["supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]


# docker run -p 3000:80 -p 8000:8000 -v $HOME/.aws:/root/.aws -v $HOME/.kube:/root/.kube -v $HOME/.minikube:/root/.minikube swamp 
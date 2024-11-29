FROM node:22

ENV POETRY_HOME='/usr/local' POETRY_VIRTUALENVS_CREATE='false'

EXPOSE 3000
EXPOSE 8000

RUN apt-get update && apt-get install -y curl supervisor python3 python3-pip python3.11-venv && apt-get clean

RUN curl -sSL https://install.python-poetry.org | python3 -

RUN python3 -m venv /usr/local/venv && . /usr/local/venv/bin/activate

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ui /ui
COPY swamp /swamp

# CMD ["/bin/bash", "-c", "sleep 123213"]
WORKDIR /ui/backend
RUN poetry install

WORKDIR /swamp
RUN poetry install

WORKDIR /ui/frontend
RUN npm install

WORKDIR /
CMD ["supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

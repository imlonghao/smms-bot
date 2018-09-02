FROM python:3.7-alpine
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
WORKDIR /app
COPY bot.py ./
COPY requirements.txt ./
RUN apk add --no-cache --virtual build-dependencies gcc libffi-dev openssl-dev && \
    pip install --no-cache-dir -r requirements.txt && \
    apk del build-dependencies
ENTRYPOINT [ "/app/bot.py" ]
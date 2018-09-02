FROM python:3.7-alpine
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
WORKDIR /app
COPY bot.py ./
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
ENTRYPOINT [ "/app/bot.py" ]
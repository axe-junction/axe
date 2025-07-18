FROM python:3.11-slim

WORKDIR /app
RUN apt-get update && apt-get install -y make
RUN apt-get update && apt-get install -y make

COPY . /app

RUN pip install --no-cache-dir -r requirements.txt
EXPOSE 8000
RUN prisma generate
CMD ["make", "run"]


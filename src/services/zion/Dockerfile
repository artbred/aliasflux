FROM python:3.11-slim-buster

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y build-essential

RUN apt-get install -y libgeos-dev binutils libproj-dev

COPY . /app
WORKDIR /app

RUN pip install --no-cache-dir -r requirements.txt

CMD ["python3", "src/services/zion/main.py"]

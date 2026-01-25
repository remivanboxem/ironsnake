run:
    docker compose up -d

run-dev:
    docker compose -f docker-compose.dev.yaml up -d --build

stop:
    docker compose down

clean:
    docker compose down --remove-orphans

build:
    docker compose build

build-dev:
    docker compose -f docker-compose.dev.yaml build

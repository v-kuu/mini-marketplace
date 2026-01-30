build:
	go build ./cmd/server

run:
	./server

up:
	docker compose up --build -d

down:
	docker compose down

clean:
	docker compose down -v --remove-orphans
	docker builder prune -af
	docker image prune -af
	docker volume prune -af

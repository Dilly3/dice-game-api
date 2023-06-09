

setup-docker:
	@docker run --name dice-game -p 4300:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=dice_game -d postgres:15-alpine
	
setup-air:
	@go install github.com/cosmtrek/air@latest
	@air init

docker-start:
	@docker start dice-game 

migrate-up:
	@migrate -path db/migrations -database postgresql://root:root@localhost:4300/dice_game?sslmode=disable -verbose up

#  install air:  go install github.com/cosmtrek/air@latest
# run air init to create air.toml file
# run air to start the server

migrate-down:
	@migrate -path db/migrations -database postgresql://root:root@localhost:4300/dice_game?sslmode=disable -verbose down
migrate-force:
	@migrate -path db/migrations -database postgresql://root:root@localhost:4300/dice_game?sslmode=disable force 1
sqlc:
	@sqlc generate

migrate-up-test:
	@migrate -path db/migrations -database postgresql://root:root@localhost:4300/test-db?sslmode=disable -verbose up

migrate-down-test:
	@migrate -path db/migrations -database postgresql://root:root@localhost:4300/test-db?sslmode=disable -verbose down
test:
	@go test -v -cover ./...
run:
	@go run main.go 
air:
	@air
mock:
	@mockgen -package mockdb -destination db/mock/store.go github.com/dilly3/dice-game-api/db/sqlc Store

# docker commands to stop and remove the container
docker-stop:
	@docker stop dice-game
docker-remove:
	@docker rm dice-game


.PHONY: migrate-up migrate-down sqlc  test run docker-start docker-stop migrate-force migrate-up-test migrate-down-test docker-remove air setup-air setup-docker
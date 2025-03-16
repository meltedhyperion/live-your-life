server-build:
	@cd server && go build -o bin/globetrotter-server ./cmd/*.go

server-start: server-build

	@cd server && ./bin/globetrotter-server

server-dev:
	@sh server/develop.sh

docker-compose-server:
	@docker-compose -f docker-compose-server.yaml up

server-run-tests:
	@cd server && go test -v ./...
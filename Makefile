server-build:
	@cd server && go build -o bin/globetrotter-server ./cmd/*.go

server-start: server-build

	@cd server && ./bin/globetrotter-server

server-dev:
	@sh server/develop.sh

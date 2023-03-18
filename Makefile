# Build the program
build:
	go build -o ./bin/metacritic-api ./cmd/metacriticapi/

# Clean up build artifacts
clean:
	rm -f ./bin/metacritic-api

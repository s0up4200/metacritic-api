# Build the program
build:
	go build -o ./bin/metacritic-api

# Clean up build artifacts
clean:
	rm -f ./bin/metacritic-api
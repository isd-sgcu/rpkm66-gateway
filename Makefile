create-doc:
	swag init -d ./src -o ./docs -md ./docs/markdown

test:
	go vet ./...
	go test  -v -coverpkg ./app/... -coverprofile coverage.out -covermode count ./app/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

server:
	go run ./.

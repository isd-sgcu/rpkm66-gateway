proto:
	protoc --proto_path=src/proto --go_out=. --go-grpc_out=. user.proto
	protoc --proto_path=src/proto --go_out=. --go-grpc_out=. auth.proto
	protoc --proto_path=src/proto --go_out=. --go-grpc_out=. baan.proto
	protoc --proto_path=src/proto --go_out=. --go-grpc_out=. file.proto
	protoc --proto_path=src/proto --go_out=. --go-grpc_out=. group.proto
	protoc --proto_path=src/proto --go_out=. --go-grpc_out=. checkin.proto
	protoc --proto_path=src/proto --go_out=. --go-grpc_out=. event.proto

create-doc:
	swag init -d ./src -o ./src/docs -md ./src/docs/markdown

test:
	go vet ./...
	go test  -v -coverpkg ./src/app/... -coverprofile coverage.out -covermode count ./src/app/...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

server:
	go run ./src/.

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

.PHONY: run down exec logs proto

run:
	docker-compose up -d --build

down:
	docker compose down

exec:
	docker compose exec -it $(name) sh

logs:
	docker compose logs -f $(name)

proto:
	# Clean up old stubs
	rm -f ./collector/gen/pb/*.go
	rm -f ./management/gen/pb/*.go
	rm -f ./gateway/gen/pb/*.go

	# Stubs for collector service
	protoc --proto_path=proto --go_out=./collector/gen/pb --go_opt=paths=source_relative \
		--go-grpc_out=./collector/gen/pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

	# Stubs for management service
	protoc --proto_path=proto --go_out=./management/gen/pb --go_opt=paths=source_relative \
		--go-grpc_out=./management/gen/pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

	# Stubs for gateway service
	protoc --proto_path=proto --go_out=./gateway/gen/pb --go_opt=paths=source_relative \
		--go-grpc_out=./gateway/gen/pb --go-grpc_opt=paths=source_relative \
	    --grpc-gateway_out=./gateway/gen/pb --grpc-gateway_opt=paths=source_relative \
        proto/*.proto

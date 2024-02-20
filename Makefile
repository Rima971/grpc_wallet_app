generate_grpc_code:
	protoc \
	--go_out=authenticator \
	--go_opt=paths=source_relative \
	--go-grpc_out=authenticator \
	--go-grpc_opt=paths=source_relative \
	authenticator.proto
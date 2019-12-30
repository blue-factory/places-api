#
# SO variables
#
# DOCKER_USER
# DOCKER_TOKEN
#

#
# Internal variables
#
VERSION=0.0.7
SVC=places-api
BIN_PATH=$(PWD)/bin
BIN=$(BIN_PATH)/$(SVC)
REGISTRY_URL=$(DOCKER_USER)

clean c:
	@echo "[clean] Cleaning bin folder..."
	@rm -rf bin/

run r:
	@echo "[running] Running service..."
	@PORT=5030 NATS_HOST=localhost NATS_PORT=4222 API_KEY=<api_key> go run cmd/main.go

build b: proto
	@echo "[build] Building service..."
	@cd cmd && go build -o $(BIN)

linux l:
	@echo "[build-linux] Building service..."
	@cd cmd && GOOS=linux GOARCH=amd64 go build -o $(BIN)

docker d:
	@echo "[docker] Building image..."
	@docker build -t $(SVC):$(VERSION) .

add-migration am: 
	@echo "[add-migration] Adding migration"
	@goose -dir "./database/migrations" create $(name) sql

migrations m:
	@echo "[migrations] Runing migrations..."
	@cd database/migrations && goose postgres $(DSN) up

docker-login dl:
	@echo "[docker] Login to docker..."
	@docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)

push p: linux docker docker-login
	@echo "[docker] pushing $(REGISTRY_URL)/$(SVC):$(VERSION)"
	@docker tag $(SVC):$(VERSION) $(REGISTRY_URL)/$(SVC):$(VERSION)
	@docker push $(REGISTRY_URL)/$(SVC):$(VERSION)

compose co:
	@echo "[docker-compose] Running docker-compose..."
	@docker-compose build
	@docker-compose up

stop s: 
	@echo "[docker-compose] Stopping docker-compose..."
	@docker-compose down

clean-proto cp:
	@echo "[clean-proto] Cleaning proto files..."
	@rm -rf proto/*.pb.go || true

proto pro: clean-proto
	@echo "[proto] Generating proto file..."
	@protoc --go_out=plugins=grpc:. ./proto/*.proto 

.PHONY: clean c run r build b linux l docker d add-migration am migrations m compose co stop s clean-proto cp proto pro
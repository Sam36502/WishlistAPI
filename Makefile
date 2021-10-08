EXE_LINUX = "api_linux"
EXE_WIN = "api_win.exe"
DOCKER_IMAGE = "wishlist_api"
CONTAINER_NAME ="wishlist_api"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Builds the executable for linux
build:
	@echo "---> Building Linux Executable"
	@GOOS="linux" go build -o ${EXE_LINUX} ./src/

## Builds the executable for windows
build-win:
	@echo "---> Building Windows Executable"
	@GOOS="windows" go build -o ${EXE_WIN} ./src/

## Builds the docker image
image: build
	@echo "---> Building Docker Image"
	@docker build -t ${DOCKER_IMAGE} .

## Starts the docker-compose cluster
up: image
	@echo "---> Starting Compose Cluster"
	@docker-compose up -d

## Stops the docker-compose cluster
down:
	@echo "---> Stopping Compose Cluster"
	@docker-compose down
	@docker-compose rm -f

## Connects to api container
bash-api:
	@docker exec -it wishlistapi_wishlist_api_1 sh

## Connects to db container
bash-db:
	@docker exec -it wishlistapi_wishlist_db_1 sh
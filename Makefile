DOCKER := docker
DOCKER_COMPOSE := ${DOCKER} compose

PROJECT_NAME := gemini
PROJECT_COMPOSE := PROJECT_NAME=${PROJECT_NAME} \
	${DOCKER_COMPOSE}\
		--project-directory .\
		-p ${PROJECT_NAME}\
		-f ./build/docker-compose.yaml

build-image:
	${DOCKER} build\
		-f ./build/go.Dockerfile\
		-t ${PROJECT_NAME}-go-mongo:latest\
		.

	-${DOCKER} images -q -f "dangling=true" | xargs ${DOCKER} rmi -f

teardown:
	${PROJECT_COMPOSE} down\
		--remove-orphans

	${PROJECT_COMPOSE} down -v

mongodb:
	${PROJECT_COMPOSE} up -d mongodb

cronjob:
	${PROJECT_COMPOSE} run importer go run\
		-C=/go/src\
		-mod=vendor\
		./cmd

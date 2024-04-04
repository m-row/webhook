include .env
export

CONTAINER_NAME:=${CONTAINER_NAME}
CONTAINER_TAG:=${APP_VER}.$(shell git rev-list --count HEAD).$(shell git describe --always)
CONTAINER_IMG:=${CONTAINER_NAME}:${CONTAINER_TAG}

# Docker Hub
.PHONY: dh dh/down dh/push prune ps inspect
dh:
	export CONTAINER_TAG=${CONTAINER_TAG}
	docker compose -f docker-compose.yml up --build -d
dh/down:
	docker compose -f docker-compose.yml down
dh/push: dh
	docker tag sadeem/${CONTAINER_IMG} ${CONTAINER_REG}/${CONTAINER_IMG}
	docker tag sadeem/${CONTAINER_IMG} ${CONTAINER_REG}/${CONTAINER_NAME}:latest
	docker push ${CONTAINER_REG}/${CONTAINER_NAME} -a
prune:
	docker system prune -a -f --volumes
ps:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.RunningFor}}\t{{.Size}}\t{{.Ports}}"
# inspect a container local ip n=name of container
inspect: 
	docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" $(n)

# Go
.PHONY: list update
list:
	go list -m -u
update:
	go get -u ./...

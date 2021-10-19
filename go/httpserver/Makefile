tag ?= latest
HUB ?= sarah4test
IMAGE ?= httpserver

build: test
	echo "building httpserver container ${HUB}/${IMAGE}:${tag}"
	docker build -t ${HUB}/${IMAGE}:${tag} .
push: build
	echo "pushing ${HUB}/${IMAGE}:${tag}"
	docker push "${HUB}/${IMAGE}:${tag}"
run: build
	docker run -p 8080:8080 ${HUB}/${IMAGE}:${tag}
test:
	go test ./...
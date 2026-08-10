.PHONY: build run docker-build docker-run docker-stop clean

IMAGE_NAME := go-latex-app
CONTAINER_NAME := latex_app_container
PORT := 8080

build:
	go build -o main .

run: build
	./main

docker-build:
	docker build -t $(IMAGE_NAME) .

docker-run: docker-build
	docker run --rm -p $(PORT):$(PORT) --name $(CONTAINER_NAME) -v /mnt/c/github/go-latex/fonts:/usr/local/share/fonts/custom:ro $(IMAGE_NAME)

docker-stop:
	docker stop $(CONTAINER_NAME)

clean:
	rm -f main
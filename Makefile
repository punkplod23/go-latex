.PHONY: build run docker-build docker-run docker-stop clean download-fonts

IMAGE_NAME := go-latex-app
CONTAINER_NAME := latex_app_container
PORT := 8080
FONTS_DIR := fonts

build:
	go build -o main .

run: build
	./main

download-fonts:
	mkdir -p $(FONTS_DIR)
	@echo "Downloading Roboto-Regular.ttf..."
	curl -sL -k -o $(FONTS_DIR)/Roboto-Regular.ttf "https://raw.githubusercontent.com/googlefonts/roboto/main/src/hinted/Roboto-Regular.ttf"
	@echo "Downloading Roboto-Bold.ttf..."
	curl -sL -k -o $(FONTS_DIR)/Roboto-Bold.ttf "https://raw.githubusercontent.com/googlefonts/roboto/main/src/hinted/Roboto-Bold.ttf"
	@echo "Fonts successfully downloaded to $(FONTS_DIR)"

docker-build: download-fonts
	docker build -t $(IMAGE_NAME) .

docker-run: docker-build
	docker run --rm -p $(PORT):$(PORT) --name $(CONTAINER_NAME) -v /mnt/c/github/go-latex/fonts:/usr/local/share/fonts/custom:ro $(IMAGE_NAME)

docker-stop:
	docker stop $(CONTAINER_NAME)

clean:
	rm -f main
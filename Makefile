# 镜像名称和版本
MODULE=stu-go
TAG := $(shell date +'%Y%m%d%H%M%S')  # 使用时间戳作为标签
IMAGE_REPO=121.36.61.64:8888

init:
	@go mod tidy
	@rm -rf .env
	@cp -a .env.dev .env

run:
	@go run .

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${MODULE}

image: build
	@docker build -t ${IMAGE_REPO}/stu/${MODULE}:${TAG} .

push: image
	@docker push ${IMAGE_REPO}/stu/${MODULE}:${TAG}
	@docker rmi ${IMAGE_REPO}/stu/${MODULE}:${TAG}
	@echo
	@echo ${IMAGE_REPO}/stu/${MODULE}:${TAG}

.PHONY: init image build

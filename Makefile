# Makefile
.PHONY: help image docs

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make docs - 更新swagger文档"
	@echo "  make image - 构建Docker镜像"

# 变量定义
TAG ?= v0.1.2
REPO ?= registry-cn-hz.7net.cc/septnet
PROJECT_NAME ?= aliyun-api
ifeq ($(TAG),master)
	TAG := latest
endif

DOCKER_IMAGE = $(REPO)/$(PROJECT_NAME):$(TAG)
DOCKER_FILE = Dockerfile

docs:
	swag init -g aliyun.go


# 构建多阶段Docker镜像（如果存在）
image:
	docker build --push -t $(DOCKER_IMAGE) .
	echo "构建成功，镜像:" $(DOCKER_IMAGE);

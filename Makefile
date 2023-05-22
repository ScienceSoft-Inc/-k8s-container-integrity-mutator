BINARY_NAME := mutator_app
APP_NAME := mutator
DEMO_NAME := demo
NAMESPACE ?= default
RELEASE_NAME_SYSLOG ?= rsyslog
SYSLOG_HOST ?= $(RELEASE_NAME_SYSLOG)
SYSLOG_PORT ?= 514
SYSLOG_ENABLED ?= false

IMAGE_NAME ?= integrity-injector
GIT_COMMIT := $(shell git describe --tags --long --dirty=-unsupported --always || echo pre-commit)
IMAGE_VERSION ?= $(GIT_COMMIT)

# helm chart path
HELM_CHART_PATH   := helm-charts
HELM_CHART_SYSLOG := helm-charts/rsyslog

# Downloads the Go module.
.PHONY : vendor
vendor:
	@echo "==> Downloading vendor"
	go mod tidy
	go mod vendor

.PHONY : docker
docker: vendor
	@eval $$(minikube docker-env) ;\
    docker build -t $(IMAGE_NAME):latest  -t $(IMAGE_NAME):$(IMAGE_VERSION) -f Dockerfile .

.PHONY: helm
helm-mutator:
	@helm upgrade -i ${APP_NAME} \
		--namespace=$(NAMESPACE) \
		--create-namespace \
		--set image.repository=$(IMAGE_NAME)    \
		--set image.tag=$(IMAGE_VERSION) \
		--set sideCar.syslog.enabled=$(SYSLOG_ENABLED) \
		--set sideCar.syslog.host=$(SYSLOG_HOST) \
		--set sideCar.syslog.port=$(SYSLOG_PORT) \
		$(HELM_CHART_PATH)/$(IMAGE_NAME)

helm-demo:
	@helm upgrade -i ${DEMO_NAME} \
		--namespace=$(NAMESPACE) \
		--create-namespace \
		--set rsyslog.enabled=$(SYSLOG_ENABLED) \
		$(HELM_CHART_PATH)/demo-app-to-inject

helm-syslog:
	helm upgrade -i ${RELEASE_NAME_SYSLOG} \
		--set fullnameOverride=$(SYSLOG_HOST) \
		--set service.port=$(SYSLOG_PORT) \
		$(HELM_CHART_SYSLOG) --wait

.PHONY : tidy
tidy: ## Cleans the Go module.
	@echo "==> Tidying module"
	@go mod tidy

.PHONY : build
build:
	go build -o ${BINARY_NAME} cmd/main.go

.PHONY : run
run:
	go build -o ${BINARY_NAME} cmd/main.go
	./${BINARY_NAME}

## Cleaning build cache.
.PHONY : clean
clean:
	go clean
	rm ${BINARY_NAME}

## Compile the binary.
compile-all: windows-32bit windows-64bit linux-32bit linux-64bit MacOS

windows-32bit:
	echo "Building for Windows 32-bit"
	GOOS=windows GOARCH=386 go build -o bin/${BINARY_NAME}_32bit.exe cmd/main.go

windows-64bit:
	echo "Building for Windows 64-bit"
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}_64bit.exe cmd/main.go

linux-32bit:
	echo "Building for Linux 32-bit"
	GOOS=linux GOARCH=386 go build -o bin/${BINARY_NAME}_32bit cmd/main.go

linux-64bit:
	echo "Building for Linux 64-bit"
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME} cmd/main.go

MacOS:
	echo "Building for MacOS X 64-bit"
	GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}_macos cmd/main.go

.PHONY: kind-load-images
kind-load-images:
	@kind load docker-image $(IMAGE_NAME):$(IMAGE_VERSION)

.PHONY: uninstall-all
uninstall-all: uninstall-mutator uninstall-demo

.PHONY: uninstall-mutator
uninstall-mutator:
	@-helm uninstall mutator

.PHONY: uninstall-demo
uninstall-demo:
	@-helm uninstall demo

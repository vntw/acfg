BINARY = acfg
ACSRV_BINARY = acServer
STRACKER_BINARY = stracker

GOARCH = amd64
BUILD_DIR = build/

.PHONY: all
all:
	$(MAKE) clean
	$(MAKE) client
	$(MAKE) static
	$(MAKE) darwin
	$(MAKE) linux
	$(MAKE) dummysrv
	#$(MAKE) docker

.PHONY: dummysrv
dummysrv:
	$(MAKE) darwin-dummysrv
	$(MAKE) linux-dummysrv

.PHONY: darwin-dummysrv
darwin-dummysrv:
	cd dummysrv/acserver && \
	GOOS=darwin GOARCH=${GOARCH} go build -o ${ACSRV_BINARY}-darwin-${GOARCH} . && \
	cd ../stracker && \
	GOOS=darwin GOARCH=${GOARCH} go build -o ${STRACKER_BINARY}-darwin-${GOARCH} .

.PHONY: linux-dummysrv
linux-dummysrv:
	cd dummysrv/acserver/ && \
	GOOS=linux GOARCH=${GOARCH} go build -o ${ACSRV_BINARY}-linux-${GOARCH} . && \
	cd ../stracker/ && \
	GOOS=linux GOARCH=${GOARCH} go build -o ${STRACKER_BINARY}-linux-${GOARCH} .

.PHONY: clean
clean:
	-rm -r build/; \
	mkdir build/

.PHONY: client
client:
	cd client/ && \
	yarn prebuild && \
	yarn build-dist

.PHONY: static
static:
	cd server && fileb0x b0x.toml

.PHONY: darwin
darwin:
	cd server/main/ && \
	GOOS=darwin GOARCH=${GOARCH} go build -o ../../build/${BINARY}-darwin-${GOARCH} .

.PHONY: linux
linux:
	cd server/main/ && \
	GOOS=linux GOARCH=${GOARCH} go build -o ../../build/${BINARY}-linux-${GOARCH} .

.PHONY: docker
docker:
	docker build --no-cache -t acfg . && \
	docker save acfg -o build/acfg-docker-image

.PHONY: test
test:
	cd server/ && \
	go test ./... && \
	cd ../client/ && \
	yarn test

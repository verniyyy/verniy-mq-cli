GOVERSION:=$(shell go version)
GOOS:=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH:=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
BUILD_DIR:=target/build
APPNAME:=$(shell cat appname)
APPNAMEPACK:=github.com/verniyyy/verniy-mq-cli/cmd.appname=$(APPNAME)
VERSIONPACK:=github.com/verniyyy/verniy-mq-cli/cmd.version=$(shell git describe --tags --abbrev=0)

dep:
	go mod download

build: dep
	rm -f $(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(APPNAME)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-X $(APPNAMEPACK) -X $(VERSIONPACK)" -o $(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(APPNAME)

multi-build: dep
	GOARCH=amd64 GOOS=darwin go build -ldflags="-X $(APPNAMEPACK) -X $(VERSIONPACK)" -o $(BUILD_DIR)/darwin-amd64/$(APPNAME)
	GOARCH=amd64 GOOS=linux go build -ldflags="-X $(APPNAMEPACK) -X $(VERSIONPACK)" -o $(BUILD_DIR)/linux-amd64/$(APPNAME)
	GOARCH=amd64 GOOS=windows go build -ldflags="-X $(APPNAMEPACK) -X $(VERSIONPACK)" -o $(BUILD_DIR)/windows-amd64/$(APPNAME)

	GOARCH=arm64 GOOS=darwin go build -ldflags="-X $(APPNAMEPACK) -X $(VERSIONPACK)" -o $(BUILD_DIR)/darwin-arm64/$(APPNAME)
	GOARCH=arm64 GOOS=linux go build -ldflags="-X $(APPNAMEPACK) -X $(VERSIONPACK)" -o $(BUILD_DIR)/linux-arm64/$(APPNAME)

install: build
	cp $(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(APPNAME) /usr/local/bin

run: build
	VERNIY_MQ_PASSWORD=password $(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(APPNAME) -H localhost -u root

clean:
	go clean
	rm $(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(APPNAME)
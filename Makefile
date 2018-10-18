.SHELL = /bin/bash
APP = sops64


default: clean build all package

.PHONY: build
build:
	mkdir -p build
	go get -d
	go test
	go build -o build/${APP} *.go

all:
	go get github.com/mitchellh/gox
	mkdir -p build
	gox \
		-output="build/{{.OS}}_{{.Arch}}/${APP}"


package:
	$(shell rm -rf build/archive)
	$(shell rm -rf build/archive)
	$(eval UNIX_FILES := $(shell ls build | grep -v ${APP} | grep -v windows))
	$(eval WINDOWS_FILES := $(shell ls build | grep -v ${APP} | grep windows))
	@mkdir -p build/archive
	@for f in $(UNIX_FILES); do \
		echo Packaging $$f && \
		(cd $(shell pwd)/build/$$f && tar -czf ../archive/$$f.tar.gz ${APP}*); \
	done
	@for f in $(WINDOWS_FILES); do \
		echo Packaging $$f && \
		(cd $(shell pwd)/build/$$f && zip ../archive/$$f.zip ${APP}*); \
	done
	ls -lah build/archive/

clean:
	rm -rf build/

install:
	chmod +x build/${APP}
	sudo mv build/${APP} /usr/local/bin/${APP}

test: test_sops test_sops64
	@$(MAKE) test_clean

test_sops: test_clean
	@echo "Testing sops"
	@sops --encrypt tests/base64.yml > tests/tmp.yml
	@sops --decrypt tests/tmp.yml > tests/tmp_out.yml
	@bash -c 'diff -w <(sort tests/base64.yml) <(sort tests/tmp_out.yml)'

test_sops64: test_clean
	@echo "Testing sops64"
	@go run main.go --encrypt tests/plain.yml > tests/tmp_sops.yml
	@sops --decrypt tests/tmp_sops.yml > tests/tmp.yml
	@bash -c 'diff -EbwB <(sort tests/tmp.yml) <(sort tests/base64.yml)'
	
	@go run main.go -e tests/plain.yml > tests/tmp.yml
	@go run main.go -d tests/tmp.yml > tests/tmp_out.yml
	@bash -c 'diff -EbwB <(sort tests/plain.yml) <(sort tests/tmp_out.yml)'

test_clean:
	@rm -rf tests/tmp*


SHELL          := /bin/bash
VERSION        := v0.1.0
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
SRC            := $(wildcard *.go) $(wildcard src/*/*.go)

CENTOS_IMAGE=docker-go-pkg-build-centos6
CENTOS_CONTAINER_NAME=docker-go-pkg-build-centos6-$(shell date +%s)

all: cloudwatch-metric2

cloudwatch-metric2: go-get $(SRC)
	GOPATH=$(RUNTIME_GOPATH) go build -a -tags netgo -installsuffix netgo -o cloudwatch-metric2
ifeq ($(GOOS),linux)
	[[ "`ldd cloudwatch-metric2`" =~ "not a dynamic executable" ]] || exit 1
endif

go-get:
	go get github.com/aws/aws-sdk-go

package: clean cloudwatch-metric2
	gzip -c cloudwatch-metric2 > cloudwatch-metric2-$(VERSION)-$(GOOS)-$(GOARCH).gz

package\:linux:
	docker run --name $(CENTOS_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(CENTOS_IMAGE) make -C /tmp/src package:linux:docker
	docker rm $(CENTOS_CONTAINER_NAME)

package\:linux\:docker: package
	mv cloudwatch-metric2-*.gz pkg/

rpm:
	docker run --name $(CENTOS_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(CENTOS_IMAGE) make -C /tmp/src rpm:docker
	docker rm $(CENTOS_CONTAINER_NAME)

rpm\:docker: clean
	cd ../ && tar zcf cloudwatch-metric2.tar.gz src
	mv ../cloudwatch-metric2.tar.gz /root/rpmbuild/SOURCES/
	cp cloudwatch-metric2.spec /root/rpmbuild/SPECS/
	rpmbuild -ba /root/rpmbuild/SPECS/cloudwatch-metric2.spec
	mv /root/rpmbuild/RPMS/x86_64/cloudwatch-metric2-*.rpm pkg/
	mv /root/rpmbuild/SRPMS/cloudwatch-metric2-*.src.rpm pkg/

docker\:build\:centos6:
	docker build -f docker/Dockerfile.centos6 -t $(CENTOS_IMAGE) .

clean:
	rm -f cloudwatch-metric2 *.gz
	rm -f pkg/*

SHELL          := /bin/bash
PROGRAM        := cloudwatch-metric2
VERSION        := v0.1.5
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
SRC            := $(wildcard *.go) $(wildcard src/*/*.go)

UBUNTU_IMAGE          := docker-go-pkg-build-ubuntu
UBUNTU_CONTAINER_NAME := docker-go-pkg-build-ubuntu-$(shell date +%s)
CENTOS_IMAGE          := docker-go-pkg-build-centos6
CENTOS_CONTAINER_NAME := docker-go-pkg-build-centos6-$(shell date +%s)

.PHONY: all
all: $(PROGRAM)

.PHONY: go-get
go-get:
	go get github.com/aws/aws-sdk-go

$(PROGRAM): $(SRC)
ifeq ($(GOOS),linux)
	GOPATH=$(RUNTIME_GOPATH) CGO_ENABLED=0 go build -ldflags "-X cwmetric2.version=$(VERSION)" -a -tags netgo -installsuffix netgo -o $(PROGRAM)
	[[ "`ldd $(PROGRAM)`" =~ "not a dynamic executable" ]] || exit 1
else
	GOPATH=$(RUNTIME_GOPATH) CGO_ENABLED=0 go build -ldflags "-X cwmetric2.version=$(VERSION)" -o $(PROGRAM)
endif

.PHONY: clean
clean:
	rm -f $(PROGRAM)

.PHONY: package
package: clean $(PROGRAM)
	gzip -c $(PROGRAM) > pkg/$(PROGRAM)-$(VERSION)-$(GOOS)-$(GOARCH).gz
	rm -f $(PROGRAM)

.PHONY: package/linux
package/linux:
	docker run \
	  --name $(UBUNTU_CONTAINER_NAME) \
	  -v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) \
	  make -C /tmp/src go-get package
	docker rm $(UBUNTU_CONTAINER_NAME)

.PHONY: deb
deb:
	docker run --name $(UBUNTU_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(UBUNTU_IMAGE) make -C /tmp/src deb/docker
	docker rm $(UBUNTU_CONTAINER_NAME)

.PHONY: deb/docker
deb/docker: clean go-get
	dpkg-buildpackage -us -uc
	mv ../$(PROGRAM)_* pkg/

package\:linux\:docker: package
	mv cloudwatch-metric2-*.gz pkg/

.PHONY: docker/ubuntu
docker/ubuntu: etc/Dockerfile.ubuntu
	docker build -f etc/Dockerfile.ubuntu -t $(UBUNTU_IMAGE) .

.PHONY: rpm
rpm:
	docker run --name $(CENTOS_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(CENTOS_IMAGE) make -C /tmp/src rpm/docker
	docker rm $(CENTOS_CONTAINER_NAME)

.PHONY: rpm/docker
rpm/docker: clean go-get
	cd ../ && tar zcf $(PROGRAM).tar.gz src
	mv ../$(PROGRAM).tar.gz /root/rpmbuild/SOURCES/
	cp $(PROGRAM).spec /root/rpmbuild/SPECS/
	rpmbuild -ba /root/rpmbuild/SPECS/$(PROGRAM).spec
	mv /root/rpmbuild/RPMS/x86_64/$(PROGRAM)-*.rpm pkg/
	mv /root/rpmbuild/SRPMS/$(PROGRAM)-*.src.rpm pkg/

.PHONY: docker/centos
docker/centos:
	docker build -f etc/Dockerfile.centos -t $(CENTOS_IMAGE) .

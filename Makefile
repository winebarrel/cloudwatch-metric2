SHELL          := /bin/bash
PROGRAM        := cloudwatch-metric2
VERSION        := 0.1.5
GOOS           := $(shell go env GOOS)
GOARCH         := $(shell go env GOARCH)
RUNTIME_GOPATH := $(GOPATH):$(shell pwd)
SRC            := $(wildcard *.go) $(wildcard src/*/*.go)

TRUSTY_IMAGE          := docker-go-pkg-build-trusty
TRUSTY_CONTAINER_NAME := docker-go-pkg-build-trusty-$(shell date +%s)
XENIAL_IMAGE          := docker-go-pkg-build-xenial
XENIAL_CONTAINER_NAME := docker-go-pkg-build-xenial-$(shell date +%s)
CENTOS6_IMAGE          := docker-go-pkg-build-centos6
CENTOS6_CONTAINER_NAME := docker-go-pkg-build-centos6-$(shell date +%s)

.PHONY: all
all: $(PROGRAM)

.PHONY: go-get
go-get:
	go get github.com/aws/aws-sdk-go

$(PROGRAM): $(SRC)
ifeq ($(GOOS),linux)
	GOPATH=$(RUNTIME_GOPATH) CGO_ENABLED=0 go build -ldflags "-X cwmetric2.version=v$(VERSION)" -a -tags netgo -installsuffix netgo -o $(PROGRAM)
	[[ "`ldd $(PROGRAM)`" =~ "not a dynamic executable" ]] || exit 1
else
	GOPATH=$(RUNTIME_GOPATH) CGO_ENABLED=0 go build -ldflags "-X cwmetric2.version=v$(VERSION)" -o $(PROGRAM)
endif

.PHONY: clean
clean:
	rm -f $(PROGRAM)

.PHONY: package
package: clean $(PROGRAM)
	gzip -c $(PROGRAM) > pkg/$(PROGRAM)-v$(VERSION)-$(GOOS)-$(GOARCH).gz
	rm -f $(PROGRAM)

.PHONY: package/linux
package/linux:
	docker run \
	  --name $(TRUSTY_CONTAINER_NAME) \
	  -v $(shell pwd):/tmp/src $(TRUSTY_IMAGE) \
	  make -C /tmp/src go-get package
	docker rm $(TRUSTY_CONTAINER_NAME)

.PHONY: deb/trusty
deb/trusty:
	docker run --name $(TRUSTY_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(TRUSTY_IMAGE) make -C /tmp/src deb/docker
	docker rm $(TRUSTY_CONTAINER_NAME)

.PHONY: deb/xenial
deb/xenial:
	docker run --name $(XENIAL_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(XENIAL_IMAGE) make -C /tmp/src deb/docker
	docker rm $(XENIAL_CONTAINER_NAME)

.PHONY: deb/docker
deb/docker: clean go-get
	dpkg-buildpackage -us -uc
	mv ../$(PROGRAM)_* pkg/
	cd pkg/ && mv $(PROGRAM)_$(VERSION)_amd64.deb $(PROGRAM)_$(VERSION)-`lsb_release -sc`_amd64.deb

package\:linux\:docker: package
	mv cloudwatch-metric2-*.gz pkg/

.PHONY: docker/trusty
docker/trusty: etc/Dockerfile.trusty
	docker build -f etc/Dockerfile.trusty -t $(TRUSTY_IMAGE) .

.PHONY: docker/xenial
docker/xenial: etc/Dockerfile.xenial
	docker build -f etc/Dockerfile.xenial -t $(XENIAL_IMAGE) .

.PHONY: rpm
rpm:
	docker run --name $(CENTOS6_CONTAINER_NAME) -v $(shell pwd):/tmp/src $(CENTOS6_IMAGE) make -C /tmp/src rpm/docker
	docker rm $(CENTOS6_CONTAINER_NAME)

.PHONY: rpm/docker
rpm/docker: clean go-get
	cd ../ && tar zcf $(PROGRAM).tar.gz src
	mv ../$(PROGRAM).tar.gz /root/rpmbuild/SOURCES/
	cp $(PROGRAM).spec /root/rpmbuild/SPECS/
	rpmbuild -ba /root/rpmbuild/SPECS/$(PROGRAM).spec
	mv /root/rpmbuild/RPMS/x86_64/$(PROGRAM)-*.rpm pkg/
	mv /root/rpmbuild/SRPMS/$(PROGRAM)-*.src.rpm pkg/

.PHONY: docker/centos6
docker/centos6:
	docker build -f etc/Dockerfile.centos6 -t $(CENTOS6_IMAGE) .

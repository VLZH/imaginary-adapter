VERSION=0.0.0
TAG=vladimirzhid/imaginary-adapter:${VERSION}

build:
	docker build -t $(TAG) .

publish:
	docker push $(TAG)

.PHONY: build publish

all: build publish
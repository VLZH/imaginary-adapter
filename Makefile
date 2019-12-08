TAG=vladimirzhid/imaginary-adapter

build:
	docker build -t $(TAG) .

publish:
	docker push $(TAG)

.PHONY: build publish

all: build publish
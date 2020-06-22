VER=$(shell git rev-parse --short HEAD)

build-service:
	cd service && go build -a -o ./main .

run-service-a:
	SERVICE_NAME="service-a" \
	PORT=9091 \
	./service/main

run-service-b:
	SERVICE_NAME="service-b" \
	PORT=9092 \
	./service/main

build-img-service:
	cd service && docker build -t layer5/learn-layer5:latest -t layer5/learn-layer5:$(VER) .

image-push:
	docker push layer5/learn-layer5

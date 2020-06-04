build-service:
	cd service && go build -a -o ./main .

run-service: build-service
	./service/main

build-img-service:
	cd service && docker build -t layer5/sample-app-service:latest

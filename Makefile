.PHONY: test
test:
	go test -v ./... -count=1

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./build/http-server ./cmd/http-server/

.PHONY: docker-build
docker-build:
	docker build -t http-mpx -f ./cmd/http-server/Dockerfile .

.PHONY: docker-run
docker-run:
	docker run --name http-mpx -p 3000:3000 -d http-mpx:latest

.PHONY: clean
clean:
	docker rm http-mpx --force
	docker image rm http-mpx --force

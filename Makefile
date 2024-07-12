all: build

clean:
	@rm -f eps

build:
	@go build -o eps main.go

image:
	@docker build --network host -f Dockerfile . -t eps:v0.1.0

cp-config:
	@cp config/example.json /var/lib/docker/volumes/eps-api_eps_volume/_data
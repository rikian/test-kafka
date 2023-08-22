bapp:
	docker build -f Dockerfile -t goapp .

bdeb:
	docker build -f Dockerfile.debug -t goapp .

runapp1:
	docker run -d -p 9095:9095 --name  goapp goapp

runapp2:
	docker run -d -p 9095:9095 -p 9096:9096 --name  goapp goapp

build:
	CGO_ENABLED=0 go build -o goapp .
all: build docker clean 

build:
	env GOOS=linux  go build -v spiderweb

docker:
	docker build . -t spiderweb:${TAG}
	docker tag spiderweb:${TAG} garybowers/spiderweb:${TAG}
	docker push garybowers/spiderweb:${TAG}

clean:
	rm ./spiderweb

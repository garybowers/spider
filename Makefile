.PHONY: deploy build tag push

build:
	docker build . -t ${IMG_NAME}:${TAG} --build-arg TF_VER=0.12.29 --build-arg GO_VER=1.14.4 --build-arg BAZEL_VER=3.2.0 --build-arg YARN_VER=1.22.4 

tag:
	docker tag ${IMG_NAME}:${TAG} ${REPO}/${IMG_NAME}:${TAG}

push:
	docker push ${REPO}/${IMG_NAME}:${TAG}

deploy:
	kubectl apply -f ./ 

run:
	docker run -it -p 8000:3000 -e USER=garybowers -u 1001 ${REPO}/${IMG_NAME}:${TAG}

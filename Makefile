.PHONY: build clean deploy

build:
	dep ensure -v
	
	env GOOS=linux go build  -o bin/main

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

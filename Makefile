node:
	cd static && yarn run production

node.dockerized:
	cd static && docker run -v $$(pwd):/www -w /www --rm -e NODE_OPTIONS=--openssl-legacy-provider node:17-alpine3.14 sh -c "yarn && yarn run production"

run.dockerized: node.dockerized
	docker-compose up -d

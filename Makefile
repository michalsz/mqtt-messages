build:
	go build -o server .

fmt:
	go fmt ./...

docker-build:
	docker build -t http_server-mqtt .

set-aws:
	export AWS_PROFILE=profile-name

docker-tag:
	docker tag http_server-mqtt:latest AWS_ACCOUNT_ID.dkr.ecr.eu-central-1.amazonaws.com/http-mqtt:latest

docker-push:
	docker push AWS_ACCOUNT_ID.dkr.ecr.eu-central-1.amazonaws.com/http-mqtt:latest

docker-container:
	docker run -p 3000:3000 http_server-mqtt
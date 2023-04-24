build:
	docker build . -t tagging_api_image:latest

run:
	docker run --name tagging_api -p 8080:8080 -d tagging_api_image:latest

stop:
	docker stop tagging_api

clean:
	docker rm tagging_api
	docker rmi --force tagging_api_image

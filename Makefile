run:
	docker run -d --name forum -p8080:8080 forum && echo "server started at http://localhost:8080/" --rm
build:
	docker build -t forum .
.PHONY: build
build:
	go build -o /dev/null ./...

.PHONY: run
run:
	go run .

.PHONY: run-scenario
run-scenario:
	curl -i -X POST http://localhost:8080/todos --data '{"title":"List up", "status": "done"}'
	curl -i -X POST http://localhost:8080/todos --data '{"title":"Go shopping", "status": "pending"}'
	curl -i http://localhost:8080/todos
	curl -i http://localhost:8080/todos/1
	curl -i -X PUT http://localhost:8080/todos/2 --data '{"title":"Go shopping", "status": "done"}'
	curl -i -X DELETE http://localhost:8080/todos/1
	curl -i http://localhost:8080/todos

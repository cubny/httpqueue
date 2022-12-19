.PHONY: test
test:
	GO111MODULE=on go vet ./...
	GO111MODULE=on go test -i ./...

.PHONY: coverage
coverage:
	GO111MODULE=on go test -v -failfast -race -count=1 -coverpkg=./... -coverprofile coverage.txt -covermode=atomic ./...
	GO111MODULE=on go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage: file://$(PWD)/coverage.html"

.PHONY: run
run:
	docker-compose up --no-cache

.PHONY: docs
docs:
	docker run -v $(CURDIR):/local -w /local  quay.io/goswagger/swagger generate spec -o ./docs/swagger.json
	docker run -v $(CURDIR):/local -w /local  quay.io/goswagger/swagger generate spec -o ./docs/swagger.yaml

#############################
# test 
############################
test:
	GO111MODULE=on go vet ./...
	GO111MODULE=on go test -i .
	GO111MODULE=on go test -v -short -failfast -race -count=1 \
		    -coverpkg=./... \
		    -coverprofile coverage.txt \
		    -covermode=atomic \
		    ./...
	GO111MODULE=on go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage: file://$(PWD)/coverage.html"
integrations:
	GO111MODULE=on go test -i .
	GO111MODULE=on go test -v ./...


#############################
# build, migrateand run
############################
build:
	go build -mod=vendor -o bin/cart github.com/cubny/cart/cmd/...
  
migrate:
	./bin/cart -data ./data/cart.db -migrate

run:
	./bin/cart -data ./data/cart.db

firstrun: build migrate run

#############################
# build, migrateand run with docker
############################
docker-build:
	docker build -t cubny/cart .

docker-migrate:
	docker run --name cart -v `pwd`/data:/app/data -it --rm  cubny/cart app/bin/cart -data /app/data/cart.db -migrate

docker-run:
	docker run --name cart -p 8080:8080 -p 8081:8081 -v `pwd`/data:/app/data -it --rm  cubny/cart app/bin/cart -data /app/data/cart.db

docker-firstrun: docker-build docker-migrate docker-run

#############################
# handy tools to test the api 
############################
addcart:
	curl -i -XPOST http://localhost:8080/carts -H "Authorisation: Key abcdef123456"

additem:
	curl -i -XPOST http://localhost:8080/carts/1/items -H "Authorisation: Key abcdef123456" -d '{"product_id":1, "quantity":1, "price":100.00}'

removeitem:
	curl -i -XDELETE http://localhost:8080/items/1 -H "Authorisation: Key abcdef123456" 


add5items:
	curl -i -XPOST http://localhost:8080/carts/1/items -H "Authorisation: Key abcdef123456" -d '{"product_id":2, "quantity":1, "price":10.00}'
	curl -i -XPOST http://localhost:8080/carts/1/items -H "Authorisation: Key abcdef123456" -d '{"product_id":3, "quantity":3, "price":12.00}'
	curl -i -XPOST http://localhost:8080/carts/1/items -H "Authorisation: Key abcdef123456" -d '{"product_id":4, "quantity":1, "price":2.50}'
	curl -i -XPOST http://localhost:8080/carts/1/items -H "Authorisation: Key abcdef123456" -d '{"product_id":5, "quantity":1, "price":0.99}'

emptycart:
	curl -i -XDELETE http://localhost:8080/carts/1/items -H "Authorisation: Key abcdef123456" 

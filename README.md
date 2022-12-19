# httpqueue 
Schedule HTTP requests in the future

**httpqueue** exposes an HTTP API to schedule HTTP requests with the desired delay.
It is supposed to be used as an internal tool.

The API exposes 2 methods as follows:

1. schedule a timer
```
POST /timers
```
2. get a timer using the timer ID 
```
GET /timers/{timer_id}
```

The full API specification is documented in `docs/swagger.json`.  [swagger.json](https://github.com/cubny/httpqueue/blob/master/docs/swagger.json).
Use a Swagger UI or the [online editor](https://editor.swagger.io/) to explore the API Doc.

If your IDE of choice is Goland you can also experiment with the API using [http-client.http](https://github.com/cubny/httpqueue/blob/master/http-client.http) file.

## How to run the service
The [Makefile](https://github.com/cubny/httpqueue/blob/master/Makefile) contains all the tooling run, build, test, and start developing.
See the full list of available make targets using:
```bash
make help
```

### a. Running the service using Docker
It's possible to run the service without installing any go tools thanks to Docker and Docker Compose.
```bash
make docker-run
```

Now the API service is available at http://localhost:8080/

Adjust the HTTP default port in the `docker-compose.yaml` file.

#### b. Running the service without Docker
```bash
make run 
```
Similarly, the API is available at http://localhost:8080/. 

Adjust the Port using `HTTP_PORT` environment variable 


## How to use the API
Of course, you can use the tools of your choice, but the project provides two convenient ways to just play with the API:

1- If you like command lines, the [Makefile](https://github.com/cubny/httpqueue/blob/master/Makefile) you can find some sub-commands 

2- If you prefer Goland, the [http-client.http](https://github.com/cubny/httpqueue/blob/master/http-client.http)

## Running the Tests
The code base includes two types of tests: unit tests and integration tests
**NOTE:** These tests are not meant to be run in containers

#### a. only unit tests ( with vet, race and coverage)
``` 
make test
```
#### b. integartions
``` 
make integrations
```

## Some notes about the code
- There are comments everywhere explaining the decisions I have made, so please make sure to read them if you want to understand why I made some choices over the others
- I implemented this service in Go, not because I believed it was the proper language for such a service, but because currently it is the language I am most productive with
- Instrumenting services is better to be done with critical metrics to the system and tracing every detail of the application. I didn'task include tracing for the lack of time, but included a sample prometheus metric to count 500 errors. the prometheus handler is exposed in 8081 port.

## Assumptions
- Since this is the first implementation of the service I didn'task include API versioning, it can be easily done in upstream 
layers like the reverse-proxy or the application load balancer for the first version. I assumed that later when the next version
of the API is about to release, the next version can become a separate application, or the versioning can make its way to the current application.

## Extra features for future
- It is only possible to add one item at a time to a cart, the payload should include price, quantity and product_id. down the road, it would be better to just pass the quantity and the product id and retrieve the price from the product micorservice
- with the current implementation, a user can have multiple carts, which is not ideal. Later it should be changed to let one user have only one open cart and multiple closed cart.
- the order microservice can get the details of a cart, but for now it is not possible to mark the cart as checked-out/ordered. now it is only possible to empty a cart.
- when the Auth service is ready, the auth client should be changed to reflect the actual users and access keys instead of stubbing the service.


## Why Redis?
- Redis can be underlying technology for both storage and queue. Using only one external dependency helps with the availability of the service. 


## TODO
- [ ] implement politeness for the task processor so in case of 429 it reads the response header or retries with exponential backoff
- [ ] Use swagger to generate API docs
- [ ] DNS caching for the HTTP Client
- [ ] Create a deny-list for hosts with many non-retryable errors
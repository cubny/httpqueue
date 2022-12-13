# cart
shopping cart service

**cart** is a RESTful API microservice extracted from a monolithic application. 
It has a basic authentication functionality, it uses sqlite3 for data storage and exposes metrics. 

The API exposes 4 methods as follows:
```
# create a cart
POST /carts
# add a product to a cart
POST /carts/:cartID/items
# remove an item from a cart
DELETE /items/:itemID
# empty a cart
DEETE /carts/:cartID/items
```
All methods expect a authorisation header in the format of `"Authorisation: Key {{key}}"`.

For more comprehensive usage of the methods, checkout the [http-client.http](https://github.com/cubny/cart/blob/master/http-client.http) file

## How to run the service
the [Makefile](https://github.com/cubny/cart/blob/master/Makefile) contains few subcommands to build, migrate and run the application.
both as a docker container or standalone. 

#### a. In container
- First run (builds the image and creates the database file)
```
make docker-firstrun
```
- Later you can run the application using the existing database
```
make docker-run
```
#### b. Standalone
- First run (creates the database file)
```
make firstrun
```
- Later you can run the application using the existing database
```
make run
```

## How to consume the API
Of course, you can use the tools of your choice, but the project provides two convenient ways to just play with the API:

1- If you like command lines, the [Makefile](https://github.com/cubny/cart/blob/master/Makefile) you can find some sub-commands 

2- If you prefer Goland, the [http-client.http](https://github.com/cubny/cart/blob/master/http-client.http)

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
- I could have mock the product service as well and when a product is being added to a cart, I could have retrieved the price from that service. I made this choice to keep things simple for the purpose of the project.
- For production, it is better to use environment variables, but I used flags to make testing and development easier. later on they should be changed to environment variables.
- Instrumenting services is better to be done with critical metrics to the system and tracing every detail of the application. I didn't include tracing for the lack of time, but included a sample prometheus metric to count 500 errors. the prometheus handler is exposed in 8081 port.

## Assumptions
- Since this is the first implementation of the service I didn't include API versioning, it can be easily done in upstream 
layers like the reverse-proxy or the application load balancer for the first version. I assumed that later when the next version
of the API is about to release, the next version can become a separate application, or the versioning can make its way to the current application.

## Extra features for future
- It is only possible to add one item at a time to a cart, the payload should include price, quantity and product_id. down the road, it would be better to just pass the quantity and the product id and retrieve the price from the product micorservice
- with the current implementation, a user can have multiple carts, which is not ideal. Later it should be changed to let one user have only one open cart and multiple closed cart.
- the order microservice can get the details of a cart, but for now it is not possible to mark the cart as checked-out/ordered. now it is only possible to empty a cart.
- when the Auth service is ready, the auth client should be changed to reflect the actual users and access keys instead of stubbing the service.


## Why Redis ZSET?
Common technologies for delayed jobs that I am aware of are Redis and SQS. the SQS delayed has a [limit of maximum 15m](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-delay-queues.html). 
I assumed that the http-request-scheduler is not limited to 15m but the timers could be set for any point in time in the future.
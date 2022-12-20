# httpqueue 
Make HTTP requests in the future.

**httpqueue** consists of different components that make it possible to schedule HTTP requests with the desired delay.
The scheduling is possible through a Rest API that `httpqueue` exposes. Under the hood the webhooks are queued and 
on the scheduled time, they are sent to workers. Workers make call to the webhooks and if the request fails they retry with exponential backoff. However, if a webhook fails because of a client error, i.e. 4XX (except 429), then the it fails permanently. 
`httpqueue` is fast because it uses Redis both for queueing and storing. It is also memory efficient, as it archives the finished webhooks in a Bloom Filter data structure.

**NOTE:** The scheduled HTTP requests are called *Timers* in `httpqueue`.

## Architecture 
Essentially `httpqueue` is a multi-component service. By default, for each component there is a separate process, but it's also possible to have all the components in one process.
For example, for development you can use one process and in production you can opt for deploying components separately with different autoscaling strategies. 

Since the timers are not only scheduled but also stored for API queries, the operation of scheduling and storing timers needs to be consistent and transactional: if the timer is stored in the datastore it must be also scheduled and if storing a timer fails it should not be scheduled and vice versa. 
Although `httpqueue` uses a single technology for both the database and message broker, it's still possible to configure different Redis server/clusters for each component, hence a distributed transaction is required to guarantee the consistency and atomicity.  

To achieve that, `httpqueue` follows the Outbox Pattern as shown in the image below (docs/httpqueue-architecture.png):

![httpqueue-architecture.png](docs%2Fhttpqueue-architecture.png)

### 1. HTTP API Server
It is a Rest API server that accepts requests for setting and getting timers. It stores the data in Redis and also enqueues the timer events.
### 2. Message Relay
Message relay is responsible to relay queued messages to the Message Broker. 
It does so by dequeuing the outbox queue, retrieving timers data and publishing them the message broker.
### 3. Message Broker and Workers
The Message broker supports delayed jobs. The server pulls tasks off the job queue and starts a worker goroutine for each task.
Once the worker is done with a task, the corresponding timer is archived in the datastore for space efficiency using a Bloom Filter.
The concurrency of workers is configurable. by default, it's 10.

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

## Tests
The code has 77.4% test coverage. 

Run the tests:
```bash 
make test
```
See the code coverage:
```bash
make coverage-report
```

## Some notes about the code
- There are some comments explaining the decisions I have made, so please make sure to read them if you want to understand why I made some choices over the others
- I implemented this service in Go, not because I believed it was the best tool for such a service, but because currently it is the language I am most productive with at the moment.
- Instrumenting services is better to be done with critical metrics to the system and tracing every detail of the application. I didn't include tracing because of time constraints, but included some prometheus metrics. the prometheus handler is exposed in 8081 port in all components.
- I chose Redis not only because it's a good option for both queueing and storing data but also because of simplicity. AS the result the service only has one dependency.

## Assumptions
- The service is supposed to be used internally, hence there is no throttling and no authentication.
- The timers are not required to be persisted permanently after the webhooks are called. Only the timer ID is kept. 
- The failed timers are retried with exponential backoff, providing that the response was retryable (5xx, 429 and others conditions)  
- It's possible to schedule a timer with zero delay.
- The timers are only expired when they are successfully called or permanently failed. In another word, if the requested delay is past due, even after some hours, the timer is not considered expired.
- A manage Redis cluster will be used, therefore I didn't configure Redis for persisting data on disk.

## Tech debts
- the permanently failed timers are not archived. they could be archived, or they could go to a DLQ queue for troubleshooting.
- Create a deny-list for hosts with many non-retryable errors
- Implement some component/integration tests.

## Extra features for future
- Politeness of the service: `httpqueue` should respect the rate limit of webhook's host. If the server responds with 429 with metadata, it should be respected. 
- Instrumenting the queue size, number of retries and the timers delays, for tuning concurrency and scalability. 
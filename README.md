# Foundry4 Interview Test 1

The JSON file in the repository has a list of products, with names and prices.
The task will be to implement the requirements in the given github issues.

The solution should be implemented by forking this repository, after which you are free to do as you like with the files contained in it.

## Guidelines
Depending on the requirements given in the issue, the purpose is to demonstrate how you would write production code, and to provide something to discuss.

For example:
* Demonstrating that you can write unit tests, or e2e tests in the given technology stack.
* Demonstrating code quality and craftsmanship
* Sensible use of third party libraries or tools

## Solution

`go-shopper` is a minimal online shopping implementation written in [Golang](https://golang.org/).

It exposes a bunch of REST API endpoints for different operations which could be easily integrated to a Web UI or another consumer capable of invoking REST Services in general.

> Note: All REST Endpoints are consumer agnostic and follows standard practices of GET, POST, PUT, DELETE on resources. Only exception is a Cart Checkout functionality which is special case.

### running the application

1. using go tool

```shell
$ git checkout https://github.com/boseabhishek/interview-projects

$ go run main.go
```

2. using docker-compose

```shell
$ docker-compose up
```

### running all tests


```shell
$ go test -v ./...
```


### running end to end tests

```shell
$ go test -v e2e_test.go
```

### exposed REST APIs

```javascript

POST /products          //save product
GET /products/{id}      //get product by id 
GET /products           //get all products
DELETE /products/{id}   //delete product by id


POST /cart              //save product to cart
GET /cart/{id}          // get product from cart by id
GET /cart               //get all products from cart
DELETE /cart/{id}       //remove product from cart

DELETE /checkout/{id}   //checkout product from cart by id

```

### some decisions taken:
- Golang maps has been used though I have tried to write a generic DB implemenation where the type can be swapped. Golang maps are good cndidate for this sort of operations as it's O(1) constant time for lookup so the performance is better. For production, MongoDB/Postgres could be used, depending upon complexity
- Standard routes -> handler -> services flow used
- pragmatic test are added as needed
- while writing E2E test using Go Convey, certain operations are avoided to preserve readability for e2e tests. They should be used by the consumer accordingly where use of goroutines and channels are recommended, if using Golang or Furture/Promise for async operations in Java/Scala/JS consumers. Please see as below:

    Scenario: Remove items from basket
	When products are chosen and added to cart, ideally they should be removed from inventory
	and vice-versa. The series of operations are:
	* fill up inventory 
		- POST fill products to inventory
	* add to cart
		- POST products to cart
	      DELETE same products from inventory (IN one consistent operation) (avoided for readability)
	* remove form cart  
		- DELETE products from cart 
		  POST same products back to inventory (IN one consistent operation) (avoided for readability)
    
** similar decisions have been taken for addition to cart, checkout etc.

- the code is a MVP version and can be refactored in future.


### answering additional requirements

1. You are expecting hundreds of thousands of people per day to view the store
contents, but only a small percentage of them would actually buy items, how
might you build this?

    Ans: 
    
    At code level, concurrency and parallelism using goroutine and channels will help. The consumer must choose to call the APIs asynchronously. Also, splitting into multople microservices and using one databse per microservice also helps.

    At infrastrucure level, horizontal pod autoscaling (HPA) when apps are deployed on K8s cluster helps in scalability. Also, a load balancer like Traefik and even Nginx helps. Services meshes are sometimes very helpful!

2. Imagine you wanted to collect real time dashboards of how many people are viewing any product at various points in time, how might you do this?

    Ans:

    Google Tag Manegr with Google Analytics helps in providing a good view of shoppers and carts. It also provides a view of dropout rates, items which were left on cart and metrics related to them.


3. The API call to the legacy warehousing system for finalizing an order after
checkout is unreliable, and frequently unavailable for periods of time. The
customer should see a successful checkout instantly however. How would you
handle this?

    Ans:

    Practising aysnchronous programming helps by expecting response in future. Caching using a fast document database e.g. MongoDB, in memeory DB like Redis or Elastic search indexing of data expected from the legacy systems helps, along with use of a circuit breaker pattern like Netflix Hystrix.
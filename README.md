# GO Challenges

This repo is going to be used as an example solution for any team member who is doing our Internal Go ramp-up exercises.

In a highlevel view, the 6 challenges are:

* Challenge A:
  * Setup a MySQL DB, either by installing a MySQL client, or by using a docker-compose file.

* Challenge B:
  * Create a DB repository to handle a CRUD of tasks and add a main.go file that run all operations. This must be implemented without any ORM.

* Challenge C:
  * Create a REST API to expose the CRUD of tasks using NO external library like gin-gonic or gorilla-mux. This must be done using only vanilla http.HandleFunc

* Challenge D:
  * Add gORM while keeping the vanilla implementation. The server API must choose which implementation to run based on an env variable. This is to enforce the team to implement decouple solutions.

* Challenge E:
  * Read [Practical GO](https://dave.cheney.net/practical-go/presentations/gophercon-singapore-2019.html) and then refactor your code based on best practices suggested in the book.

* Challeng F (Optional):
  * Implement a second API layer, this time using gRPC

## Setup

Create an `.env` file in the root directory. Use the `.env.example` file as an example.

## Running the code

### CLI

To run the "CLI" version, just run the following command:

```
go run cmd/cli/main.go
```

### REST API Server

To run the REST API Server, just run the following command:

```
go run cmd/server/main.go
```

### gRPC API Server

To run the gRPC API Server, just run the following command:

```
go run cmd/grpc/main.go
```


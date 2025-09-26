# clean-architecture-go

This is an example implementation of an E-Commerce basket 
using Uncle Bobs' Clean Architecture.

Please keep in mind that this is a demo application
and some things would be implemented differently in a real world application.

![Clean Architecture Diagram](CleanArchitecture.jpg)

Source: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

## E-Commerce Basket

### Use Case Diagram

![](usecase-diagram.svg)

### Data Model

todo

### Implementation

todo

## Start application

### Start application using Go

```shell
go run ./cmd/server
```

#### Start application using Docker

```shell
./build.sh
```

```shell
./run.sh
```

## Usage

### Web

Go to http://localhost:8080/

### REST API

```shell
curl http://localhost:8080/basket

curl -XPOST http://localhost:8080/basket/A12345/1
curl -XPOST http://localhost:8080/basket/A12345/1
curl -XPATCH http://localhost:8080/basket/A12345/10

curl -XPOST http://localhost:8080/basket/A12346/1
curl -XDELETE http://localhost:8080/basket/A12346

curl -XDELETE http://localhost:8080/basket
```


## Maintenance

### Recreate diagrams

```shell
./recreate-diagrams.sh
```

# Clean Architecture (Go)

This is an example implementation of an E-Commerce basket 
using Uncle Bobs' Clean Architecture written in Go.

![Clean Architecture Diagram](docs/CleanArchitecture.jpg)

Source: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

## E-Commerce Basket

### Use Case Diagram

![](docs/usecase-diagram.svg)

### ER Diagram

Simple Model:

![](docs/er-diagram-simple.svg)

With adapted Basket:

![](docs/er-diagram-adapted.svg)

### Implementation

The golang project structure is based on http://github.com/golang-standards/project-layout.

All code, except the `cmd` entrypoint, is "hidden" inside the `internal` directory.

Using a domain-driven design approach, the domains are separated inside the `internal/domain` directory.

The high-level layers "Entities" and "Use Cases" are combined inside the `business` directory.

The low-level layers "Adapters" and "Drivers" are separated.

#### Entities

The entities are stored inside this layer.

Also, there are Factory classes
to create new entities and Repository classes
to retrieve entities from the data layer and save entities into the data layer.

#### Use Cases

The use cases are stored inside this layer.

Every Use Case has a separate class which improves the readability and understandability.

Also, there are additional Service classes containing more business logic.

And for the output, there are some "Data Transfer Object" (DTO) classes. 

#### Adapters

The interface adapters are stored inside this layer.

The implemented adapters are a full REST API and a web adapters, but only for the "Show Basket" use case.

#### Drivers

The drivers are stored inside this layer.

The implemented drivers are an in-memory driver, MongoDB driver, and MySQL driver for the basket.

## Start application

### Start application using Go

#### In-Memory Driver (Default)
```shell
go run ./cmd/server
```

#### MongoDB Driver
```shell
DRIVER=mongodb go run ./cmd/server
```

#### MySQL Driver
```shell
# Set MySQL environment variables (optional - defaults shown)
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USERNAME=root
export MYSQL_PASSWORD=password
export MYSQL_DATABASE=ecommerce

# Start with MySQL driver
DRIVER=mysql go run ./cmd/server
```

**Note:** For MySQL, make sure to create the database and tables first by running:
```shell
mysql -u root -p < internal/domain/basket/drivers/mysql/schema.sql
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DRIVER` | `inmemory` | Database driver to use (`inmemory`, `mongodb`, `mysql`) |
| `HTTP_ADDR` | `localhost:8080` | HTTP server address |
| `MYSQL_HOST` | `localhost` | MySQL server host (when using MySQL driver) |
| `MYSQL_PORT` | `3306` | MySQL server port (when using MySQL driver) |
| `MYSQL_USERNAME` | `root` | MySQL username (when using MySQL driver) |
| `MYSQL_PASSWORD` | `password` | MySQL password (when using MySQL driver) |
| `MYSQL_DATABASE` | `ecommerce` | MySQL database name (when using MySQL driver) |

### Start application using Docker

#### Option 1: Using Docker Compose (Recommended)

**Quick Start Scripts:**

```shell
# Start with MySQL (includes database setup)
./start-with-mysql.sh

# Start with MongoDB
./start-with-mongodb.sh

# Start with in-memory (no database needed)
./start-with-inmemory.sh
```

**Manual Setup:**

Start the database services:

```shell
# Start MySQL and MongoDB services
docker-compose up -d

# Start the application with MySQL driver
DRIVER=mysql go run ./cmd/server

# Or start with MongoDB driver
DRIVER=mongodb go run ./cmd/server
```

Stop the services when done:

```shell
docker-compose down
```

#### Docker Services

The docker-compose.yaml file includes:

- **MySQL 8.0**: Available on port 3306
  - Database: `clean_architecture_go`
  - Root password: `password`
  - Automatically loads the schema.sql on first startup
  - Data persisted in `mysql_data` volume

- **MongoDB 8.0**: Available on port 27017
  - Database: `ecommerce`
  - No authentication required

#### Option 2: Using Docker Build Scripts

First build the docker image:

```shell
./docker-build.sh
```

Then run the docker container using the built image:

```shell
./docker-run.sh
```

## Usage

### Web

The web implementation only shows the basket. (first use case)

To view it, open http://localhost:8080/ in your web browser.

If you want to interact with the basket, please use the REST API described in the following section.

### REST API

The REST API fully implements all basket use cases with the following routes:

```shell
GET    /basket
POST   /basket/:productId
POST   /basket/:productId/:count
PATCH  /basket/:productId/:count
DELETE /basket/:productId
DELETE /basket
```

If you use `curl` in the shell, you can use [jq](https://github.com/jqlang/jq) to prettify the output.

#### Show Basket

```shell
curl http://localhost:8080/basket
```

#### Add first product A12345 with default count=1 to the basket

```shell
curl -XPOST http://localhost:8080/basket/A12345
```

#### Add more of product A12345 with count=2 to the basket

```shell
curl -XPOST http://localhost:8080/basket/A12345/2
```

#### Set count of the existing product A12345 in the basket to 10

```shell
curl -XPATCH http://localhost:8080/basket/A12345/10
```

#### Add product A12346 to the basket

```shell
curl -XPOST http://localhost:8080/basket/A12346/1
```

#### Delete product A12346 from the basket

```shell
curl -XDELETE http://localhost:8080/basket/A12346
```

#### Clear the basket

```shell
curl -XDELETE http://localhost:8080/basket
```

## Maintenance

### Recreate diagrams

The diagrams are built using `plantuml`.

To recreate them, just run:

```shell
./recreate-diagrams.sh
```

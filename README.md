# whobest

whobest is a powerful tool designed for performance testing of Redis databases. It allows users to benchmark different Redis-compatible databases, such as DragonflyDB, KeyDB, and standard Redis, by executing a series of operations and measuring their performance metrics.

## Features

* Multi-Database Support: Test and compare performance across multiple Redis-compatible databases.
* Configurable Parameters: Customize the number of operations, connection settings, and time duration for tests.
* Result Logging: Automatically logs the results of each test to a file for easy analysis.

## Requirements

* Go 1.23.2 or higher
* Redis or compatible databases (DragonflyDB, KeyDB)
* Go modules enabled

## Usage

To use whobest, clone the repository and navigate to the project directory:
```
go mod tidy
go run main.go
```


Results of DragonFlyDB vs KeyDB vs Redis :

(Tested on a i5-14400F cpu)

DragonflyDB:
  - Total Operations: 2244541
  - Errors: 0
  - Successful Operations: 2244541

KeyDB:
  - Total Operations: 2255358
  - Errors: 0
  - Successful Operations: 2255358

Redis:
  - Total Operations: 2054570
  - Errors: 0
  - Successful Operations: 2054570
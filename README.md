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

## Installation

To install whobest, clone the repository and navigate to the project directory:
```
go mod tidy
go run main.go
```

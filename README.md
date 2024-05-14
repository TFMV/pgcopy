# pgcopy

`pgcopy` is a Go-based solution designed to replicate data from one PostgreSQL database to another. It supports data migration across multiple tables, ensuring that data types are handled correctly during the copy process.

## Features

- Connects to both source and target PostgreSQL databases.
- Replicates data from specified tables in the source database to the target database.
- Handles various data types including integers, floats, strings, and booleans.
- Ensures type safety during data replication.
- Supports configuration via a YAML file.

## Requirements

- Go 1.21.2 or later
- PostgreSQL 12 or later

## Installation

### Clone the repository:

```bash
git clone https://github.com/TFMV/pgcopy.git
cd pgcopy
```

### Install dependencies

```bash
go mod tidy
```

### Build the executable

```bash
go build
```

## Configuration

```yaml
source:
  host: "localhost"
  port: "5432"
  user: "postgres"
  pass: "password"
  db: "source_db"
  isUnixSocket: false

target:
  host: "localhost"
  port: "5432"
  user: "postgres"
  pass: "password"
  db: "target_db"
  isUnixSocket: false

tables:
  - "customer"
  - "orders"
  - "products"
```

- source: Connection details for the source PostgreSQL database.
- target: Connection details for the target PostgreSQL database.
- tables: List of table names to be replicated from the source database to the target database.

## Usage

nsure that PostgreSQL is running and accessible based on the configuration provided in config.yaml.

Run the executable

```bash
./pgcopy
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.

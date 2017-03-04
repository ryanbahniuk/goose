# Goose
A lightweight database migration manager.

* * *

Goose is a runner for migrations written in raw SQL. It allows you to manage any SQL-based database operations via a simple command line interface. At the moment, it only supports Postgres.

## Installation

Install the executable via curl and move it into your path:

 ```sh
curl -O https://raw.githubusercontent.com/ryanbahniuk/goose/master/goose && chmod +x goose && mv goose /usr/local/bin
 ```

## Setup

Goose requires a configuration file to connect to your database. This should live in your application root and be called `.gooserc`. This file will most likely be `.gitignore`d so it can differ by environment and not leak credentials publicly. A sample `.gooserc` file is below (note that only Postgres is supported right now so the `db` value should not change):

```
[goose]
db = "postgres"
path = "/usr/local/bin/psql"
hostname = "localhost"
port = "5432"
name = "db_name"
username = "db_username"
password = "db_password"
```

## Usage

Goose has two main commands: `hatch` and `migrate`. These command should be run in your app directory that has a `.gooserc` file.

### hatch
This command will create a uniquely named migration file in your `migrations/` folder. If this folder does not exist yet, it will create it for you.

### migrate
This command will run all migrations after your last run migration. Goose stores your already run migrations in a local file called `.gaggle`. This will also most likely need to be in your `.gitignore` to keep database state in each environment.

## License

MIT © Ryan Bahniuk

# Async API example
### Prerequisite
[Install migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

#### Set env variables from .envrc
````bash
direnv allow
````

#### Check database connection
````bash
psql $DB_URL
asyncapi=# \d
asyncapi=# \c
````

#### Create migrations
````bash
name=init_schema make db_create_migration
````
# Async API example
### Prerequisite

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
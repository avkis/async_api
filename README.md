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

### Create (Sign up) testing user
````bash
curl -X POST -d '{ "email":"test@test.com", "password":"test"}' http://localhost:5000/auth/signup | jq
````

### Sign in testing user
````bash
curl -X POST -d '{ "email":"test@test.com", "password":"test"}' http://localhost:5000/auth/signin | jq
````

### Refresh token
````bash
curl -X POST -d '{ "refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTAwMCIsInN1YiI6ImQ0NGI4NzBkLTNiNDQtNDlhZC04ZDdlLWRjNTg5Y2MwYmIyNSIsImV4cCI6MTc0MTE5MDA5NCwiaWF0IjoxNzQwNzU4MDk0fQ.B0_eFXY1no7Kx-yflf4nIFn0tlsT1sIGn77N2l-UKLE"}' http://localhost:5000/auth/refresh | jq
````

### Ping the server by unauthorized user
````bash
curl http://localhost:5000/ping -v
````

### Ping the server by authorized user
````bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTAwMCIsInN1YiI6ImQ0NGI4NzBkLTNiNDQtNDlhZC04ZDdlLWRjNTg5Y2MwYmIyNSIsImV4cCI6MTc0MTE5MDA5NCwiaWF0IjoxNzQwNzU4MDk0fQ.B0_eFXY1no7Kx-yflf4nIFn0tlsT1sIGn77N2l-UKLE" http://localhost:5000/ping
````
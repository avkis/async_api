db_login:
	psql ${DB_URL}

db_test_login:
	psql ${DB_URL_TEST}

db_create_migration:
	migrate create -ext sql -dir migrations -seq $(name)

db_migrate:
	migrate -database ${DB_URL} -path migrations up

db_migrate_down:
	migrate -database ${DB_URL} -path migrations down

db_test_migrate:
	migrate -database ${DB_URL_TEST} -path migrations up

run_test:
	go test ./store
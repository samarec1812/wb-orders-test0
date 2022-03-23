.PHONY: build
build:
	$(info Building...)
	docker-compose build

.PHONY: run
run:
	$(info Running...)
	docker-compose up


.PHONY: drop-app
drop-app:
	$(info Droping...)
	docker-compose down

.PHONY: migrate-up
migrate-up:
	$(info Migrate up...)
	migrate -path ./schema -database 'postgres://postgres:garrix@localhost:5436/postgres?sslmode=disable' up

.PHONY: migrate-down
migrate-down:
	$(info Migrate down...)
	migrate -path ./schema -database 'postgres://postgres:garrix@localhost:5436/postgres?sslmode=disable' down
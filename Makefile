SHELL := /bin/bash

.PHONY: admin user

admin:
	source cmd/admin/sample.env && reflex -r '\.go' -s -- sh -c 'go run cmd/admin/main.admin.go'

stop-admin:
	pkill -f 'go run cmd/admin/main.admin.go'

user:
	source cmd/user/sample.env && reflex -r '\.go' -s -- sh -c 'go run cmd/user/main.user.go'

stop-user:
	pkill -f 'go run cmd/user/main.user.go'



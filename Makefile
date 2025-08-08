create_mocks:
	mockery --all --recursive --output ./mocks

linter:
	golangci-lint run internal/... --timeout 5m

temporal-linter:
	workflowcheck ./...

run-tests:
	go test ./internal/...

run-tests-cover:
	go test -cover ./internal/...

run-tests-cover-number:
	go test -coverprofile=coverage.out ./internal/... > /dev/null && go tool cover -func=coverage.out | grep total: | awk '{print $3}'
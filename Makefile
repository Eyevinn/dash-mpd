.PHONY: all
all: test check coverage

.PHONY: prepare
prepare:
	go mod tidy

.PHONY: test
test: prepare
	go test ./...

.PHONY: coverage
coverage:
	# Ignore (allow) packages without any tests
	go test ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out -o coverage.txt
	tail -1 coverage.txt

.PHONY: venv
venv: .venv/bin/activate

.venv/bin/activate:
	python3 -m venv .venv
	.venv/bin/pip install --upgrade pip
	.venv/bin/pip install pre-commit==4.2.0
	source .venv/bin/activate
	touch .venv/bin/activate

.PHONY: pre-commit
pre-commit: venv
	source .venv/bin/activate && .venv/bin/pre-commit run --all-files

.PHONY: check
check: prepare pre-commit
	golangci-lint run

.PHONY: update
update:
	go get -t -u ./...
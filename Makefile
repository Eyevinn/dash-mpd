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
venv:
	python3 -m venv venv
	venv/bin/pip install --upgrade pip
	venv/bin/pip install pre-commit==4.2.0

.PHONY: check
check: prepare
	golangci-lint run

venv:
	python3 -m venv venv
	venv/bin/pip install --upgrade pip
	venv/bin/pip install pre-commit==4.2.0
	venv/bin/pip install codespell

.PHONY: pre-commit-install
pre-commit-install: venv
	venv/bin/pre-commit install

.PHONY: pre-commit
pre-commit: venv
	venv/bin/pre-commit run --all-files

.PHONY: codespell
codespell: venv
	venv/bin/codespell -S testdata,references,xml,./venv -L nD

.PHONY: check
check: prepare pre-commit codespell

.PHONY: update
update:
	go get -t -u ./...
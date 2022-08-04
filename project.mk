.PHONY: build
build:
	go build -o out/terraform-provider-git

.PHONY: install
install:
	go install -v ./...

.PHONY: test
test:
	go test -v -cover -timeout=120s -parallel=4 ./internal/provider

.PHONY: single
single:
	go test -v -timeout=120s -run $(filter-out $@,$(MAKECMDGOALS)) ./internal/provider

.PHONY: docs
docs:
	go generate ./...

.PHONY: fmt
fmt:
	gofmt -s -w -e .

.PHONY: update
update:
	go get -u
	go mod tidy

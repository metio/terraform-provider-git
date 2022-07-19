.PHONY: build
build:
	go build -o out/terraform-provider-git

.PHONY: test
test:
	go test -v -cover ./internal/provider/

.PHONY: docs
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

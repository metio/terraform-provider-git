TEST?=$$(go list ./... | grep -v 'vendor')

.PHONY: build
build:
	go build -o out/terraform-provider-git

.PHONY: test
test:
	go test $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test -timeout=30s -parallel=4

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TEST) -v -timeout 120m

.PHONY: docs
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

NAMESPACE     = metio
NAME          = git
PROVIDER      = terraform-provider-${NAME}
VERSION       = 9999.99.99
OS_ARCH       ?= linux_amd64
XDG_DATA_HOME ?= ~/.local/share

out/${PROVIDER}: $(shell find internal -type f -name '*.go' -and -not -name '*test.go')
	mkdir --parents $(@D)
	go build -o out/${PROVIDER}

out/docs-sentinel: $(shell find internal -type f) $(shell find examples -type f -name '*.tf' -or -name '*.sh')
	mkdir --parents $(@D)
	go generate ./...
	touch $@

# see https://www.terraform.io/cli/config/config-file#implied-local-mirror-directories
out/install-sentinel: out/${PROVIDER}
	mkdir --parents $(@D)
	mkdir --parents ${XDG_DATA_HOME}/terraform/plugins/localhost/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp out/${PROVIDER} ${XDG_DATA_HOME}/terraform/plugins/localhost/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${PROVIDER}
	touch $@

tests/.terraform.lock.hcl: out/install-sentinel
	rm -rf ./tests/.terraform.lock.hcl
	terraform -chdir=./tests init

out/acceptance-sentinel: tests/.terraform.lock.hcl $(shell find tests -type f -name '*.tf')
	mkdir --parents $(@D)
	terraform -chdir=./tests apply -auto-approve -var="git_repo_path=${CURDIR}"
	touch $@

out/tests-sentinel: $(shell find internal -type f -name '*.go')
	mkdir --parents $(@D)
	go test -v -cover -timeout=120s -parallel=4 ./internal/provider
	touch $@

##@ hacking
.PHONY: install
install: out/install-sentinel ## install the provider locally

.PHONY: docs
docs: out/docs-sentinel ## generate the documentation

.PHONY: acceptance
acceptance: out/acceptance-sentinel ## run the acceptance tests

.PHONY: tests
tests: out/tests-sentinel ## run the integration tests

.PHONY: test
test: ## run specific tests
	go test -v -timeout=120s -run $(filter-out $@,$(MAKECMDGOALS)) ./internal/provider

.PHONY: format
format: ## format go code
	gofmt -s -w -e .

.PHONY: update
update: ## update all dependencies
	go get -u
	go mod tidy

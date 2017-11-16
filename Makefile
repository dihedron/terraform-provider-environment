TEST?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# defaul: build
default: test vet build

# buinl compiles the prvider and places it under $GOPATH/bin.
build: fmtcheck
	go install

# test runs unit tests.
test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

# testacc runs acceptance tests.
testacc: test-acceptance

# test-acceptance runs acceptance tests.
test-acceptance: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

# vet runs the Go source code static analysis tool `vet` to find any common errors.
vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

# fmt formats the Go source files.
fmt:
	gofmt -w $(GOFMT_FILES)

# fmtcheck checks if the sources are formatted according to Go rules.
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

# test-race runs the race checker
test-race: fmtcheck build
	TF_ACC= go test -race $(TEST) $(TESTARGS)

# cover runs the coverage and provides outputs as HTML
cover:
	@echo "coverage..."
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

test-docker:
	$(TEST_ENV) go test -v

clean:
	rm $(GOPATH)/bin/terraform-provider-environment

# disallow any parallelism (-j) for Make. This is necessary since some
# commands during the build process create temporary files that collide
# under parallel conditions.
.NOTPARALLEL:

.PHONY: build test testacc test-acceptance vet fmt fmtcheck errcheck vendor-status test-compile test-race cover test-docker


BINARY=terraform-provider-environment
TEST_ENV := VAR1=value1 VAR2=value2 VAR3="my test string"

.DEFAULT_GOAL: $(BINARY)

$(BINARY):
	go build -o bin/$(BINARY)

test:
	go test -v

docker_test:
	$(TEST_ENV) go test -v

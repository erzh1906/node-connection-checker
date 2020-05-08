.PHONY: build
build:
	for service in "web" "scheduler" ; do \
		CGO_ENABLED=0 GOOS=linux go build -o build/$$service node-connection-checker/cmd/$$service ; \
	done

.PHONY: clean
clean:
	rm -rf build

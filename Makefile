.PHONY install:
install:
	go install .

.PHONY generate:
generate:
	go generate ./...

GO_FILES := $(shell find . -name '*.go')
$(info GO_FILES = ${GO_FILES})

.PHONY:	test

datatable: $(GO_FILES)
	go build

test: 
	go test ./...

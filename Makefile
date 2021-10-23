GO_FILES := $(shell find . -name '*.go')
$(info GO_FILES = ${GO_FILES})

.PHONY:	test

datatable: $(GO_FILES)
	go build

test: export OUTPUT = $(shell tempfile)
test:
	@echo "####Output File: $${OUTPUT}"
	go test -v ./... > $${OUTPUT}
	cat $${OUTPUT}

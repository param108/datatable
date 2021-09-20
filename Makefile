GO_FILES := $(shell find . -name '*.go')
$(info GO_FILES = ${GO_FILES})
datatable: $(GO_FILES)
	go build

# Makefile for go-clay project
# Based on c-for-go Makefile template: https://github.com/xlab/c-for-go/wiki/Makefile-template

MANIFEST=clay.yml
PACKAGE_DIR=clay

all: generate

generate:
	go tool c-for-go $(MANIFEST)

clean:
	rm -f $(PACKAGE_DIR)/cgo_helpers.go $(PACKAGE_DIR)/cgo_helpers.h $(PACKAGE_DIR)/cgo_helpers.c
	rm -f $(PACKAGE_DIR)/const.go $(PACKAGE_DIR)/doc.go $(PACKAGE_DIR)/types.go
	rm -f $(PACKAGE_DIR)/clay.go

test:
	cd $(PACKAGE_DIR) && go build

.PHONY: all generate clean test


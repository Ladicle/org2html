PKGROOT := github.com/Ladicle/org2html

.PHONY: all fmt vet lint test

all: fmt vet lint test

fmt:
	go fmt $(PKGROOT)/org/...

vet:
	go vet -printfuncs Infof,Warningf,Errorf,Fatalf,Exitf,Logf $(PKGROOT)/org/...

lint:
	hack/golangci-lint.sh

test:
	go test $(PKGROOT)/org/...

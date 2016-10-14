GOENV := env GOPATH=$(shell pwd)

.PHONY: compile
compile: bootstrap
	$(GOENV) go build -o bin/zhrss zhrss

.PHONY: bootstrap
bootstrap:
	$(GOENV) go get github.com/gorilla/feeds
	$(GOENV) go get github.com/PuerkitoBio/goquery

.PHONY: dev-bootstrap
dev-bootstrap: bootstrap
	cp githooks/* .git/hooks/

.PHONY: clean
clean:
	rm -rf ./bin ./pkg

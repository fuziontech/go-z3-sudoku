Z3_REF ?= master

all: libz3.a test

clean:
	rm -rf vendor
	rm -f libz3.a

gofmt:
	@echo "Checking code with gofmt.."
	gofmt -s *.go >/dev/null

libz3.a: vendor/z3
	cd vendor/z3 && python scripts/mk_make.py --staticlib
	cd vendor/z3/build && ${MAKE}
	cp vendor/z3/build/libz3.a .
	cp libz3.a vendor/github.com/fuziontech/go-z3/

tidy:
	go mod tidy
	go mod vendor

vendor/z3: tidy
	mkdir -p vendor
	git clone https://github.com/Z3Prover/z3.git vendor/z3
	cd vendor/z3 && git reset --hard && git clean -fdx
	cd vendor/z3 && git checkout ${Z3_REF}

test: gofmt
	go test -v

.PHONY: all clean libz3.a test

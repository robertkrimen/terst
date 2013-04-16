.PHONY: test build release install clean

build:
	go build

release: build
	$(MAKE) -C .. $@

install: build
	go install

clean:
	rm -f terst-import

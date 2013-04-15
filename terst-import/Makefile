.PHONY: test build release install clean

build:
	go build

release: build
	godocdown -signature . > README.markdown

install: build
	go install

clean:
	rm -f terst-import
	rm -rf terst

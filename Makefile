.PHONY: test release install clean

test:
	go test -i
	go test

release: test
	(cd terst-import && godocdown -signature . > README.markdown) || false
	godocdown --signature > README.markdown

install: test
	go install
	$(MAKE) -C terst-import $@

clean:
	$(MAKE) -C terst-import $@

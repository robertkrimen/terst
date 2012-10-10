.PHONY: test synopsis release

test:
	go test .

synopsis:
	cd .hide && go test

release:
	godocdown > README.markdown

.PHONY: test synopsis release

test:
	go test .

synopsis:
	cd .hide && go test

release: test
	godocdown --signature > README.markdown

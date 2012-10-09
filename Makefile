.PHONY: test synopsis README

test:
	go test .

synopsis:
	cd .hide && go test

README:
	godocdown > README.md

.PHONY: test synopsis

test:
	go test .

synopsis:
	cd synopsis && go test

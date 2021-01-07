.PHONY: test
test:
	go test -v ./test

certs:
	(cd test/greta; certdeployer/certdeployer.json)
	(cd test/patty; certdeployer/certdeployer.json)
	(cd test/legacy; certdeployer/certdeployer.json)
test-greta:
	go build -o test/check-nomad-job cmd/CheckNomadJob.go
	cd test; \
	(cd greta; certdeployer certdeployer.json); \
	./check-nomad-job -job=martin-example -address="https://nomad01.greta-internal.hc.swops.cloud:4646" -caCert=greta/ca.pem -clientCert=greta/cert.pem -clientKey=greta/key.pem

test-patty:
	go build -o test/check-nomad-job cmd/CheckNomadJob.go
	cd test; \
	(cd patty; certdeployer certdeployer.json); \
	./check-nomad-job -j swops-plugin-csi-gluster -t service --address "https://nomad01.patty-production.awseuc1.swops.cloud:4646" --ca patty/ca.pem --cert patty/cert.pem --key patty/key.pem

test-legacy:
	go build -o test/check-nomad-job cmd/CheckNomadJob.go
	cd test; \
	(cd legacy; certdeployer certdeployer.json); \
	./check-nomad-job -j csi.serviceware.gluster -t csi-plugin --address "https://server01.nomadserver-internal.hc.sabio.de:4646" --ca legacy/ca.pem --cert legacy/cert.pem --key legacy/key.pem

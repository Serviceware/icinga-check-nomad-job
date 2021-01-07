test-greta:
	go build -o test/check-nomad-job cmd/CheckNomadJob.go
	cd test; \
	(cd greta; certdeployer certdeployer.json); \
	./check-nomad-job -job=martin-example -address="https://nomad01.greta-internal.hc.swops.cloud:4646" -caCert=greta/ca.pem -clientCert=greta/cert.pem -clientKey=greta/key.pem

test-patty:
	go build -o test/check-nomad-job cmd/CheckNomadJob.go
	cd test; \
	(cd patty; certdeployer certdeployer.json); \
	./check-nomad-job -job=swops-plugin-csi-gluster -type=csi-plugin -address="https://nomad01.patty-production.awseuc1.swops.cloud:4646" -ca-cert=patty/ca.pem -client-cert=patty/cert.pem -client-key=patty/key.pem

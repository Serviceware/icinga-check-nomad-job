# Icinga Check Nomad Job

## Goal

This project provides a [Nagios plugin](https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/pluginapi.html) to check if
a Nomad job is available. The plugin can be used in Icinga as a `CheckCommand`.

## Implement check

There are different types of Nomad jobs: service, system, batch (not supported at the moment). Each type implements its own check using
the [Check](internal/check.go) interface in the internal package. The check determines a
[status](https://icinga.com/docs/icinga-2/latest/doc/05-service-monitoring/#status) and prints some information about the job on `stdout`.

## Testing

To run local tests against one of our Nomad cluster (patty, greta, legacy)

* create certs: `make certs`
* run tests: `go test -v ./test`
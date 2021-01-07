# Icinga Check Nomad Job

## Goal

This project provides a [Nagios plugin](https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/pluginapi.html) to check if
a Nomad job is available. The plugin can be used as an Icinga ``

## Development

There are different types of Nomad jobs: service, system, batch (not supported at the moment). Each type has its own check implementation
in the internal package. Each check determines a 

## Testing

To run local tests against one of our Nomad cluster (patty, greta, legacy)

* create certs: `make certs`
* run tests: `go test -v ./test`
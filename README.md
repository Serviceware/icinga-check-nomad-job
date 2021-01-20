# Icinga Check Nomad Job

## Goal

This project provides a [Nagios plugin](https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/pluginapi.html) to check if
a Nomad job is available. The plugin can be used in Icinga as a `CheckCommand`.

## Contribute

There are different types of Nomad jobs: service, system, batch (not supported at the moment). Each type implements its own check using
the [Check](internal/check.go) interface in the internal package. The check determines a
[status](https://icinga.com/docs/icinga-2/latest/doc/05-service-monitoring/#status) and prints some information about the job on `stdout`.

## Build

```shell
go build -o check_nomad_job cmd/CheckNomadJob.go
```

## Testing

To run local tests against one of our Nomad cluster (patty, greta, legacy)

* create certs: `make certs`
* run tests: `go test -v ./test`

## Icinga config

```
object CheckCommand "nomad_job" {
    command = [ "/etc/icinga2/scripts/check_nomad_job" ]

    arguments = {
        "check" = {
            value = "$nomad_check$"
            skip_key = true
            order = -1
        }
        "--address" = "$nomad_address$",
        "--ca" = "$nomad_ca$",
        "--cert" = "$nomad_cert$",
        "--key" = "$nomad_key$",
        "--job" = "$nomad_job$",
        "--plugin" = {
            value = "$nomad_plugin$",
            set_if = {{ macro("$nomad_check$") == "csi-plugin" }}
        }
        "--unhealthyNodesWarning" = {
            value = "$nomad_unhealthy_nodes_warning$",
            set_if = {{ macro("$nomad_check$") == "csi-plugin" }}
        }
        "--unhealthyNodesCritical" = {
            value = "$nomad_unhealthy_nodes_critical$",
            set_if = {{ macro("$nomad_check$") == "csi-plugin" }}
        }
    }

    vars.nomad_address = "$nomad_address$"
    vars.nomad_ca = "$nomad_ca$"
    vars.nomad_cert = "$nomad_cert$"
    vars.nomad_key = "$nomad_key$"

    vars.nomad_check = "$nomad_check$"
    vars.nomad_job = "$nomad_job$"
    vars.nomad_plugin = ""
    vars.nomad_unhealthy_nodes_warning = 0
    vars.nomad_unhealthy_nodes_critical = 0
}
```

```
object Host "nomad-cluster" {
  import "generic-host"

  display_name = "nomad-cluster"

  check_command = "dummy"
  vars.dummy_state = 0
  vars.dummy_text = "Virtual dummy host I can attach passive checks for nomad-jobs to"

  vars.nomad_address = "https://server01.nomadserver-internal.hc.sabio.de:4646"
  vars.nomad_ca = "/etc/nomad/certificates/nomad-ca.pem"
  vars.nomad_cert = "/etc/nomad/certificates/client-cert.pem"
  vars.nomad_key = "/etc/nomad/certificates/client-key.pem"
}
```

```
object Service "job" {
    import "generic-service"
    host_name = "nomad-internal"
    
    check_command = "nomad_job"
    vars.nomad_check = "service"
    vars.nomad_job = "job name"
}
```

```
object Service "gluster-csi-plugin" {
  import "generic-service"
  host_name = "nomad-internal"

  check_command = "nomad_job"
  vars.nomad_check = "csi-plugin"
  vars.nomad_job = "gluster-csi-plugin"
  vars.nomad_plugin = "csi.serviceware.gluster"
  vars.nomad_unhealthy_nodes_warning = 1
  vars.nomad_unhealthy_nodes_critical = 2
}
```
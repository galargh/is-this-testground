# Is this Testground?

!["Is this Testground?" meme](is-this-testground.jpg)

This repo is a result of going through a mental excercise of trying to understand what the minimal set of components is that we could call [Testground](https://github.com/testground/testground). It is an answer to a hypothetical question: if I only had a week to build a Testground, what would I do?

## How to run the test?

```sh
docker compose up
```

## What happens during a test run?

1. The `coordinator` container starts. It is connected to the `control` bridge network.
1. 2 `runner` containers start. They are connected to the `control` bridge network and the `data` bridge network.
1. The `runner` containers modify their `data` network interfaces. They add 1000ms of latency.
1. The `runner` containers retrieve their own IP addresses for the `data` network and the `control` network.
1. The `runner` containers send their addresses as registration to the `coordinator` container through the `control` network.
1. Once the `coordinator` receives all the registrations it's been waiting for, it sends them to each `runner` container through `control` network.
1. The `runner` containers ping each other through the `data` network.
1. The `runner` containers send _done_ message to the `coordinator` container through the `control` network.
1. Once the `coordinator` receives all the _done_ messages it's been waiting for, it sends _shutdown` message to each `runner` container through `control` network and exits.
1. The `runnner` containers exit.

## Components

### Coordinator

The coordinator covers the following responsibilities:
- waiting for all runners to register
- advertising the registration data of all runners to all runners
- waiting for all runners to report results

### Runner

The runner covers the following responsibilities:
- registering with the coordinator
- waiting for the coordinator to advertise the registration data of all runners
- running the test
- reporting results to the coordinator

## Requirements

- runners need to be able to communicate with the coordinator
- the coordinator needs to be able to communicate with the runners
- eth0 is the control network interface (eth1 happens to be the data network interface but it does not have to be a hard requirement)

## Observations

- if runners implemented a healthcheck endpoint, the coordinator could monitor the health of the runners without assuming any specific environment (e.g. Docker)
- this example uses docker compose to set up the test but it could be "setup" agnostic
- neighter the coordinator nor runners have to be docker containers
- the Testground composition we know today could become a competitor to docker compose
- the network interface manipulation could be abstracted away in the SDK (and registration, healthz, result reporting, etc.)
- the runner doesn't even have to be a single process, it could be a cluster of processes
- how does the coordinator know when the test can be started? i.e. what is all the runners? would giving a number of registrations to expect be enough?
- the control network interface name could be configurable too

## Next steps

There's no next steps. This is just an excercise to understand what minimal Testground could be and how it could split responsibilities and abstractions.

# Actual Testground

## Breakdown by function

| function | component | description | required | complexity | alternatives | value |
| --- | --- | --- | --- | --- | --- | --- |
| building | daemon | building test nodes | optional | complex because of the need to support multiple languages and runtimes (builders) | `docker build`, `go build` | ? |
| orchestration | daemon | starting and stopping test nodes | optional | complex because of the need to support multiple environemnts (runners) | `docker compose`, `nomad`, `kubernetes` | ? |
| network setup | daemon | creating networks over which test nodes communicate | optional | only applicable to a subset of environments | `docker network`, `cni` | ? |
| synchronisation | sync | facilitating communication between test nodes | essential | ? | this is what makes a test run in a distributed environment |
| network configuration | sidecar | modifying network characteristics | essential | complex because of the need to support multiple environments (runners) | `tc`, `iptables` | an abstraction that differentiates testground |
| collection | daemon | collecting test results (from the sync) | essential | ? | restults reporting is a pretty important part of testing framework |
| metrics | sync | collecting test node application, runtime metrics | nice to have | `prometheus` | makes performance tests possible |
| logging | daemon | collecting test node application logs | optional | `docker logs`, `fluentd` | makes debugging tests easier |

## Breakdown by component

| component | function |
| --- |
| cli | init, wait |
| daemon | building, orchestration, network setup, collection |
| sync | synchronisation, metrics |
| sidecar | network configuration |
| test node | test logic |
| sdk | glue |

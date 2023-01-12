# Chaos Toolkit with Pumba

## Run

```sh
docker compose up -d
# Wait for daemons to be ready
./bin/chaos run experiments/delay.json
docker compose down
```

## Components

- [Chaos Toolkit](https://chaostoolkit.org/)
- [Pumba](https://github.com/alexei-led/pumba)

## Reading

- [How to run Pumba in Kubernetes under control of the Chaos Toolkit
](https://github.com/chaosiq/chaosiq/blob/4fa5e15e1f5fed1daec325a1bda3000fa4410c0c/pumba-kubernetes-integration/pumba-kubernetes.md)

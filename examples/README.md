# Examples

## Bckupr Only

```shell
docker compose -f docker-compose.yml
```

## Bckupr Read/Write

```shell
docker compose -f docker-compose.yml -f docker-compose.rw.yml
```

## Bckupr - Metrics

```shell
docker compose -f docker-compose.yml -f docker-compose.metrics.yml
```

1. Open Grafana UI
1. Add Prometheus Datasource: `Connections` -> `Add new connection` -> Search: `Prometheus` -> `Add New Datasource`
1. Set server URL to `http://prometheus:9090`


## Bckupr - Notifications

```shell
docker compose -f docker-compose.yml -f docker-compose.notifications.yml
```

## Bckupr - Following Dependancies

```shell
docker compose -f docker-compose.yml -f docker-compose.dependancies.yml
```
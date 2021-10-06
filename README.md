# lb

a reverse proxy load-balancing server, It implements the Weighted Round Robin Balancing algorithm referenced [here](https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35).

[Installing](#installing) | [How to Use](#how-to-use) | [Configuration](#configuration)

## Installing

## How to Use
After installing, run the command below to start the server: 
```sh
lb <PATH TO CONFIG FILE>
``` 
Now the server proxies every request to downstream backends.

## Configuration
Configuration format is `YAML` and path defaults to `lb.yml`.

### Available Values
| **Name**             | **Description**                                                               | **Default** |
|----------------------|-------------------------------------------------------------------------------|:-----------:|
| `port`               | Port the server will listen on                                                |    `3000`   |
| `retries`            | Amount retries per request                                                    |     `1`     |
| `health.path`        | URL path to listen to for liveness checks                                     |     `/`     |
| `health.interval`    | Interval to perform the liveness checks                                       |    `30s`    |
| `backends.[].host`   | Downstream Server URL to proxy requests to (ex. `localhost:4023`)             |     `-`     |
| `backends.[].weight` | Weight to influence the amount of requests received by the downstream backend |     `1`     |

Example Configuration:

```yaml
port: 5000
retires: 5
health:
  path: /health
  interval: 30s
backends:
  - host: "localhost:5000"
    weight: 2
  - host: "localhost:5001"
    weight: 1
```

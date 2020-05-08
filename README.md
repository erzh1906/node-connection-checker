# Check connectivity between nodes and send results to statsd

**Installation**

Docker Compose is the easiest way to try:

```bash
git clone https://git.aviata.team/devops/node-connection-checker.git
cd node-connection-checker
docker-compose up --build
```

Wait few minutes and visit http://localhost:7080/dashboard (Graphite)

**Components**

* [web](https://git.aviata.team/devops/node-connection-checker/-/blob/master/cmd/web/main.go) is the simple web app returns 200 ok.
* [scheduler](https://git.aviata.team/devops/node-connection-checker/-/blob/master/cmd/scheduler/main.go) is the scheduler checks targets listed in configuration

**Requirements for local development:**
  - Golang 1.14.2
  - Docker

**Scheduler environment variables:**

  - `STATSD_ADDR:` StatsD server address (host:port). Default `127.0.0.1:8125`
  - `METRIC_PREFIX:` prefix of sending metrics. Default `production`
  - `METRIC_SOURCE:` prefix of sending metrics. Default `server`
  - `SENTRY_DSN:` raven Sentry DSN address. Default `https://123:abc@sentry.example.com:443/1`

**Scheduler config ([example](https://git.aviata.team/devops/node-connection-checker/-/blob/master/pkg/scheduler/scheduler.yml)):**

  - `check_interval:` interval for periodic checks (seconds). Default `5`
  - `check_timeout_ms:` check timeout (milliseconds). Default `1500`
  - `targets:` list of targets to check.

**Metrics example (without statsd prefixes):**

  - `production.server.*.latency` latency between scheduler and remote targets (statsd timer type)
  - `production.server.*.http_200.count` success count between scheduler and remote targets (statsd counter type)
  - `production.server.*.error.count:` error count between scheduler and remote targets (statsd counter type)

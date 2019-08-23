# observ

observ is a minimal go application reporting [RED and USE](https://www.vividcortex.com/blog/monitoring-and-observability-with-use-and-red) prometheus metrics

RED metrics include:
* (Requests) request count and rate
* (Errors) HTTP error count and rate
* (Duration) HTTP request/response time 

USE measurements include:
* (Utilization) Percentage of workers that are busy
* (Saturation) Size of the work queue
* (Errors) Number of jobs that ended in error

RED metrics deal with requests and are generated automatically with [prometheus.InstrumentHander](https://godoc.org/github.com/prometheus/client_golang/prometheus#example-InstrumentHandler)

USE metrics deal with resources like CPU, MEM, disk, or in our case ... workers!

worker metrics are custom metrics and are instrumented using [prometheus.Gauge](https://godoc.org/github.com/prometheus/client_golang/prometheus#Gauge) and [prometheus.Counter](https://godoc.org/github.com/prometheus/client_golang/prometheus#Counter)

## Start observ and prometheus:

Note: default worker count and queue length is 4

```
docker-compose up
```

## Prometheus metrics

observ reports prometheus metrics at `/metrics` after the first request:

```
curl -s localhost:8111/req && curl -s localhost:8111/metrics
```

These ones are the RED Metrics:
```
http_request_duration_microseconds_sum{handler="/req"} 15.9
http_request_duration_microseconds_count{handler="/req"} 1
http_requests_total{code="200",handler="/req",method="get"} 1
```

and these are the USE metrics:
```
observ_worker_errors 0
observ_worker_saturation 0
observ_worker_utilization 0
```

## observ API

`GET` `/req` simulates http request traffic with the following params:

| param    | units     | default | description                       |
|----------|-----------|---------|-----------------------------------|
| duration | millisecs | 0       | ms delay before HTTP response     |
| httpcode | int       | 200     | http code of response             |
| worksecs | seconds   | 0       | seconds of work to add to queue   |
| workfail | bool      | false   | set true to simulate work failing |

## Examples:

Simulate a request that takes 472ms to return 302 but results in no added work

```bash
curl -si "localhost:8111/req?duration=472&httpcode=302"
```

Simulate a request that returns instantly but results in 6 seconds of added work.

```bash
curl -si "localhost:8111/req?worksecs=6"
```

Add a requtest that takes 150ms to return HTTP 201, and produces 6 seconds of work that ultimately fails.

```bash
curl -si "localhost:8111/req?duration=150&httpcode=201&worksecs=6&workfail=true"
```

Saturate workqueue (Note: responds HTTP 507)

```bash
while true; do curl -si localhost:8111/req?worksecs=11; sleep 1; done
```

## Testing

Use `localhost:8111/metrics` or the links below to verify prometheus metrics:

* [Request] rates go up the more often `/req` endpoint is hit
* [Errors] go up the more the `httpcode` param is set to 5xx
* [Duration] goes up when `duration` is set higher

* [Utilization] (number of workers in use) and,
* [Saturation] (number of jobs in queue) go up with work and work is created with each non-zero `worksecs` request.
* [Error] (count) increases each `workfail=true`

## TODO
* CD
* Serverless version to simulate job metrics
* Grafana

[Request]: http://localhost:9090/graph?g0.range_input=2h&g0.expr=rate(http_requests_total%7Bjob%3D%22observ%22%7D%5B5m%5D)&g0.tab=0'
[Errors]: http://localhost:9090/graph?g0.range_input=30m&g0.expr=http_requests_total%7Bjob%3D%22observ%22%2C%20code%3D~%225..%22%7D&g0.tab=0
[Duration]: http://localhost:9090/graph?g0.range_input=2h&g0.expr=http_request_duration_microseconds_sum%7Bjob%3D%22observ%22%7D%20%2F%20http_request_duration_microseconds_count%7Bjob%3D%22observ%22%7D&g0.tab=0
[Utilization]: http://localhost:9090/graph?g0.range_input=1h&g0.expr=observ_worker_utilization&g0.tab=0
[Saturation]: http://localhost:9090/graph?g0.range_input=1h&g0.expr=observ_worker_saturation&g0.tab=0
[Error]: http://localhost:9090/graph?g0.range_input=1h&g0.expr=observ_worker_errors&g0.tab=0

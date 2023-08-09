# Reverse proxy

This project represents a simple demo of a [reverse proxy](./cmd/proxy/proxy.go). A reverse proxy is a server that accepts a request from a client, forwards the request to another one of many other servers, and returns the results from the server that actually processed the request to the client as if the proxy server had processed the request itself.

Some of the features this project have are:

    - Inspecting and masking sensitive json data
    - Logging all input and output traffic
    - Block requests based on set of predefined rules

Proxy also adds header to your forwarded responses. `X-Proxy-Error` which can be `true` or `false`. It signals internal proxy errors proxy-ing or inspecting the request.

## Running Project

This project features a [json endpoint](./cmd/jsonendpoint/jsonendpoint.go) that allows easy testing of [reverse proxy](./cmd/proxy/proxy.go).

To run both proxy and the json endpoint all you need to do is `docker compose up`

1. It will create a container of proxy listening on port `8000`
2. It will create a container of jsonendpoint listening on port `8001`
3. Proxy will load [config.json](./config.json) file containing request blocking rules as well as config to which host and scheme to forward requests
4. To test the proxy masking you can send a `GET` request to proxy

```
curl localhost:8000
```

5. To test blocking you can send a `DELETE` request to proxy

```
curl -X POST localhost:8000
```

6. To test forwarding request without masking send a `PUT` request

```
curl -X PUT localhost:8000
```

## Running tests

Whole project is comprehensively using interfaces in order to allow easy unit testing.
Due to that whole project is also unit tested. For some tests snapshoting library [cupaloy](https://github.com/bradleyjkemp/cupaloy) was used. It generates snapshots in `./snapshots` dir.

To run tests

`go test ./...`

## Forwarding requests

Proxy will forward requests to a host specified in [config.json](./config.json) in field `forward_host`
It will use a scheme specified as `forward_scheme`

## Masking Rules

Default masking rules for PII (Personally identifiable information) are quite simple and if it were a real world project i would aim to use a more comprehensive set of detections instead of a couple of simple detections.
They are located [here](./internal/mask/classifier.go) and are easily extendible

## Blocking rules explained

As explained in [top comment](./internal/block/guards.go):
If this was a real world user facing software i would use [open policy agent rego](https://www.openpolicyagent.org/docs/latest/policy-language/), however my understanding of the task was that it was required to build some kind of rules mechanism myself.

Blocking rules are loaded from [config.json]([./config.json) `block` field
Blocking rules are designed in a way that you can compose different rules in order to build different combinations of blocking rules.

For example if you want to block all delete or posts requests you use this blocking rules:

```

{
    "forward_host": "jsonendpoint:8000",
    "forward_scheme": "http",
    "block": [
        [
            {
                "method": "POST"
            }
        ],
        [
            {
                "method": "DELETE"
            }
        ]
    ]
}

```

If on the other hand you would like to block all `POST` requests that start with path `/api` or any `DELETE` request

```

{
    "forward_host": "jsonendpoint:8000",
    "forward_scheme": "http",
    "block": [
        [
            {
                "method": "POST"
            },
            {
                "path": "/api"
            }
        ],
        [
            {
                "method": "DELETE"
            }
        ]
    ]
}

```

First level of block property acts as `OR` (`||`) and second level acts as `AND` (`&&`) when matching

### Possible blocks

All guards used for blocks are located [here](./internal/block/guards.go)

1. Method block

```

{
    "method": "DELETE"
}

```

2. Path block

```

{
    "path": "/api"
}

```

3. Query Parameter block

```

{
    "query_param": "userID"
    "value": "userID"
}

```

4. Header block

```

{
    "header": "Content-Type"
    "value": "text/html; charset=utf-8"
}

```

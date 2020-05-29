# elastic-opencensus-exporter-go

The `escensus` package provides an exporter for sending Opencensus stats and traces to Elastic.

## Install

```bash
$ go get -u github.com/newrelic/newrelic-opencensus-exporter-go
```

## Using the exporter

```go

package main

import (
    "github.com/ResultadosDigitais/elastic-opencensus-exporter-go/escensus"
    "go.opencensus.io/stats/view"
    "go.opencensus.io/trace"
)

func main() {
    exporter, err := escensus.NewElasticApmExporter(
    if err != nil {
        panic(err)
    }
    
    trace.RegisterExporter(exporter)
}
```

You need to set some envvars to identify your service:

| Attribute          | ENV                         | Description                                                     |
|--------------------|-----------------------------|-----------------------------------------------------------------|
| ServiceName        | ELASTIC_APM_SERVICE_NAME    | The Service Name that identifies your application               |
| ServiceVersion     | ELASTIC_APM_SERVICE_VERSION | Represents the version of your application                      |
| ServiceEnvironment | ELASTIC_APM_ENVIRONMENT     | Represents the environment of your application. eg: development |

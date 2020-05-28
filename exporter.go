package elasticexporter

import (
	"go.elastic.co/apm"
	"go.opencensus.io/trace"
)

// ElasticApmExporter is an implementation of opencensus/trace.Exporter that exports spans to Elastic APM.
type ElasticApmExporter struct {
	Tracer *apm.Tracer
}

// NewElasticApmExporter is a factory that will return the elasticApmExporter struct
func NewElasticApmExporter() *ElasticApmExporter {
	return &ElasticApmExporter{
		Tracer: apm.DefaultTracer,
	}
}

// ExportSpan reports s as a transaction or span to e.Tracer.
func (e *ElasticApmExporter) ExportSpan(s *trace.SpanData) {
	// TODO(axw) need to do better for the type -- should be based on attributes.
	// TODO(axw) pull out well-known attributes, like "http.*", and construct the appropriate context.
	// TODO(axw) add support for marks, translate annotations to marks?
	var traceOptions apm.TraceOptions
	if s.TraceOptions.IsSampled() {
		traceOptions = traceOptions.WithRecorded(true)
	}

	if s.HasRemoteParent || s.ParentSpanID == (trace.SpanID{}) {
		tx := e.Tracer.StartTransactionOptions(
			s.Name, "request",
			apm.TransactionOptions{
				Start:         s.StartTime,
				TransactionID: apm.SpanID(s.SpanID),
				TraceContext: apm.TraceContext{
					Trace:   apm.TraceID(s.TraceID),
					Span:    apm.SpanID(s.ParentSpanID),
					Options: traceOptions,
				},
			},
		)
		for key, value := range s.Attributes {
			tx.Context.SetLabel(key, value)
		}
		tx.Duration = s.EndTime.Sub(s.StartTime)
		tx.End()
	} else {
		// BUG(axw) if the parent is another span, then the transaction
		// ID here will be invalid. We should instead stop specifying the
		// transaction ID (requires APM Server 7.1+)
		transactionID := apm.SpanID(s.ParentSpanID)
		span := e.Tracer.StartSpan(
			s.Name, "request",
			transactionID,
			apm.SpanOptions{
				Start:  s.StartTime,
				SpanID: apm.SpanID(s.SpanID),
				Parent: apm.TraceContext{
					Trace:   apm.TraceID(s.TraceID),
					Span:    apm.SpanID(s.ParentSpanID),
					Options: traceOptions,
				},
			},
		)
		for key, value := range s.Attributes {
			span.Context.SetLabel(key, value)
		}
		span.Duration = s.EndTime.Sub(s.StartTime)
		span.End()
	}
}

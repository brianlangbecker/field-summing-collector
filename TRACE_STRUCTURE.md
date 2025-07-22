# Trace Structure Documentation

## Visual Trace Structure

```
root_span (appserver)
└── compute.v1.ComputeEngine/getData (appserver) [6.106s] ← TARGET FOR ENRICHMENT
    └── compute.v1.ComputeEngine/getData (compute) [6.105s]
        ├── caching.series.v1.Series/read (compute) [1.857ms]
        │   └── caching.series.v1.Series/read (series-cache) [0.8663ms]
        └── data.signal.v1.Signal/read (compute) [6.029s]
            └── data.signal.v1.Signal/read (datasource-proxy) [6.027s] ← SOURCE FOR SUMMING
                └── POST (datasource-proxy) [2.967ms]
                    └── POST /system/monitors/gauge (datasource-proxy) [2.645ms]
                        ├── SELECT seeq.public (appserver) [0.2662ms]
                        ├── SELECT seeq.public (appserver) [0.1667ms]
                        └── SELECT seeq (appserver) [0.1687ms]
```

## Trace Details

### Root Level
- **Service**: `appserver`
- **Span**: `root_span`
- **Child**: `compute.v1.ComputeEngine/getData`

### Main Operation (Target for Enrichment)
- **Service**: `appserver`
- **Span**: `compute.v1.ComputeEngine/getData`
- **Duration**: 6.106s
- **Trace ID**: `774717dab04497439c19dbda0be5ac82`
- **Span ID**: `bb4d247fed155b95`
- **Parent**: `d491d6a90824aafe`

### Compute Service Call
- **Service**: `compute`
- **Span**: `compute.v1.ComputeEngine/getData`
- **Duration**: 6.105s
- **Span ID**: `50f699f76be95bf2`
- **Parent**: `bb4d247fed155b95`

### Cache Operations
- **Service**: `compute`
- **Span**: `caching.series.v1.Series/read`
- **Duration**: 1.857ms
- **Span ID**: `b479aeb00406737e`

- **Service**: `series-cache`
- **Span**: `caching.series.v1.Series/read`
- **Duration**: 0.8663ms
- **Span ID**: `060387beb846d2fd`

### Signal Operations (Source for Summing)
- **Service**: `compute`
- **Span**: `data.signal.v1.Signal/read`
- **Duration**: 6.029s
- **Span ID**: `30bec16e4ff86558`

- **Service**: `datasource-proxy` ← **THIS IS WHAT WE SUM**
- **Span**: `data.signal.v1.Signal/read`
- **Duration**: 6.027s
- **Span ID**: `9dc44e5a101233dd`

### HTTP Operations
- **Service**: `datasource-proxy`
- **Span**: `POST`
- **Duration**: 2.967ms
- **Span ID**: `4f16d15cda9fe058`

### Database Operations
- **Service**: `appserver`
- **Span**: `SELECT seeq.public`
- **Duration**: 2.645ms
- **Span ID**: `e31709174d8065bb`

- **Service**: `appserver`
- **Span**: `SELECT seeq.public`
- **Duration**: 0.2662ms
- **Span ID**: `e83994ab6e3d254b`

- **Service**: `appserver`
- **Span**: `SELECT seeq`
- **Duration**: 0.1667ms
- **Span ID**: `0073b6922044cdac`

## OTEL Collector Configuration Logic

### Critical Assumption
**The collector assumes exactly ONE instance of the high-level `getData` operation per trace, with ALL `Signal/read` operations contained as descendants underneath it.**

This means:
- Each trace has exactly one `compute.v1.ComputeEngine/getData` span from `appserver`
- All `data.signal.v1.Signal/read` spans from `datasource-proxy` are descendants of this `getData` span
- The collector will sum ALL `Signal/read` spans in the trace and add the result to the single `getData` span

### What Gets Summed
```yaml
Select(
  spans,
  span.name == "data.signal.v1.Signal/read" and 
  resource.attributes["service.name"] == "datasource-proxy",
  duration / 1000000  # Convert ns to ms
)
```

This selects **ALL** `data.signal.v1.Signal/read` spans from the `datasource-proxy` service within the trace and sums their durations.

### Where It Gets Added
```yaml
where span.name == "compute.v1.ComputeEngine/getData" and 
      resource.attributes["service.name"] == "appserver"
```

The summed value gets added as an attribute to the **single** `compute.v1.ComputeEngine/getData` span from the `appserver` service.

### Result
The appserver's `getData` span will have these new attributes:
- `read_signal_sum_ms`: Sum of all Signal/read durations (6027ms in this example)
- `read_signal_max_end_ns`: Maximum end time of all Signal/read spans

### Important Limitations
⚠️ **This configuration will NOT work correctly if:**
- A trace has multiple `getData` operations at the same level
- `Signal/read` spans exist outside the `getData` operation hierarchy
- Multiple traces are processed together without proper trace grouping

The `groupbytrace` processor ensures traces are processed individually, maintaining the one-to-one relationship between `getData` and its descendant `Signal/read` spans.

## Usage with Trace Generator

To run the trace generator:
```bash
cd trace-generator
go mod tidy
go run main.go
```

This will generate traces matching the structure above and send them to the OTEL collector at `localhost:4318`.
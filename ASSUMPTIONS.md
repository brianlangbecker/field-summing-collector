# Field Summing Collector - Assumptions and Limitations

## Critical Assumptions

### 1. Single getData Operation Per Trace
The collector assumes **exactly one** `compute.v1.ComputeEngine/getData` operation per trace from the `appserver` service.

✅ **Supported:**
```
root_span (appserver)
└── compute.v1.ComputeEngine/getData (appserver) ← SINGLE getInstance
    └── ... descendant spans including Signal/read operations
```

❌ **NOT Supported:**
```
root_span (appserver)
├── compute.v1.ComputeEngine/getData (appserver) ← MULTIPLE getData
├── compute.v1.ComputeEngine/getData (appserver) ← MULTIPLE getData
└── ... other spans
```

### 2. Hierarchical Containment
All `data.signal.v1.Signal/read` spans from `datasource-proxy` must be **descendants** of the `getData` span.

✅ **Supported:**
```
getData (appserver)
└── ... intermediate spans
    └── data.signal.v1.Signal/read (datasource-proxy) ← DESCENDANT
```

❌ **NOT Supported:**
```
root_span (appserver)
├── getData (appserver)
└── data.signal.v1.Signal/read (datasource-proxy) ← SIBLING, NOT DESCENDANT
```

### 3. Service Name Consistency
- The target `getData` span must come from service `appserver`
- The source `Signal/read` spans must come from service `datasource-proxy`

### 4. Trace Grouping
The `groupbytrace` processor ensures all spans in a trace are processed together before transformation.

## Configuration Behavior

### What Gets Summed
```yaml
Select(
  spans,
  span.name == "data.signal.v1.Signal/read" and 
  resource.attributes["service.name"] == "datasource-proxy",
  duration / 1000000
)
```

This selects **ALL** matching spans within the trace and sums their durations.

### Where It Gets Added
```yaml
where span.name == "compute.v1.ComputeEngine/getData" and 
      resource.attributes["service.name"] == "appserver"
```

The summed value is added as an attribute to the **single** `getData` span.

## Failure Scenarios

### Multiple getData Operations
If a trace has multiple `getData` spans:
- Each `getData` span will receive the sum of **ALL** Signal/read spans in the trace
- This creates duplicate/incorrect summed values

### Signal/read Outside getData Hierarchy
If Signal/read spans exist outside the getData operation:
- They will still be included in the sum
- The summed value may not represent the actual work done by the getData operation

### Missing groupbytrace Processor
Without `groupbytrace`:
- Spans from different traces may be processed together
- Summed values may include spans from multiple traces
- Results will be incorrect

## Validation

To verify your traces match these assumptions:

1. **Check for single getData per trace:**
   ```
   count(spans where span.name == "compute.v1.ComputeEngine/getData" and service.name == "appserver") == 1
   ```

2. **Verify hierarchical containment:**
   - All Signal/read spans should have the getData span as an ancestor
   - Use trace visualization tools to confirm hierarchy

3. **Validate service names:**
   - getData spans should come from `appserver`
   - Signal/read spans should come from `datasource-proxy`

## Recommendations

1. **Use trace visualization** to confirm your trace structure matches assumptions
2. **Test with sample data** before production deployment
3. **Monitor for anomalies** in summed values that might indicate assumption violations
4. **Consider extending the configuration** if your traces have different patterns
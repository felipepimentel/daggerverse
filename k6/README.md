# k6 Dagger Module

This Dagger module provides integration with [k6](https://k6.io/), a modern load testing tool built for developer happiness. The module uses the [xk6-dashboard](https://github.com/grafana/xk6-dashboard) extension to provide real-time test metrics visualization.

## Features

- Run k6 load tests with customizable parameters
- Real-time metrics visualization through xk6-dashboard
- Configurable virtual users (VUs) and test duration
- Environment variable support
- HTML report generation
- JSON summary export
- Error logging

## Usage

### Basic Usage

```typescript
import { k6 } from "@felipepimentel/daggerverse/k6";

// Initialize k6
const client = k6();

// Run a simple test
const container = await client.run({
  workingDir: dag.host().directory("./tests"), // directory containing k6 scripts
  script: "script.js", // script file to execute
  vus: 1, // 1 virtual user
  duration: "30s", // 30 seconds duration
});
```

### Advanced Usage

```typescript
// Run with environment variables and more VUs
const container = await client.run({
  workingDir: dag.host().directory("./tests"),
  script: "load-test.js",
  env: ["API_URL=https://api.example.com", "TOKEN=secret123"],
  vus: 50, // 50 virtual users
  duration: "5m", // 5 minutes duration
});
```

## API Reference

### Constructor

The module accepts no parameters.

### Methods

#### `run(options: RunOptions)`

Runs k6 load tests with specified parameters.

Parameters:

- `workingDir`: Directory containing the k6 script files
- `script`: Name of the script file to execute
- `env` (optional): List of environment variables in the format `KEY=VALUE`
- `vus` (optional): Number of virtual users to simulate (default: 1)
- `duration` (optional): Duration of the test (default: "1s")

Returns a container with the test execution and results.

### Output Files

The module generates several output files in the container:

- `/output/report.html`: Interactive HTML dashboard with test metrics
- `/output/summary.json`: JSON file containing test summary
- `/output/errors.txt`: Text file with error logs

## Example k6 Script

```javascript
import http from "k6/http";
import { check, sleep } from "k6";

export default function () {
  const res = http.get("http://test.k6.io");
  check(res, {
    "status is 200": (r) => r.status === 200,
  });
  sleep(1);
}
```

## Complete Example

```typescript
import { k6 } from "@felipepimentel/daggerverse/k6";

export async function runLoadTest() {
  // Initialize k6
  const client = k6();

  // Run load test
  const container = await client.run({
    workingDir: dag.host().directory("./load-tests"),
    script: "api-test.js",
    env: ["BASE_URL=https://api.example.com", "API_KEY=your-api-key"],
    vus: 10, // 10 VUs
    duration: "2m", // 2 minutes
  });

  // Export results
  await container.directory("/output").export("./test-results");
}
```

## Testing

To test the module:

1. Create a k6 script file (e.g., `script.js`)
2. Run the test:

```typescript
const container = await client.run({
  workingDir: dag.host().directory("."),
  script: "script.js",
  vus: 1,
  duration: "10s",
});
```

3. Check the results in the `/output` directory of the container

## License

See [LICENSE](../LICENSE) file in the root directory.

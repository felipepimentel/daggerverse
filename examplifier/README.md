# Examplifier Module for Dagger

A Dagger module that automatically generates example pipelines and documentation for other Dagger modules using GPT. This module helps developers understand how to use Dagger modules by creating comprehensive examples and documentation.

## Features

- Automatic example pipeline generation
- Module documentation creation
- GPT-powered analysis
- Customizable number of examples
- Support for context directories
- Integration with knowledge base
- System prompt customization

## Usage

### Basic Setup

```typescript
import { examplifier } from "@felipepimentel/daggerverse/examplifier";

// Initialize the Examplifier module
const client = examplifier({
  gptToken, // OpenAI API token
  systemPromptFile, // Optional: Custom system prompt file
  knowledgeDir, // Optional: Knowledge base directory
});
```

### Generate Examples

```typescript
// Generate examples for a module
const result = await client.examplify({
  address: "github.com/username/module", // Module address
  context: contextDir, // Optional: Context directory
  n: 5, // Optional: Number of examples (default: 5)
});
```

## Configuration

### Constructor Options

The module accepts:

- `token`: OpenAI API token for GPT integration
- `systemPrompt`: Custom system prompt file (default: "system-prompt.txt")
- `knowledgeDir`: Knowledge base directory (default: "./knowledge")

### Examplify Options

The `examplify` method accepts:

- `address`: The address of the Dagger module to analyze
- `context`: Optional directory containing additional context
- `n`: Number of examples to generate (default: 5)

## Examples

### Basic Usage

```typescript
import { examplifier } from "@felipepimentel/daggerverse/examplifier";

export async function generateExamples() {
  // Initialize with default settings
  const client = examplifier({
    gptToken: dag.setSecret("OPENAI_API_KEY", openaiToken),
    // Use default system prompt
    // Use default knowledge directory
  });

  // Generate examples for a module
  const result = await client.examplify({
    address: "github.com/dagger/dagger/sdk",
    // No additional context
    n: 5, // Generate 5 examples
  });
}
```

### Custom Configuration

```typescript
import { examplifier } from "@felipepimentel/daggerverse/examplifier";

export async function generateCustomExamples() {
  // Initialize with custom settings
  const client = examplifier({
    gptToken: dag.setSecret("OPENAI_API_KEY", openaiToken),
    systemPromptFile: dag.file("custom-prompt.txt"),
    knowledgeDir: dag.directory("./custom-knowledge"),
  });

  // Generate examples with context
  const result = await client.examplify({
    address: "github.com/username/module",
    context: dag.directory("./module-context"),
    n: 10, // Generate 10 examples
  });
}
```

### Multiple Modules

```typescript
import { examplifier } from "@felipepimentel/daggerverse/examplifier";

export async function generateMultipleExamples() {
  const client = examplifier({
    gptToken: dag.setSecret("OPENAI_API_KEY", openaiToken),
  });

  // Generate examples for multiple modules
  const modules = [
    "github.com/username/module1",
    "github.com/username/module2",
    "github.com/username/module3",
  ];

  for (const module of modules) {
    const result = await client.examplify({
      address: module,
      n: 5,
    });
  }
}
```

## Output

The module generates a markdown file at `./examplify/README.md` containing:

1. A summary of the module's features and purpose
2. Detailed example pipelines
3. Usage instructions
4. Function and argument documentation

## Dependencies

The module requires:

- Dagger SDK
- OpenAI API token
- GPT module (internal dependency)

## License

See [LICENSE](../LICENSE) file in the root directory.

# GPT Dagger Module

This Dagger module provides integration with OpenAI's GPT models, allowing you to incorporate AI capabilities into your Dagger pipelines. It supports interactive conversations, knowledge base integration, and shell command execution.

## Features

- Integration with OpenAI's GPT models
- Custom system prompts
- Knowledge base management
- Shell command execution
- Conversation history tracking
- Workdir management
- OpenTelemetry integration
- Support for multiple GPT models

## Usage

### Basic Usage

```typescript
import { gpt } from "@felipepimentel/daggerverse/gpt";

// Initialize GPT with API token
const client = gpt({
  token: dag.setSecret("OPENAI_API_KEY", "your-api-key"),
  model: "gpt-4", // model
  knowledgeDir: dag.host().directory("./knowledge"), // knowledge directory
  systemPrompt: dag.host().file("./system-prompt.txt"), // system prompt
});

// Ask a question
const response = await client.ask("What is Dagger?");
```

### Using Knowledge Base

```typescript
// Add knowledge directly
const withKnowledge = client.withKnowledge(
  "dagger-intro",
  "Introduction to Dagger",
  "Dagger is a programmable CI/CD engine that runs your pipelines in containers."
);

// Add knowledge from directory
const withKnowledgeDir = await client.withKnowledgeDir(
  dag.host().directory("./docs")
);
```

### Shell Command Execution

```typescript
// Execute shell commands through GPT
const response = await client.ask("List all files in the current directory");
```

## API Reference

### Constructor Options

The module accepts:

- `token`: OpenAI API token
- `model` (optional): GPT model to use (default: "gpt-4")
- `knowledgeDir` (optional): Directory containing knowledge base files
- `systemPrompt` (optional): File containing system prompt

### Methods

#### `ask(prompt: string)`

Sends a message to the GPT model and processes the response.

Parameters:

- `prompt`: Message to send to the model

#### `withKnowledge(name: string, description: string, contents: string)`

Adds a piece of knowledge to the knowledge base.

Parameters:

- `name`: Unique identifier for the knowledge
- `description`: Short description of the knowledge
- `contents`: Detailed content of the knowledge

#### `withKnowledgeDir(dir: Directory)`

Adds knowledge from text files in a directory.

Parameters:

- `dir`: Directory containing knowledge files

#### `withWorkdir(workdir: Directory)`

Sets the working directory for shell commands.

Parameters:

- `workdir`: Directory to use as working directory

## Knowledge Base Format

Knowledge base files can be either `.txt` or `.md` files. Each file should follow this format:

```
Short description of the knowledge
This is the first paragraph and will be used as the description.

Detailed content of the knowledge.
This can be multiple paragraphs and will be used as the contents.
```

## Example

### Complete Interaction Example

```typescript
import { gpt } from "@felipepimentel/daggerverse/gpt";

export async function example() {
  // Initialize GPT
  const client = gpt({
    token: dag.setSecret("OPENAI_API_KEY", process.env.OPENAI_API_KEY),
    model: "gpt-4",
    knowledgeDir: dag.host().directory("./knowledge"),
    systemPrompt: dag.host().file("./prompts/system.txt"),
  });

  // Add custom knowledge
  const withKnowledge = client.withKnowledge(
    "custom-info",
    "Custom project information",
    "This project uses Dagger for CI/CD pipelines."
  );

  // Set working directory
  const withWorkdir = withKnowledge.withWorkdir(dag.host().directory("."));

  // Have a conversation
  const response1 = await withWorkdir.ask("What is this project about?");

  // Execute a command
  const response2 = await withWorkdir.ask(
    "List all Python files in the project"
  );

  // Print conversation log
  for (const entry of withWorkdir.log) {
    console.log(entry);
  }
}
```

### System Prompt Example

```txt
You are an AI assistant helping with Dagger pipelines.
Your responses should be clear and concise.
When executing commands, prefer using standard Unix tools.
```

## License

See [LICENSE](../LICENSE) file in the root directory.

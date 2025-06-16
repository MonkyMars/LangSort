# File Sorting Tool

A cross-platform Go utility that automatically organizes project directories based on their programming language or framework.

## How it works

1. Scans directories for `.filesort` files
2. Reads the project type from each `.filesort` file
3. Moves the project to an organized directory structure: `~/Coding/{language}/{projectName}`

## Setup

1. Create a `config.json` file:

```json
{
  "sortDir": "~/Coding",
  "acceptedLanguages": ["golang", "nextjs", "python", "javascript"]
}
```

2. In each project directory, create a `.filesort` file:

```
type=golang
```

## Usage

```bash
go build
./filesorting
```

## Example

If you have a project in `/tmp/my-go-project/` with a `.filesort` file containing `type=golang`, the tool will move it to `~/Coding/golang/my-go-project/`.

## Requirements

- Go 1.24.4 or later
- Write permissions to the target directory

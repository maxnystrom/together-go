
# together-go

> **Note**: This library is a proof-of-concept and is not yet feature-complete nor should be considered production ready.

A Go library for interacting with

[Together AI's API](https://docs.together.ai/reference/). This library allows you to:

- Interact with chat and moderation models
- Interact with language, code, and image models
- Embed models
- Fine-tune models

## Installation

You need a working Go environment. Only Go versions according to [Go project's release policy](https://go.dev/doc/devel/release#policy) are currently supported.

```shell
go get github.com/maxnystrom/together-go
```

## Getting Started

```go
package main

import (
   "context"
   "encoding/json"
   "fmt"
   "log"
   "os"
  
   "go get github.com/maxnystrom/together-go"
)

func main() {
  // Construct a new API object using a global API key
  api, err := together.New(os.Getenv("TOGETHER_API_KEY"))
  if err != nil {
    log.Fatal(err)
}

  // Most API calls require a Context
  ctx := context.Background()
  
  // List running instances
  r, err := api.ListRunningInstances(context)
  if err != nil {
    log.Fatal(err)
  }
  
  // Print instance details
  prettyPrint, _ := json.Marshal(r)
  fmt.Println(string(prettyPrint))
}
```

## Contributing

Pull Requests are welcome, but please open an issue (or comment in an existing
issue) to discuss any non-trivial changes before submitting code.

## License

BSD licensed. See the [LICENSE](LICENSE) file for details.

# Spill (WIP)

Spill is a simple and efficient key-value store for Go. It provides a minimalist API for storing and retrieving data.
Why "Spill"? Since it will spill the data to a new map when the total size of the map exceeds the limit.

## Features

- Simple API for data operations
- Comming soon...


## Usage

```go
import "github.com/tituscarl/spill"

// Create a new spill instance
s := spill.New()

// Add data to the spill
s.Put("key", "value")

// Get value
val := s.Get("key")
```


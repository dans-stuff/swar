![adc63e1f-a7e3-4272-b3d4-8e60f88c9b92](https://github.com/user-attachments/assets/3a6af901-d95e-46eb-9c1c-9d395fea8739)

# `swar`: Processs `[]byte` quicker!

**Process 8 bytes at a time using an old technique called Simd Within a Register.**

[![Go Reference](https://pkg.go.dev/badge/github.com/dans-stuff/swar.svg)](https://pkg.go.dev/github.com/dans-stuff/swar) [![Go Report Card](https://goreportcard.com/badge/github.com/dans-stuff/swar)](https://goreportcard.com/report/github.com/dans-stuff/swar) 

- 🚀 **Up to 6x faster** than optimized byte-by-byte code
- 🔌 **Zero dependencies** - no CGO or assembly required
- 🧩 **Dead simple API** - works with your existing code
- ⚡ **Fully portable** - runs anywhere Go runs

```go
chunks, remainder := swar.BytesToLanes(text)
for _, chunk := range chunks {
    // We can work with 8 bytes at a time!
    matches := swar.HighBitWhereEqual(chunk, spaces)
}
```

## Installation

```bash
go get github.com/dans-stuff/swar@latest
```

## Core Operations

| Category | Operations | Use Cases |
|----------|------------|-----------|
| **Comparison** | Equal, Less, Greater | Pattern matching, thresholds |
| **Math** | Add, Subtract, Min/Max, Average | Signal processing, stats |
| **Bit Ops** | Swap nibbles, Reverse bits, Count ones | Encoding, hashing |
| **Selection** | Branchless conditional select | Transformations |

## Real Performance

| Operation | Standard Go | SWAR | Speedup |
|-----------|-------------|------|---------|
| Count character occurrences | 19.30 ns | 7.58 ns | **2.55x** |
| Find uppercase letters | 31.41 ns | 20.32 ns | **1.55x** |
| Convert case | 61.96 ns | 31.53 ns | **1.96x** |
| Detect anomalies | 6.99 ns | 4.17 ns | **1.68x** |

## Full Example: Character Counter

This example counts spaces in a string. For short strings it even outperforms stdlib `bytes.Count`, which is written in assembly! Find more examples in the `examples_test.go` file.

```go
package main

import (
    "fmt"
    "github.com/dans-stuff/swar"
)

func main() {
    text := []byte("Hello, World!")
    
    // Process in 8-byte chunks
    lanes, remainder := swar.BytesToLanes(text)
    
    // Find spaces in parallel
    spaces := swar.Dupe(' ')
    count := 0
    
    for _, lane := range lanes {
        // Sets high bit in bytes equal to space
        matches := swar.HighBitWhereEqual(lane, spaces)
        // Count matches
        count += bits.OnesCount64(matches >> 7)
    }
    
    // Process any leftover bytes
    for _, c := range text[remainder:] {
        if c == ' ' {
            count++
        }
    }
    
    fmt.Printf("Found %d spaces\n", count)
}
```

## Perfect For

- **Text Processing**: UTF-8 validation, parser tokenization
- **Network Protocols**: Header parsing, packet filtering
- **Image Processing**: Thresholding, pixel transformations
- **Data Analysis**: Time series anomaly detection

## How It Works

SWAR treats a 64-bit integer as 8 parallel lanes, using clever bit manipulation to perform the same operation on all bytes simultaneously without branching.

## License & Contributing

MIT Licensed. [Contributions](https://github.com/dans-stuff/swar/fork) welcome!

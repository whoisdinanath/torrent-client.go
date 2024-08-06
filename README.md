# torrentClient

torrentClient is (wip) a Go library for marshalling and unmarshalling BitTorrent bencode data into Go structs using reflection.

## Features

- Encode Go structs to bencode format
- Decode bencode data to Go structs
- Utilizes Go's reflection capabilities

## Installation

To install TorrentClient, use `go get`:

```sh
go get github.com/yourusername/torrentclient
```

## Usage

Here's a basic example of how to use TorrentClient:

```go
package main

import (
	"fmt"
	"github.com/yourusername/torrentclient"
)

func main() {
	// Example struct
	type Torrent struct {
		Announce string
		Info     struct {
			Name   string
			Length int
		}
	}

	// Decode bencode data
	var t Torrent
	err := bencoder.Unmarshal([]byte("d8:announce13:http://...e"), &t)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Decoded:", t)

	// Encode struct to bencode
	data, err := bencoder.Marshal(t)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Encoded:", string(data))
}
```

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss your changes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
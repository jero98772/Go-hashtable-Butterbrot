
# Go-hashtable-Butterbrot

Web Distributed Hash Table testing in Go with design patterns and Redis.

![](https://github.com/jero98772/Go-hashtable-Butterbrot/blob/main/docs/screenshots/2.png)

## Overview

This project implements a simple key-value store API using a combination of a Distributed Hash Table (DHT) for local caching and Redis for persistent storage. It includes REST API endpoints for managing data and serves a static homepage.

## Prerequisites

- **Redis**: Ensure you have Redis installed and running.
  Start Redis with the command:

  ```bash
  redis-server
  ```

- **Go**: Install Go (1.19+ recommended).

## Quick Start

1. Clone the repository:

   ```bash
   git clone https://github.com/jero98772/Go-hashtable-Butterbrot.git
   cd Go-hashtable-Butterbrot
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Start the application:

   ```bash
   go run main.go
   ```

4. Open your browser and visit [http://localhost:8080](http://localhost:8080).

## Project Structure

```
.
├── core/
│   ├── dht.go           # Distributed hash table implementation
│   ├── redis.go         # Redis wrapper for data operations
│   ├── webhandler.go    # API handlers
|   ├── combined.go      # merge dht and redis methots
├── static/
│   └── index.html       # Static homepage
│   └── css/styles.css   # Css styles
│   └── js/scripts.js    # javascript code
├── docs/screenshots/
│   └── 2.png            # Screenshot for the README
├── main.go              # Entry point of the application
├── go.mod               # Go module file
├── go.sum               # Dependency lock file
└── README.md            # Project documentation
```


## License

This project is licensed under the GPLv3 License. See the [LICENSE](LICENSE) file for details.

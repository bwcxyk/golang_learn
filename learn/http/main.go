package main

import (
    "flag"
    "fmt"
    "net/http"
    "sync"
    "time"
)

func main() {
    var (
        url      string
        numConns int
        timeout  int
    )
    flag.StringVar(&url, "url", "", "the url to test")
    flag.IntVar(&numConns, "concurrency", 1, "the number of concurrent connections")
    flag.IntVar(&timeout, "timeout", 10, "the request timeout in seconds")
    flag.Parse()

    if url == "" {
        fmt.Println("Usage: go run main.go --url=http://example.com --concurrency=10 --timeout=5")
        return
    }

    client := http.Client{Timeout: time.Duration(timeout) * time.Second}

    var wg sync.WaitGroup
    wg.Add(numConns)

    start := time.Now()

    for i := 0; i < numConns; i++ {
        go func() {
            defer wg.Done()

            for {
                req, err := http.NewRequest("GET", url, nil)
                if err != nil {
                    fmt.Println("Error creating request:", err)
                    return
                }

                resp, err := client.Do(req)
                if err != nil {
                    fmt.Println("Error sending request:", err)
                    return
                }
                defer resp.Body.Close()
            }
        }()
    }

    wg.Wait()
    elapsed := time.Since(start)

    fmt.Printf("Finished %d requests in %s\n", numConns, elapsed)
    fmt.Printf("Average response time: %s\n", elapsed/time.Duration(numConns))
    fmt.Printf("Requests per second: %.2f\n", float64(numConns)/elapsed.Seconds())
}
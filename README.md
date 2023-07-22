# Concurrent Web Crawler

Concurrent web crawler in Go using goroutines.

## Usage

In your main func

```go
crawler.NewCrawler(&crawler.Config{
  Depth:                    3,
  Breadth:                  0, // Using breadth = 0 to get all links on a page.
  NumWorkers:               50,
  RequestThrottlePerWorker: 100 * time.Millisecond,
}).OutputTo(output_file).ErrorTo(error_file).Crawl([]string{"https://example.com"})

// 50 workers, each sending 10 requests per second, results in 500 requests per second.
```

Output is in CSV format of `<url>,<text>`.

## Design

![design](/diagram.jpg)

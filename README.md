# Concurrent Web Crawler

Concurrent web crawler in Go using goroutines.

## Usage

In your main func

```go
crawler.NewCrawler(&crawler.Config{
  Depth:                    3,
  Breadth:                  0,                      // Using breadth = 0 to get all links on a page.
  NumWorkers:               100,
  RequestThrottlePerWorker: 50 * time.Millisecond,  // Limited to 100 workers * (1000 / 50) = 2000 requests per second.
  SameHostname:             true,                   // Keep only links from the same hostname.
}).OutputTo(output_file).ErrorTo(error_file).Crawl([]string{"https://example.com"})
```

Output is in CSV format of `<url>,<text>,<page_title>`. Can be directly loaded to Pandas. Feel free to use.

## Design

![design](/graphes/design.jpg)

## How it works

![diagram](/graphes/diagram.jpg)

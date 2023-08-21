# Concurrent Web Crawler

Concurrent web crawler in Go using goroutines.

## Usage

In your main func

```go
crawler.NewCrawler(&crawler.Config{
		Depth:            3,
		Breadth:          0,                      // Use breadth = 0 to get all links on a page.
		NumWorkers:       50,                     // Use more workers if having heavy workload per request such as pruning.
		RequestThrottler: 100 * time.Millisecond, // 10 requests per second.
		SameHostname:     true,                   // Keep only links from the same hostname.
	}, os.Stdout, os.Stderr, &crawler.NoOpPruner{}).Crawl([]string{
		"https://example.com",
	})
```

Output is in CSV format of `<depth>,<url>,<text>,<page_title>`. Can be directly loaded to Pandas. Feel free to use.

## Design

![design](/graphes/design.jpg)

## How it works

![diagram](/graphes/diagram.jpg)

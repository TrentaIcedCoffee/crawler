# Concurrent Web Crawler

Concurrent web crawler in Go using goroutines.

## Usage

In your main func

```go
crawler.Crawl([]string{"https://example.com"}, &crawler.Config{
  Depth:      3,
  Breadth:    0, // Breadth == 0 gets all links on a page.
  IsDebug:    true,
  NumWorkers: 10,
})
```

Then

```sh
go run . 2> ./error > ./output
```

Output is in CSV format of `<url>,<text>`.

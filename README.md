Crawl is a simple web crawler written in go. 

```
crawler --help
Usage of ./crawler:
  -assets=false: show page assets in sitemap output
  -concurrency=10: number of concurrent requests
  -insecure=false: ignore invalid site certificates
  -links=false: show page links in sitemap output
  -url="https://example.com": url to crawl
```

```
crawler -links=true -assets=true -concurrency=20 -url=http://example.com
/
 . /stylesheets/screen.css/
 . /js/jquery/jquery-2.1.3.min.js/
 > http://link1.example.com
 > http://link2.example.com
 /foo
     /bar
         /baz

```

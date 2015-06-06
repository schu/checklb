# checklb

## Requirements

```
go get -u github.com/tools/godep
```

## Build

```
make
```

## Usage examples

Send a HTTP request with header 'Host: news.ycombinator.com' to each IP
`news.ycombinator.com` resolves to:

```
./checklb news.ycombinator.com
```

Send a HTTP request with header 'Host: news.ycombinator.com' to each IP
`example.com` resolves to:

```
./checklb news.ycombinator.com example.com
```

Send a HTTP request with header 'Host: news.ycombinator.com' to IP
`198.41.191.47` and `198.41.190.47`:

```
./checklb news.ycombinator.com 198.41.191.47 198.41.190.47
```

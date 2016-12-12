
# tcpstream

tcpstream is simple handling of tcp communication data.<br/>
This is very easy usability.

If you want hook or look tcp raw data.<br/>
tcpstream will help you :D

## Installation

```
go get github.com/narita-takeru/tcpstream
```

## Usage

```
  t := tcpstream.Thread{}

  t.SrcToDstHook = func(b []byte) {
    // Argument "b []byte" is tcp communication data.
    // Do something process.
  }

  t.DstToSrcHook = func(b []byte) {
    // Do something process.
  }

  src := "127.0.0.1:80"
  dst := "127.0.0.1:8080"
  
  t.Do(src, dst) // Start tcp listen by src, and hook both directions.
```


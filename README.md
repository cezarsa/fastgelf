# FastGELF is an *experimental* GELF library

This is a Go library for writing log messages in the GELF format to a remote
endpoint.

## What this library doesn't do:

* No chunking support.
* No compression support.
* No TCP support.
* No support for listening and reading messages.

If you are looking for a feature complete GELF library try
[https://github.com/Graylog2/go-gelf](https://github.com/Graylog2/go-gelf).

## What it does

This library is optimized for writing a large number of UDP GELF messages on
Linux.

It uses the [sendmmsg](http://man7.org/linux/man-pages/man2/sendmmsg.2.html)
syscall, which is only available on Linux, to reduce the number of times we
cross kernel boundaries when sending a huge amount of log messages.

This library will still work on non-linux environments (tested on Darwin) but
the performance improvement won't be so noticeable.

## Benchmarks

These benchmarks compare this library to the Graylog2/go-gelf library.

Linux:

```
$ go test -test.bench . -test.benchmem -test.benchtime 3s
goos: linux
goarch: amd64
pkg: github.com/cezarsa/fastgelf/experiments
Benchmark_FastGELF_WriteMessageExtra/simple         	  500000	      7633 ns/op	     203 B/op	       4 allocs/op
Benchmark_FastGELF_WriteMessageExtra/with_extra     	 1000000	      7535 ns/op	     200 B/op	       4 allocs/op
Benchmark_gogelf_WriteMessageExtra/simple           	  500000	     10941 ns/op	     176 B/op	       1 allocs/op
Benchmark_gogelf_WriteMessageExtra/with_extra       	  300000	     17108 ns/op	     752 B/op	      14 allocs/op
PASS
```

Darwin:

```
$ go test -test.bench . -test.benchmem -test.benchtime 3s
goos: darwin
goarch: amd64
pkg: github.com/cezarsa/fastgelf/experiments
Benchmark_FastGELF_WriteMessageExtra/simple-4         	  300000	     13937 ns/op	     241 B/op	       8 allocs/op
Benchmark_FastGELF_WriteMessageExtra/with_extra-4     	  300000	     13383 ns/op	     242 B/op	       8 allocs/op
Benchmark_gogelf_WriteMessageExtra/simple-4           	  500000	     14276 ns/op	     176 B/op	       1 allocs/op
Benchmark_gogelf_WriteMessageExtra/with_extra-4       	  300000	     15290 ns/op	     752 B/op	      14 allocs/op
PASS
```

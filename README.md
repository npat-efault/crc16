# crc16 [![GoDoc](https://godoc.org/github.com/npat-efault/crc16?status.png)](https://godoc.org/github.com/npat-efault/crc16)
Package crc16 implements the 16-bit Cyclic Redundancy Check

Download:
```shell
go get github.com/npat-efault/crc16
```

* * *

Package crc16 is a Golang implementation of the 16-bit Cyclic
Redundancy Check, or
[CRC-16 checksum](http://en.wikipedia.org/wiki/Cyclic_redundancy_check).

The package's API is almost identical to the standard-library's
[hash/crc32](https://golang.org/pkg/hash/crc32/) and
[hash/crc64](https://golang.org/pkg/hash/crc64/)package.

Package crc16 supports CRC calculation in all possible configurations:
Polynomial, bit-order, byte-order, initial-values, and final value can
all be selected. Predefined configurations are supplied for the most
common uses (e.g. PPP).


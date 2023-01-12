# GROMISE

A go library that immitates the behavior of javascript promises.

## Disclaimer

This is an anti-pattern. Javascript is single threaded and therefore requires the use of promises to manage 
it's single-threaded nature. Go on the other hand is a power language that can run on multiple threads and utilize all cores at the same time.
This library is mainly intended to be a bridge for javascript developers new to the language and/or a way to learn how to utilize go's powerful go-routines.
In some cases, it might prove to be easier on the eye to write code using this library.

## Roadmap

- [x] new
- [x] resolve/reject
- [x] all
- [x] allSettled
- [x] any
- [x] race

## Test

```shell
go test ./... -v
```

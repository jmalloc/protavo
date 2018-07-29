# Protavo

[![Build Status](http://img.shields.io/travis/com/jmalloc/protavo/master.svg)](https://travis-ci.com/jmalloc/protavo)
[![Code Coverage](https://img.shields.io/codecov/c/github/jmalloc/protavo/master.svg)](https://codecov.io/github/jmalloc/protavo)
[![Latest Version](https://img.shields.io/github/tag/jmalloc/protavo.svg?label=semver)](https://semver.org)
[![GoDoc](https://godoc.org/github.com/jmalloc/protavo?status.svg)](https://godoc.org/github.com/jmalloc/protavo/src/protavo)
[![Go Report Card](https://goreportcard.com/badge/github.com/jmalloc/protavo)](https://goreportcard.com/report/github.com/jmalloc/protavo)

Protavo is a simple embedded document store for Go where document content is
represented as protocol buffers messages.

 It includes support for multiple drivers. The reference implementation is built
 on top of [BoltDB](https://github.com/coreos/bbolt).

    go get -u github.com/jmalloc/protavo/src/protavo

> This project is EXPERIMENTAL. Expect frequent breaking changes to the API.

[![GoDoc][1]][2]
[![GoCard][3]][4]

[1]: https://godoc.org/github.com/ditrit/go-linda?status.svg
[2]: https://godoc.org/github.com/ditrit/go-linda
[3]: https://goreportcard.com/badge/ditrit/go-linda
[4]: https://goreportcard.com/report/github.com/ditrit/go-linda


This is a trivial and incomplete implementation of the linda language.

The purpose is to implement the dinner of the philosopher as described in page 451 of the document [Linda in Context](http://www.inf.ed.ac.uk/teaching/courses/ppls/linda.pdf) from Nicholas Carriero and David Gelernter.

# Versions

## v0.2

the v0.2 is using and embedded language based on Lisp (see here [zygomys](https://github.com/glycerine/zygomys)) and [etcd](https://github.com/coreos/etcd) as tuplespace.

See this [blog post](https://blog.owulveryck.info/2017/02/28/to-go-and-touch-lindas-lisp/index.html) for more details.

## v0.1

For more information about the v0.1, please refer to this [blog post](https://blog.owulveryck.info/2017/02/03/linda-31yo-with-5-starving-philosophers.../index.html)

## Running the example

`go get -v github.com/ditrit/go-linda`

`cd $GOPATH/src/github.com/ditrit/go-linda/cmd && go run *.go ../example/dinner/dinner.zy`

# Caution

This project is not production ready at all and has not been tested.
The API may change at each commit.

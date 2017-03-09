[![GoDoc][1]][2]
[![GoCard][3]][4]

[1]: https://godoc.org/github.com/ditrit/go-linda?status.svg
[2]: https://godoc.org/github.com/ditrit/go-linda
[3]: https://goreportcard.com/badge/ditrit/go-linda
[4]: https://goreportcard.com/report/github.com/ditrit/go-linda


This is a trivial and incomplete implementation of the linda language.

The purpose is to implement the dinner of the philosopher as described in page 451 of the document [Linda in Context](http://www.inf.ed.ac.uk/teaching/courses/ppls/linda.pdf) from Nicholas Carriero and David Gelernter.

# Versions

## v0.3

the v0.3 introduces a new primitive: **evalc**.

**evalc** works like **eval** but insted of triggering a new goroutine, it post an event in etcd that is captured by another agent.

More documentation will follow.

![screenshot](https://github.com/ditrit/go-linda/raw/master/doc/v0.3.png)

## v0.2

the v0.2 is using and embedded language based on Lisp (see here [zygomys](https://github.com/glycerine/zygomys)) and [etcd](https://github.com/coreos/etcd) as tuplespace.

See this [blog post](https://blog.owulveryck.info/2017/02/28/to-go-and-touch-lindas-lisp/index.html) for more details.

## v0.1

For more information about the v0.1, please refer to this [blog post](https://blog.owulveryck.info/2017/02/03/linda-31yo-with-5-starving-philosophers.../index.html)

# About

The executables are composed of two parts.

* One agent that is actually the `zygo/linda` interpreter
* One injector that injects code to the tuple space (actually etcd)

## The injector

The injector reads a file and injects its content in etcd. The key is a _uuid_ and is returned if the PUT operation succeed. 

**There is no verification of any sort made on the file before puting it into etcd**

## The agent

The agent takes a _uuid_ as argument.
Then it tries to get the zygo/lisp file from etcd. If it succeeds, it evaluates the content.

## Running the example

`go get -v github.com/ditrit/go-linda`

`cd $GOPATH/src/github.com/ditrit/go-linda/worker && go build`

Make sure an `etcd` daemon is available and accessible.

Then export the following configuration vaiable to reflect your settings. For example:

`export GLINDA_ETCD_ENDPOINT="localhost:2379"`

Then:
 
launch 5 workers (one per philosopher):

`./worker/worker`

and launch the main routine:

`./worker/worker ./example/dinner/dinner.zy`

(or use two separates commands)

# Caution

This project is not production ready at all and has not been tested.
The API may change at each commit.

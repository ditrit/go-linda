[![GoDoc][1]][2]
[![GoCard][3]][4]

[1]: https://godoc.org/github.com/ditrit/go-linda?status.svg
[2]: https://godoc.org/github.com/ditrit/go-linda
[3]: https://goreportcard.com/badge/ditrit/go-linda
[4]: https://goreportcard.com/report/github.com/ditrit/go-linda


This is a trivial and incomplete implementation of the linda language.

The purpose is to implement the dinner of the philosopher as described in page 451 of the document [Linda in Context](http://www.inf.ed.ac.uk/teaching/courses/ppls/linda.pdf) from Nicholas Carriero and David Gelernter.

# Running the example

`go get -v github.com/owulveryck/go-linda`

`cd $GOPATH/src/github.com/owulveryck/go-linda/example/dinner && go run *.go`

<pre>
Philosopher 1 is born
[1] is thinking
Philosopher 0 is born
[0] is thinking
Philosopher 3 is born
[3] is thinking
Philosopher 2 is born
[2] is thinking
Philosopher 4 is born
[4] is thinking
[2] has finished thinking
[2] is hungry
[2] is eating
[1] has finished thinking
[1] is hungry
[1] is eating
[4] has finished thinking
[4] is hungry
[4] is eating
...
</pre>

# Caution

This project is not production ready at all and has not been tested.
The API may change at each commit.

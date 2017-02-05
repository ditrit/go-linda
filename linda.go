package linda

import (
	"reflect"
)

// Linda holds the communication channels that allows to get and put tuples in the Linda
type Linda struct {
	Input  <-chan interface{}
	Output chan<- interface{}
}

// Tuple shoud be a flat structure composed of first class citizens
type Tuple interface{}

// In method extracts a tuple from the tuple space. It finds a tuple that "matches" the object passed as a parameter
// The parameter must be a pointer so the value will be overwritten
// the function blocks until a value is present in the tuple space
// If a template matches a tuple then any formals in the template are
// assigned values from the the tuple.
func (l *Linda) In(m Tuple) {
	for t := range l.Input {
		if match(m, t) {
			// Assign t to m
			m = t
			return
		}
		// Not for me, put the tuple back
		l.Output <- m
	}
}

// Rd is similar to In, except that it does not remove the matched tuple from the tuple space.
// it leaves it unchanged in the tuple space.
func (l *Linda) Rd(t Tuple) {
}

// Out operator inserl a tuple into the tuple space.
func (l *Linda) Out(t Tuple) {
	l.Output <- t
}

// Eval is similar to Out except it launches any function in the struct
// within a goroutine
// TODO: For now eval is only used as a wrapper to launch a goroutine
func (l *Linda) Eval(fns []interface{}) {
	// The first argument of eval should be the function
	if reflect.ValueOf(fns[0]).Kind() == reflect.Func {
		fn := reflect.ValueOf(fns[0])
		var args []reflect.Value
		for i := 1; i < len(fns); i++ {
			args = append(args, reflect.ValueOf(fns[i]))
		}
		go fn.Call(args)
	}
}

// match compares a template m and a tuple t and returns true if:
// 1. m and t have the same number of fields ;
// 2. Corresponding fields have the same types ;
// 3. Each pair of corresponding fields Fm and Ft (in m and t respectively)
//   match. Two fields match only if:
//      - If both Fm and Ft are actuals with "equal" values
//      - TODO: If Fm is a formal and Ft an actual
//      - TODO: If Ft is a formal and Fm an actual
func match(m, t interface{}) bool {
	return reflect.DeepEqual(m, t)
}

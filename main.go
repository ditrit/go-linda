package linda

// Linda holds the communication channels that allows to get and put tuples in the Linda
type Linda struct {
	Input  <-chan interface{}
	Output chan<- interface{}
}

// Tuple shoud be a flat structure composed of first class citizens
type Tuple interface{}

// In  extracl a tuple from the tuple space. It finds a tuple that "matches" the object passed as a parameter
// The parameter must be a pointer so the value will be overwritten
// the function blocks until a value is present in the tuple space
// If a template matches a tuple then any formals in the template are
// assigned values from the the tuple.
func (l *Linda) In(m Tuple) {
	for t := range l.Input {
		if match(m, t) {
			// Assign t to m
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

// Eval is similar to Out except that is evaluates it launches any function in the struct
// within a goroutine
// For now eval only works when reflect.TypeOf(t).Kind() == reflect.Func
func (l *Linda) Eval(t Tuple) {
}

// match compares a template m and a tuple t and returns true if:
// 1. m and t have the same number of fields ;
// 2. Corresponding fields have the same types ;
// 3. Each pair of corresponding fields Fm and Ft (in m and t respectively)
//   match. Two fields match only if:
//      - If both Fm and Ft are actuals with "equal" values
//      - If Fm is a formal and Ft an actual
//      - If Ft is a formal and Fm an actual
func match(m, t interface{}) bool {
	return true
}

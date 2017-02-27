package linda

import (
	"github.com/coreos/etcd/clientv3"
	zygo "github.com/glycerine/zygomys/repl"
	"github.com/google/uuid"
	"log"
)

// New creates a new Linda instance
func New(cli *clientv3.Client) *Linda {
	input := make(chan *zygo.Sexp, 10)
	output := make(chan *zygo.Sexp, 10)
	go func() {
		for i := range output {
			input <- i
		}
	}()
	return &Linda{
		id:     uuid.New(),
		cli:    cli,
		input:  input,
		output: output,
	}
}

// Linda holds the communication channels that allows to get and put tuples in the Linda
type Linda struct {
	id     uuid.UUID
	cli    *clientv3.Client
	input  <-chan *zygo.Sexp
	output chan<- *zygo.Sexp
}

// InRd method extracts a tuple from the tuple space. It finds a tuple that "matches" the object passed as a parameter
// The parameter must be a pointer so the value will be overwritten
// the function blocks until a value is present in the tuple space
// If a template matches a tuple then any formals in the template are
// assigned values from the the tuple.
// The function checks its name.
// if name is "in" tuple is removed, if it is "rd" it does not remove the tuple
func (l *Linda) InRd(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	//log.Println("[InRd]", &args)
	m := zygo.MakeList(args)
	for t := range l.input {
		if match(m, *t) {
			if name == "rd" {
				l.output <- &m
			}
			if len(args) >= 2 {
				log.Printf("[InRd] Matched! %v %v", args[0].SexpString(&zygo.PrintState{}), args[1].SexpString(&zygo.PrintState{}))
			}
			return m, nil
		}
		// Not for me, put the tuple back
		l.output <- &m
	}
	return zygo.SexpNull, nil
}

// Out operator inserl a tuple into the tuple space.
func (l *Linda) Out(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	//log.Println("[Out]", &args)
	lst := zygo.MakeList(args)
	l.output <- &lst
	return zygo.SexpNull, nil
}

// Eval is similar to Out except it launcheS any function in the struct
// within a goroutine
// TODO: For now eval is only used as a wrapper to launch a goroutine
func (l *Linda) Eval(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	//log.Println("[Eval]", args)
	// The first element of the args should be a SexpFunction
	fn := args[0].(*zygo.SexpFunction)
	go func(env *zygo.Glisp, fn *zygo.SexpFunction, args []zygo.Sexp) {
		//_, err := env.Apply(fn, args[:]) // Put the result in the tuplespace
		expr, err := env.Apply(fn, args[:]) // Put the result in the tuplespace
		if err != nil {
			//log.Fatal(err)
		}
		if expr != zygo.SexpNull {
			l.output <- &expr
		}
	}(env.Clone(), fn, args[1:])
	//log.Println("[/Eval]", args)
	return zygo.SexpNull, nil
}

// match compares a template m and a tuple t and returns true if:
// 1. m and t have the same number of fields ;
// 2. Corresponding fields have the same types ;
// 3. Each pair of corresponding fields Fm and Ft (in m and t respectively)
//   match. Two fields match only if:
//      - If both Fm and Ft are actuals with "equal" values
//      - TODO: If Fm is a formal and Ft an actual
//      - TODO: If Ft is a formal and Fm an actual
func match(m, t zygo.Sexp) bool {
	//log.Printf("[match] %v,%v", m.SexpString(&zygo.PrintState{}), t.SexpString(&zygo.PrintState{}))
	return m.SexpString(&zygo.PrintState{}) == t.SexpString(&zygo.PrintState{})
}

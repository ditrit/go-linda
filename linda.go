package linda

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	zygo "github.com/glycerine/zygomys/repl"
	"github.com/google/uuid"
	"log"
)

const prefix = "LINDA"

// New creates a new Linda instance
func New(cli *clientv3.Client) *Linda {
	return &Linda{
		id:  uuid.New(),
		cli: cli,
	}
}

// Linda holds the communication channels that allows to get and put tuples in the Linda
type Linda struct {
	id  uuid.UUID
	cli *clientv3.Client
}

// InRd method extracts a tuple from the tuple space. It finds a tuple that "matches" the object passed as a parameter
// The parameter must be a pointer so the value will be overwritten
// the function blocks until a value is present in the tuple space
// If a template matches a tuple then any formals in the template are
// assigned values from the the tuple.
// The function checks its name.
// if name is "in" tuple is removed, if it is "rd" it does not remove the tuple
func (l *Linda) InRd(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	m := zygo.MakeList(args)
	// Try to see if the tuple already exists in the tuplespace
	resp, err := l.cli.Get(context.TODO(), prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return zygo.SexpNull, err
	}
	for _, ev := range resp.Kvs {
		if match(m, string(ev.Value)) {
			fmt.Printf("%q : %q\n", ev.Key, ev.Value)
			// TODO: Remove the tuple from the space is IN is called
			return zygo.SexpNull, nil
		}
	}
	// The tuple does not exists, watch for a new event until a tuple match
	rch := l.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			// TODO: Check if ev.Type is PUT otherwise continue
			if match(m, string(ev.Kv.Value)) {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				// TODO: Remove the tuple from the space is IN is called
				return zygo.SexpNull, nil
			}
		}
	}
	return zygo.SexpNull, nil
}

// Out operator inserl a tuple into the tuple space.
func (l *Linda) Out(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	lst := zygo.MakeList(args)
	_, err := l.cli.Put(context.TODO(), prefix+"-"+l.id.String()+"-"+uuid.New().String(), lst.SexpString(&zygo.PrintState{}))
	return zygo.SexpNull, err
}

// Eval is similar to Out except it launcheS any function in the struct
// within a goroutine
// TODO: For now eval is only used as a wrapper to launch a goroutine
func (l *Linda) Eval(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	// The first element of the args should be a SexpFunction
	fn := args[0].(*zygo.SexpFunction)
	go func(env *zygo.Glisp, fn *zygo.SexpFunction, args []zygo.Sexp) error {
		//_, err := env.Apply(fn, args[:]) // Put the result in the tuplespace
		expr, err := env.Apply(fn, args[:]) // Put the result in the tuplespace
		if err != nil {
			// TODO
			log.Println(err)
		}
		if expr != zygo.SexpNull {
			_, err := l.cli.Put(context.TODO(), prefix+"-"+l.id.String()+"-"+uuid.New().String(), expr.SexpString(&zygo.PrintState{}))
			return err
		}
		return nil
	}(env.Clone(), fn, args[1:])
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
func match(m zygo.Sexp, t string) bool {
	return m.SexpString(&zygo.PrintState{}) == t
}

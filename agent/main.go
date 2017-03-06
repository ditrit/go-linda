package main

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/ditrit/go-linda"
	zygo "github.com/glycerine/zygomys/repl"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"log"
	"math/rand"
	"os"
	"time"
)

type configuration struct {
	EtcdEndpoints []string      `envconfig:"etcd_endpoint" default:"localhost:2379,localhost:22379,localhost:32379"`
	Timeout       time.Duration `default:"5s"`
}

func sleep(env *zygo.Glisp, name string, args []zygo.Sexp) (zygo.Sexp, error) {
	t := args[0].(*zygo.SexpInt)
	time.Sleep(time.Duration(rand.Int31n(int32(t.Val))) * time.Millisecond)
	return zygo.SexpNull, nil
}

func main() {
	me := uuid.New()
	var s configuration
	err := envconfig.Process("glinda", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Get the Working ID from the cli
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v WORK_ID", os.Args[0])

	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   s.EtcdEndpoints,
		DialTimeout: s.Timeout,
	})
	if err != nil {
		log.Fatalf("Cannot connect to etcd: %v", err)
	}
	defer cli.Close()
	// Getting the lisp source from ETCD
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	resp, err := cli.Get(ctx, os.Args[1])
	cancel()
	if err != nil {
		log.Fatal("Cannot get lisp code from the tuple space (etcd)", err)
	}
	if len(resp.Kvs) != 1 {
		log.Fatalf("Found %v lisp code in the tuple space mtching ID %v", len(resp.Kvs), os.Args[1])
	}
	var lisp []byte
	for _, ev := range resp.Kvs {
		lisp = ev.Value
		//fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

	lda := linda.New(cli)
	env := zygo.NewGlisp()
	env.AddFunction("in", lda.InRd)
	env.AddFunction("rd", lda.InRd)
	env.AddFunction("out", lda.Out)
	env.AddFunction("eval", lda.Eval)
	env.AddFunction("sleep", sleep)
	//env.SourceFile(f)
	env.LoadString(string(lisp))
	// The tuple does not exists, watch for a new event until a tuple match
	rch := cli.Watch(context.Background(), "LINDA-evalc-", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			// Check if the evalc is anonymous or for me...
			if string(ev.Kv.Key) == "LINDA-evalc-" || string(ev.Kv.Key) == "LINDA-evalc-"+me.String() {

				// TODO: Check if ev.Type is PUT otherwise continue
				//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				// TODO: Remove the tuple from the space is IN is called
				_, err := cli.Delete(context.TODO(), string(ev.Kv.Key), clientv3.WithPrefix())
				if err != nil {
					log.Println(err)
					continue
				}
				// Try to decode the message
				log.Println(string(ev.Kv.Value))
				sexp, err := zygo.JsonToSexp(ev.Kv.Value, env)
				_, err = lda.Eval(env, "eval", sexp)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	done := make(chan bool)
	<-done
	//zygo.Repl(env, cfg)
}

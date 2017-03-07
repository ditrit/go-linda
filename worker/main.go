package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/coreos/etcd/clientv3"
	"github.com/ditrit/go-linda"
	zygo "github.com/glycerine/zygomys/repl"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type configuration struct {
	EtcdEndpoints []string      `envconfig:"etcd_endpoint" default:"localhost:2379,localhost:22379,localhost:32379"`
	Timeout       time.Duration `default:"5s"`
	WorkID        string        `envconfig:"work_id" default:"test" required:"true"`
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

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   s.EtcdEndpoints,
		DialTimeout: s.Timeout,
	})
	if err != nil {
		log.Fatalf("Cannot connect to etcd: %v", err)
	}
	defer cli.Close()

	lda := linda.New(cli)
	env := zygo.NewGlisp()
	env.AddFunction("in", lda.InRd)
	env.AddFunction("rd", lda.InRd)
	env.AddFunction("out", lda.Out)
	env.AddFunction("eval", lda.Eval)
	env.AddFunction("evalc", lda.EvalC)
	env.AddFunction("sleep", sleep)

	// We are the main process
	if len(os.Args) > 1 {
		// Put the file in the tuple space
		lisp, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			log.Fatal("Cannot read zygo file: ", err)
		} // Putting the content of the lisp file in the tuple space
		ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
		_, err = cli.Put(ctx, s.WorkID, string(lisp))
		cancel()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		env.SourceFile(f)
		_, err = env.Run()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		// We are an "agent"

		// Getting the lisp source from ETCD
		ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
		resp, err := cli.Get(ctx, s.WorkID)
		cancel()
		if err != nil {
			log.Fatal("Cannot get lisp code from the tuple space (etcd)", err)
		}
		if len(resp.Kvs) != 1 {
			log.Fatalf("Found %v lisp code in the tuple space mtching ID %v", len(resp.Kvs), s.WorkID)
		}
		var lisp []byte
		for _, ev := range resp.Kvs {
			lisp = ev.Value
		}
		//env.SourceFile(f)
		log.Println("Got LISP")
		err = env.LoadString(string(lisp))
		if err != nil {
			log.Fatal("Cannot load lisp", err)
		}
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
					var network = bytes.NewBuffer(ev.Kv.Value)
					dec := gob.NewDecoder(network) // Will read from network.
					var msg []string
					var args []zygo.Sexp
					dec.Decode(&msg)
					for _, m := range msg {
						log.Println(m)
						s, _ := zygo.JsonToSexp([]byte(m), env)
						args = append(args, s)
					}

					log.Println(string(ev.Kv.Value))
					_, err = lda.Eval(env, "eval", args)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
	done := make(chan bool)
	<-done
	//zygo.Repl(env, cfg)
}

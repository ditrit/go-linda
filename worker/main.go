package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/coreos/etcd/clientv3"
	"github.com/ditrit/go-linda"
	zygo "github.com/glycerine/zygomys/repl"
	//"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"log"
	"math/rand"
	"os"
	"regexp"
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
	//me := uuid.New()
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

		// The tuple does not exists, watch for a new event until a tuple match
		rch := cli.Watch(context.Background(), "LINDA-evalc-", clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				// Check if the evalc is anonymous or for me...
				//if string(ev.Kv.Key) == "LINDA-evalc-" || string(ev.Kv.Key) == "LINDA-evalc-"+me.String() {

				// TODO: Check if ev.Type is PUT otherwise continue
				log.Printf("Event: %v: %v / %v", ev.Type.String(), string(ev.Kv.Key), string(ev.Kv.Value))
				if ev.Type.String() != "PUT" {
					break
				}
				//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				// TODO: Remove the tuple from the space is IN is called
				dresp, err := cli.Delete(context.TODO(), string(ev.Kv.Key), clientv3.WithPrefix())
				if err != nil {
					log.Println(err)
					break
				}
				log.Printf("Deleted: %v (%v)", dresp, dresp.Deleted)
				if dresp.Deleted == 0 {
					break
				}
				// Try to decode the message
				var network = bytes.NewBuffer(ev.Kv.Value)
				dec := gob.NewDecoder(network) // Will read from network.
				var msg []string
				dec.Decode(&msg)
				var funct string
				// Ugly hack! get the name of the function to execute
				re := regexp.MustCompile("defn ([^ ]+) ")
				var def string
				for _, m := range msg {
					//log.Println("Got event:", m)
					out := re.FindStringSubmatch(m)
					if len(out) == 2 {
						def = m
						funct = out[1]
						// TODO: fine check the error but by now assume it is because it is a SexpFunc
					} else {
						funct = funct + " " + m
					}
				}
				log.Println(def + "(" + funct + ")")
				env.EvalString(def + "(" + funct + ")")
				//}
			}
		}
	}
	done := make(chan bool)
	<-done
	//zygo.Repl(env, cfg)
}

package shared

import (
	"fmt"
	"io/ioutil"

	"github.com/ChimeraCoder/anaconda"
	"github.com/nats-io/nats"
	"gopkg.in/yaml.v2"
)

type creds map[string]string

func ReadCreds(file *string) creds {
	var cfg creds
	data, _ := ioutil.ReadFile(*file)
	yaml.Unmarshal(data, &cfg)
	return cfg
}

func MakeConn(send_chan chan anaconda.Tweet) (conn *nats.EncodedConn, err error) {
	nc, err := nats.Connect(nats.DefaultURL)

	wch := make(chan bool)
	nc.Opts.AsyncErrorCB = func(c *nats.Conn, s *nats.Subscription, e error) {
		fmt.Print(e)
		wch <- true
	}
	if err != nil {
		fmt.Print(err)
		return
	}

	conn, err = nats.NewEncodedConn(nc, "json")
	if err != nil {
		fmt.Print(err)
		return
	}

	return
}

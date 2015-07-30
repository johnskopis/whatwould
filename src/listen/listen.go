package listen

import (
	"fmt"
	"net/url"
	"shared"

	"github.com/ChimeraCoder/anaconda"
)

func Listen(cfg *string, quit chan bool) {
	creds := shared.ReadCreds(cfg)
	send_chan := make(chan anaconda.Tweet)
	conn, _ := shared.MakeConn(send_chan)
	conn.BindSendChan("tweets", send_chan)

	anaconda.SetConsumerKey(creds["consumer_key"])
	anaconda.SetConsumerSecret(creds["consumer_secret"])
	api := anaconda.NewTwitterApi(creds["access_token"], creds["access_token_secret"])

	v := url.Values{}
	//s := api.UserStream(v)
	s := api.PublicStreamSample(v)
	go func() {
		<-quit
		fmt.Print("Caught an interrupt...cleanign up")
		s.Interrupt()
		s.End()
		conn.Close()
	}()

	for {
		select {
		case o := <-s.C:
			switch twt := o.(type) {
			case anaconda.Tweet:
				send_chan <- twt
				//fmt.Printf("twt: %s\n", twt.Text)
				//spew.Dump(twt.Entities.User_mentions)
			}
		}
	}
}

package speak

import (
	"fmt"
	"io/ioutil"

	"shared"

	"github.com/ChimeraCoder/anaconda"
	"github.com/sajari/fuzzy"
	"gopkg.in/yaml.v2"
)

func read_handles() []string {
	var handles []string
	data, _ := ioutil.ReadFile("handles.yml")
	yaml.Unmarshal(data, &handles)
	return handles
}

func Speak(cfg *string, quit chan bool) {
	handles := read_handles()
	creds := shared.ReadCreds(cfg)
	recv_chan := make(chan anaconda.Tweet)
	conn, err := shared.MakeConn(recv_chan)
	if err != nil {
		fmt.Printf("go error %s", err)
	}

	conn.BindRecvQueueChan("tweets", "worker", recv_chan)

	model := fuzzy.NewModel()
	model.SetDepth(4)
	model.Train(handles)

	anaconda.SetConsumerKey(creds["consumer_key"])
	anaconda.SetConsumerSecret(creds["consumer_secret"])
	//api := anaconda.NewTwitterApi(creds["access_token"], creds["access_token_secret"])

	go func() {
		<-quit
		fmt.Print("Caught an interrupt...cleanign up")
		conn.Close()
	}()

	for {
		select {
		case twt := <-recv_chan:
			if len(twt.Entities.User_mentions) != 0 {
				for _, u := range twt.Entities.User_mentions {
					suggest := model.Suggestions(u.Screen_name, false)
					if len(suggest) != 0 && suggest[0] != u.Screen_name {
						fmt.Printf("suggest: %s\n", model.Suggestions(u.Screen_name, false))
						fmt.Printf("mention: %s\n", u.Screen_name)
						fmt.Printf("twt: %s\n", twt.Text)
						fmt.Printf("%s\n", twt.User.ScreenName)
					}
				}
			}
		}
	}
}

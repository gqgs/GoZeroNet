package site

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/gqgs/go-zeronet/pkg/user"
)

func download(addr string) error {
	pubsubManager := pubsub.NewManager()

	userManager, err := user.NewManager()
	if err != nil {
		return err
	}

	siteManager, err := site.NewManager(pubsubManager, userManager)
	if err != nil {
		return err
	}

	newSite, err := siteManager.NewSite(addr)
	if err != nil {
		return err
	}

	newSite.Announce()

	peers := newSite.Peers()
	log.Println("found ", len(peers), " peers")

	for _, peer := range peers {
		log.Println("connecting to peer", peer)
		if err := peer.Connect(); err != nil {
			log.Println(err)
			continue
		}
		defer peer.Close()

		resp, err := fileserver.GetFile(peer, addr, "content.json", 0, 0)
		if err != nil {
			log.Println(err)
			continue
		}
		jsonDump(resp)
		break
	}

	return nil
}

func jsonDump(v interface{}) {
	d, _ := json.Marshal(v)
	fmt.Println(string(d))
}

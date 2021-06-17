package site

import (
	"log"

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

	newSite.AnnounceTrackers()
	newSite.AnnouncePex()

	peers := newSite.Peers()
	log.Println("found ", len(peers), " peers")

	return newSite.Download()
}

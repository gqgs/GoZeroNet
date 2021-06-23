package site

import (
	"errors"
	"time"

	"github.com/gqgs/go-zeronet/pkg/content"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/peer"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/gqgs/go-zeronet/pkg/user"
)

func download(addr string) error {
	pubsubManager := pubsub.NewManager()

	userManager, err := user.NewManager()
	if err != nil {
		return err
	}

	contentDB, err := database.NewContentDatabase()
	if err != nil {
		return err
	}
	defer contentDB.Close()

	siteManager, err := site.NewManager(pubsubManager, userManager, contentDB)
	if err != nil {
		return err
	}

	newSite, err := siteManager.NewSite(addr)
	if err != nil {
		return err
	}

	contentManager := content.NewManager(contentDB, pubsubManager)
	defer contentManager.Close()

	peerManager := peer.NewManager(pubsubManager, addr)
	defer peerManager.Close()

	go newSite.Announce()

	if err = newSite.Download(peerManager); err != nil {
		return err
	}

	if err := newSite.OpenDB(); err != nil {
		return err
	}
	defer newSite.CloseDB()
	return newSite.RebuildDB()
}

func downloadRecent(addr string) error {
	pubsubManager := pubsub.NewManager()

	userManager, err := user.NewManager()
	if err != nil {
		return err
	}

	contentDB, err := database.NewContentDatabase()
	if err != nil {
		return err
	}
	defer contentDB.Close()

	siteManager, err := site.NewManager(pubsubManager, userManager, contentDB)
	if err != nil {
		return err
	}

	site := siteManager.Site(addr)
	if site == nil {
		return errors.New("site not found")
	}

	contentManager := content.NewManager(contentDB, pubsubManager)
	defer contentManager.Close()

	peerManager := peer.NewManager(pubsubManager, addr)
	defer peerManager.Close()

	go site.Announce()

	if err = site.DownloadSince(peerManager, time.Now().AddDate(0, 0, -7)); err != nil {
		return err
	}

	if err := site.OpenDB(); err != nil {
		return err
	}
	defer site.CloseDB()
	return site.RebuildDB()
}

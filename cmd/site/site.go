package site

import (
	"errors"
	"log"
	"time"

	"github.com/gqgs/go-zeronet/pkg/content"
	"github.com/gqgs/go-zeronet/pkg/database"
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

	contentDB, err := database.NewContentDatabase()
	if err != nil {
		return err
	}
	defer contentDB.Close()

	contentManager := content.NewManager(contentDB, pubsubManager)
	defer contentManager.Close()

	newSite.AnnounceTrackers()
	newSite.AnnouncePex()

	newSite.SetContentDB(contentDB)

	peers := newSite.Peers()
	log.Println("found ", len(peers), " peers")

	return newSite.Download()
}

func downloadRecent(addr string) error {
	pubsubManager := pubsub.NewManager()

	userManager, err := user.NewManager()
	if err != nil {
		return err
	}

	siteManager, err := site.NewManager(pubsubManager, userManager)
	if err != nil {
		return err
	}

	site := siteManager.Site(addr)
	if site == nil {
		return errors.New("site not found")
	}

	contentDB, err := database.NewContentDatabase()
	if err != nil {
		return err
	}
	defer contentDB.Close()

	contentManager := content.NewManager(contentDB, pubsubManager)
	defer contentManager.Close()

	site.AnnounceTrackers()
	site.AnnouncePex()

	site.SetContentDB(contentDB)

	peers := site.Peers()
	log.Println("found ", len(peers), " peers")

	return site.DownloadSince(time.Now().AddDate(0, 0, -7))
}

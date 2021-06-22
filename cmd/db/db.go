package db

import (
	"errors"

	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/gqgs/go-zeronet/pkg/user"
)

func rebuild(addr string) error {
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

	if err := site.OpenDB(); err != nil {
		return err
	}
	defer site.CloseDB()
	return site.RebuildDB()
}

package db

import (
	"errors"

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

	siteManager, err := site.NewManager(pubsubManager, userManager)
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

func query(site, query string) error {
	panic("implement me")
}

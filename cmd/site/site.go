package site

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gqgs/go-zeronet/pkg/content"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/gqgs/go-zeronet/pkg/user"
	"golang.org/x/term"
)

func download(addr string, daysAgo int) error {
	ctx := context.Background()

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

	siteManager, err := site.NewManager(ctx, pubsubManager, userManager, contentDB)
	if err != nil {
		return err
	}
	defer siteManager.Close()

	newSite, err := siteManager.NewSite(addr)
	if err != nil {
		return err
	}

	contentWorker := content.NewWorker(contentDB, pubsubManager)
	defer contentWorker.Close()

	go newSite.Announce()

	now := time.Now()
	if err = newSite.Download(time.Now().AddDate(0, 0, -daysAgo)); err != nil {
		return err
	}

	if err := newSite.OpenDB(); err != nil {
		return err
	}
	defer newSite.CloseDB()
	return newSite.UpdateDB(now)
}

func downloadRecent(addr string, daysAgo int) error {
	ctx := context.Background()

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

	siteManager, err := site.NewManager(ctx, pubsubManager, userManager, contentDB)
	if err != nil {
		return err
	}
	defer siteManager.Close()

	site := siteManager.Site(addr)
	if site == nil {
		return errors.New("site not found")
	}

	contentWorker := content.NewWorker(contentDB, pubsubManager)
	defer contentWorker.Close()

	go site.Announce()
	return site.Update(daysAgo)
}

func verify(addr, innerPath string) error {
	ctx := context.Background()

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

	siteManager, err := site.NewManager(ctx, pubsubManager, userManager, contentDB)
	if err != nil {
		return err
	}
	defer siteManager.Close()

	site := siteManager.Site(addr)
	if site == nil {
		return errors.New("site not found")
	}

	return site.Verify(innerPath)
}

func sign(addr, innerPath string) error {
	ctx := context.Background()

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

	siteManager, err := site.NewManager(ctx, pubsubManager, userManager, contentDB)
	if err != nil {
		return err
	}
	defer siteManager.Close()

	site := siteManager.Site(addr)
	if site == nil {
		return errors.New("site not found")
	}

	fmt.Fprintln(os.Stdout, "Private key:")
	privateKey, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	return site.Sign(innerPath, string(privateKey))
}

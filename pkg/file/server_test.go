package file

import (
	"fmt"

	"github.com/gqgs/go-zeronet/pkg/config"
)

func testURL() string {
	return fmt.Sprintf("http://localhost:%d/", config.FileServer.Port)
}

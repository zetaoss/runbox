package runid

import (
	"strings"
	"time"

	"github.com/zetaoss/runbox/pkg/util"
)

func New(parts ...string) string {
	t := time.Now().Format("20060102150405")
	return t[:8] + "-" + t[8:] + "-" + util.NewHash(5) + "-" + strings.Join(parts, "-")
}

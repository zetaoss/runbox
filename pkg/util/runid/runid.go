package runid

import (
	"strings"
	"time"

	"github.com/zetaoss/runbox/pkg/util"
)

func New(parts ...string) string {
	return time.Now().Format("20060102150405") + "_" + util.NewHash(3) + "_" + strings.Join(parts, "_")
}

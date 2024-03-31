package single

import (
	"github.com/zetaoss/zetarun/pkg/run/lang/multi"
	"github.com/zetaoss/zetarun/pkg/run/lang/types"
)

func Run(input types.SingleInput, extraOpts ...map[string]int) (*types.Output, error) {
	return multi.Run(types.MultiInput{Lang: input.Lang, Files: []types.File{{Content: input.Source}}}, extraOpts...)
}

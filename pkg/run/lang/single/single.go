package single

import (
	"github.com/zetaoss/runbox/pkg/run/lang/multi"
	"github.com/zetaoss/runbox/pkg/run/lang/types"
	"github.com/zetaoss/runbox/pkg/util/runid"
)

func Run(input types.SingleInput, extraOpts ...map[string]int) (*types.Output, error) {
	if input.Source == "" {
		return nil, types.ErrNoSource
	}
	if input.RunID == "" {
		input.RunID = runid.New("single", input.Lang)
	}
	return multi.Run(types.MultiInput{Lang: input.Lang, Files: []types.File{{Content: input.Source}}}, extraOpts...)
}

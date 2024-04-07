package single

import (
	"github.com/zetaoss/runbox/pkg/runner/lang/multi"
	"github.com/zetaoss/runbox/pkg/runner/lang/types"
	"github.com/zetaoss/runbox/pkg/util/runid"
	"k8s.io/klog/v2"
)

func Run(input Input, extraOpts ...map[string]int) (*types.Output, error) {
	if input.Source == "" {
		return nil, ErrNoSource
	}
	if input.RunID == "" {
		input.RunID = runid.New("single", input.Lang)
	}
	output, err := multi.Run(multi.Input{Lang: input.Lang, Files: []multi.File{{Text: input.Source}}}, extraOpts...)
	if err != nil {
		if err != multi.ErrInvalidLanguage {
			klog.Warningf("unknown err: %s", err.Error())
		}
		return nil, ErrInvalidLanguage
	}
	return output, nil
}

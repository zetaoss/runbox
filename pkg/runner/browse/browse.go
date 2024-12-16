package browse

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/zetaoss/runbox/pkg/runner/box"
	"golang.org/x/net/html"
	"k8s.io/utils/ptr"
)

type Browse struct {
	box *box.Box
}

func New(box *box.Box) *Browse {
	return &Browse{box}
}

func textFromNode(n *html.Node) string {
	fmt.Println("n", n)
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(textFromNode(c))
	}
	return buf.String()
}

func (b *Browse) Run(urlString string) (string, error) {
	image := "selenium/standalone-chrome:3.141.59"
	command := fmt.Sprintf(`/opt/google/chrome/chrome --headless --dump-dom --disable-gpu --no-sandbox '%s'`, urlString)
	opts := &box.Opts{
		CollectStats: ptr.To(false),
		Command:      command,
		Image:        image,
		Timeout:      30000,
	}
	result, err := b.box.Run(opts)
	if err != nil {
		return "", fmt.Errorf("box.Run err: %w", err)
	}
	if result.Timedout {
		return "", errors.New("timed out")
	}
	stdout, _ := result.StreamStrings()
	return stdout, nil
}

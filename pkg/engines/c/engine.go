package C

import (
	"fmt"
	"strings"

	"github.com/kampanosg/lazytest/pkg/engines"
	"github.com/kampanosg/lazytest/pkg/models"
)

const icon = "C"

type cNode struct {
	Ref      any
	Name     string
	Children map[string]*cNode
}

type CEngine struct {
	Runner engines.Runner
}

func NewCEngine(r engines.Runner) *CEngine {
	return &CEngine{
		Runner: r,
	}
}

func (c *CEngine) GetIcon() string { return icon }

func (c *CEngine) Load(dir string) (*models.LazyTree, error) {
	// In C the command to run tests is un standardised
	// since its often tied to the makefile. should add a config file
	// so people can tweek this. At the moment they just have to setup
	// the makefile in a compatible way
	o, err := c.Runner.RunCmd("make test TFLAGS=--list")
	if err != nil {
		return nil, nil
	}
	root := &cNode{
		Name:     dir,
		Ref:      nil,
		Children: make(map[string]*cNode),
	}
	lines := strings.Split(o, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "/")
		// segfault if len parts == 1
		if len(parts) == 0 {
			continue
		}

		currentNode := root
		var testSuite *models.LazyTestSuite

		for i, part := range parts {
			childNode, exists := currentNode.Children[part]

			if !exists {
				childNode = &cNode{Name: part, Children: make(map[string]*cNode)}
				if len(parts)-2 == i {
					childNode.Ref = &models.LazyTestSuite{
						Path:  strings.Join(parts[:i+1], "::"),
						Tests: make([]*models.LazyTest, 0),
					}
				}
				currentNode.Children[part] = childNode
			}
			currentNode = childNode

			if i == len(parts)-2 {
				testSuite = currentNode.Ref.(*models.LazyTestSuite)
			}

			if i == len(parts)-1 {
				test := &models.LazyTest{
					Name:   part,
					RunCmd: fmt.Sprintf("make test TFLAGS=%s", line),
				}
				childNode.Ref = test
				testSuite.Tests = append(testSuite.Tests, test)
			}
		}
	}

	if len(root.Children) == 0 {
		return nil, nil
	}

	lazyRoot := toLazyTree(root)
	return models.NewLazyTree(lazyRoot), nil
}

func toLazyTree(r *cNode) *models.LazyNode {
	children := make([]*models.LazyNode, 0)
	for _, child := range r.Children {
		children = append(children, toLazyTree(child))
	}

	return &models.LazyNode{
		Name:     r.Name,
		Ref:      r.Ref,
		Children: children,
	}
}

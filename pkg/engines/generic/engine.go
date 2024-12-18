package generic

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kampanosg/lazytest/pkg/config"
	"github.com/kampanosg/lazytest/pkg/engines"
	"github.com/kampanosg/lazytest/pkg/models"
)

type genNode struct {
	Ref      any
	Children map[string]*genNode
	Name     string
}

// TODO add a way to switch between list command/file naming convention
type GenEngine struct {
	Runner engines.Runner
	config conf.EngineConfig
}

func NewGenEngine(
	r engines.Runner,
	config conf.EngineConfig,
) *GenEngine {
	return &GenEngine{
		Runner: r,
		config: config,
	}
}

func (g *GenEngine) GetIcon() string { return g.config.Icon }

func (g *GenEngine) Load(dir string) (*models.LazyTree, error) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	o, err := g.Runner.RunCmd(fmt.Sprintf("cd %s && %s", dir, g.config.ListCommand))
	if err != nil {
		return nil, nil
	}
	root := &genNode{
		Name:     dir,
		Ref:      nil,
		Children: make(map[string]*genNode),
	}
	lines := strings.Split(o, "\n")[g.config.SkipLines:]
	for _, oline := range lines {
		oline = strings.TrimSpace(oline)
		if oline == "" {
			continue
		}
		line := strings.Replace(oline, g.config.TestSeperator, g.config.DirSeperator, 1)

		parts := strings.Split(line, g.config.DirSeperator)
		if len(parts) == 0 {
			continue
		}

		currentNode := root
		var testSuite *models.LazyTestSuite

		for i, part := range parts {
			if part == "" {
				continue
			}
			childNode, exists := currentNode.Children[part]

			if !exists {
				childNode = &genNode{Name: part, Children: make(map[string]*genNode)}
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
					RunCmd: fmt.Sprintf("cd %s && %s%s", dir, g.config.RunCommand, oline),
				}
				childNode.Ref = test
				log.Println(test)
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

func toLazyTree(r *genNode) *models.LazyNode {
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

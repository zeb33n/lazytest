package main

import (
	"flag"
	"fmt"
	"slices"
	"strings"

	"github.com/kampanosg/lazytest/internal/clipboard"
	"github.com/kampanosg/lazytest/internal/runner"
	"github.com/kampanosg/lazytest/internal/tui"
	"github.com/kampanosg/lazytest/internal/tui/elements"
	"github.com/kampanosg/lazytest/internal/tui/handlers"
	"github.com/kampanosg/lazytest/internal/tui/state"
	"github.com/kampanosg/lazytest/pkg/engines"
	"github.com/kampanosg/lazytest/pkg/engines/golang"
	"github.com/rivo/tview"
)

const (
	Version = "v.0.2.0"
)

func main() {
	dir := flag.String("dir", ".", "the directory to start searching for tests")
	exc := flag.String("excl", "", "engines to exclude")
	vsn := flag.Bool("version", false, "the current version of LazyTest")
	flag.Parse()

	if *vsn {
		fmt.Printf("LazyTest %s\n", Version)
		return
	}

	excludedEngines := strings.Split(*exc, ",")
	var engines []engines.LazyEngine

	if !slices.Contains(excludedEngines, "golang") {
		engines = append(engines, golang.NewGolangEngine())
	}

	a := tview.NewApplication()
	h := handlers.NewHandlers()
	r := runner.NewRunner()
	e := elements.NewElements()
	c := clipboard.NewClipboardManager()
	s := state.NewState()

	t := tui.NewTUI(a, h, r, c, e, s, *dir, engines)

	if err := t.Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/Larufu1880/git-remove-branch/cmd"
	"github.com/jroimartin/gocui"
)

var gr = &cmd.GitRepository{}

func main() {
	gr.SetBranches()
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {
	_, maxY := g.Size()

	branchNames := gr.GetBranchesName()

	if v, err := g.SetView("select", 0, 0, 4, maxY*9/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		for range branchNames {
			fmt.Fprintln(v, "[]")
		}
	}
	if v, err := g.SetView("main", 4, 0, 100, maxY*9/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		for _, name := range branchNames {
			fmt.Fprintln(v, name)
		}
		v.Highlight = true
		v.Editable = true
		v.Wrap = true
	}
	if v, err := g.SetView("info", 0, maxY*9/10, 100, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(v, "Select a branch name that you wanna delete by space key. Undo by space key. Delete by Enter")
	}
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeySpace, gocui.ModNone, selectBranch); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, deleteBranch); err != nil {
		return err
	}
	return nil
}

func deleteBranch(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		if err := gr.DeleteBranches(); err != nil {
			log.Fatalln(err)
		}
		g.Update(initView)
	}
	return nil
}

func selectBranch(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		gr.SetDeleteFlg(cy)
		g.Update(updateSelectView)
	}
	return nil
}

func initView(g *gocui.Gui) error {
	gr.SetBranches()
	v, err := g.View("main")
	if err != nil {
		log.Panic(err)
	}
	v.Clear()
	for _, branchName := range gr.GetBranchesName() {
		fmt.Fprintln(v, branchName)
	}
	v, err = g.View("select")
	v.Clear()
	if err != nil {
		log.Panic(err)
	}
	for _, flg := range gr.GetDeleteFlg() {
		if flg {
			fmt.Fprintln(v, "[x]")
		} else {
			fmt.Fprintln(v, "[]")
		}
	}
	return nil
}

func updateSelectView(g *gocui.Gui) error {
	v, err := g.View("select")
	if err != nil {
		log.Panicln(err)
	}
	v.Clear()
	for _, flg := range gr.GetDeleteFlg() {
		if flg {
			fmt.Fprintln(v, "[x]")
		} else {
			fmt.Fprintln(v, "[]")
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

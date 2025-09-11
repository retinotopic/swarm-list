package main

import (
	"github.com/GianlucaP106/gotmux/gotmux"
	"log"
)

type PaneInfo struct {
	isvertical bool
	P          *gotmux.Pane
}

func main() {
	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		log.Fatal(err,"1")
	}
	// sess ,err := tmux.ListSessions()
	// for i := range sess {
	// 	if sess[i] != nil {
	// 		sess[i].Kill()
	// 	}
	// }
	session, err := tmux.New()
	if err != nil {
		log.Fatal(err,"2")
	}
	window, err := session.New()
	if err != nil {
		log.Fatal(err,"3")
	}
	var direction gotmux.PaneSplitDirection
	for range 5 {
		panes,err := window.ListPanes()
		if err != nil {
			log.Fatal(err,"4")
		}
		for j := range panes {
			panes[j].SplitWindow(&gotmux.SplitWindowOptions{SplitDirection: direction,
			ShellCommand: "go run ../sim.go",})
		}
		switch direction {
			case gotmux.PaneSplitDirectionHorizontal:
				direction = gotmux.PaneSplitDirectionVertical
			default:
				direction = gotmux.PaneSplitDirectionHorizontal
		}
	}
	session.Attach()

}

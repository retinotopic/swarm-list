package main

import (
	"bufio"

	"strconv"

	"log"

	"github.com/gdamore/tcell/v2"

	crand "crypto/rand"

	"github.com/retinotopic/GoChat/app"

	"github.com/retinotopic/GoChat/app/list"
	mrand "math/rand/v2"

	"math/rand/v2"
	"os"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

var lettersrune = []rune(letters)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var buf [32]byte
	crand.Read(buf[:])
	cha := mrand.NewChaCha8(buf)
	rand := mrand.New(cha)

	sm := Sim{rnd: rand, tcevs: make([]*tcell.EventKey, 10),
		Lines: [4][2]int{
			{2, 34}, {37, 53}, {56, 66}, {69, 83},
		}}
	sm.CreateLineIndex(dir + "/sim")

	usr := sm.rnd.IntN(100)
	usrstr := strconv.Itoa(usr)
	apphost := os.Getenv("APP_HOST")
	if len(apphost) == 0 {
		apphost = "localhost"
	}
	appport := os.Getenv("APP_PORT")
	if len(appport) == 0 {
		appport = "8080"
	}
	wsUrl := "ws" + "://" + "localhost" + ":" + "80" + "/connect"
	dflog := log.New(os.Stdout, "app log: ", 0)
	sm.Chat = app.NewChat(usrstr, wsUrl, 20, true, true, dflog, dflog)
	errch := sm.Chat.TryConnect()
	sm.Errch = errch
	for {
		n := sm.rnd.IntN(100)
		x := 0
		switch {
		case n < 50: // 50% message
			x = 2
		case n < 75: // 25% (50-74) add/delete to/from room
			x = 3
		case n < 87: // 12% (75-86) find user and create duo room
			x = 0
		default: // 13% (87-99) create group room
			x = 1
		}
		sm.ReadLinesByIndex(sm.LinesIndex, sm.Lines[x])
	}

}

type Sim struct {
	rnd        *rand.Rand
	Chat       *app.Chat
	file       *os.File
	Lines      [4][2]int
	tcevs      []*tcell.EventKey
	LinesIndex []int64
	Errch      <-chan error
}

func (s *Sim) CreateLineIndex(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	s.file = file
	defer file.Close()

	scanner := bufio.NewScanner(file)
	offsets := []int64{0}

	for scanner.Scan() {
		currentOffset := offsets[len(offsets)-1] + int64(len(scanner.Bytes())) + 1
		offsets = append(offsets, currentOffset)
	}
	s.LinesIndex = offsets
}

func (s *Sim) ReadLinesByIndex(offsets []int64, lines [2]int) {
	fromline := lines[0]
	toline := lines[1]
	s.file.Seek(offsets[fromline], 0)

	scanner := bufio.NewScanner(s.file)
	linesRead := 0
	totalLines := toline - fromline
	s.tcevs = s.tcevs[:0]

	for linesRead < totalLines && scanner.Scan() {
		ev := scanner.Text()
		runeev := []rune(ev)
		typekey := string(runeev[3:6])
		switch string(typekey) {
		case "Key", "Run", "%%%":
			if string(typekey) == "%%%" {
				cn := s.rnd.IntN(2)
				if cn == 0 {
					continue
				}
			}
			var tcev *tcell.EventKey
			numb := runeev[8:]
			n, err := strconv.Atoi(string(numb))
			if err != nil {
				return
			}
			if string(typekey) == "Key" {
				tcev = tcell.NewEventKey(tcell.Key(n), ' ', tcell.ModNone)
			} else {
				tcev = tcell.NewEventKey(tcell.KeyRune, rune(n), tcell.ModNone)
			}
			s.tcevs = append(s.tcevs, tcev)
		case "###":
			length := s.rnd.IntN(15)
			s.tcevs = append(s.tcevs, tcell.NewEventKey(tcell.Key(13), ' ', tcell.ModNone))

			for range length {
				idx := s.rnd.IntN(len(letters))
				rune := lettersrune[idx]

				s.tcevs = append(s.tcevs, tcell.NewEventKey(tcell.KeyRune, rune, tcell.ModNone))
			}
		case "$$$":
			times := s.rnd.IntN(15)
			for range times {
				s.tcevs = append(s.tcevs, tcell.NewEventKey(tcell.Key(13), ' ', tcell.ModNone))
				s.tcevs = append(s.tcevs, tcell.NewEventKey(tcell.Key(258), ' ', tcell.ModNone))
			}
		case "@@@":
			runeusr := []rune{'u', 's', 'e', 'r'}
			s.tcevs = append(s.tcevs, tcell.NewEventKey(tcell.Key(13), ' ', tcell.ModNone))
			for i := range len(runeusr) {
				s.tcevs = append(s.tcevs, tcell.NewEventKey(tcell.KeyRune, runeusr[i], tcell.ModNone))
			}
			usrid := s.rnd.IntN(100)
			s.tcevs = append(s.tcevs, tcell.NewEventKey(tcell.KeyRune, rune(usrid), tcell.ModNone))
		}
		linesRead++
	}
	if len(s.tcevs) != 0 {
		for i := range s.tcevs {
			ky := s.tcevs[i]
			key := s.Chat.MainFlexNavigation(ky)
			prm := s.Chat.MainFlex.GetItem(s.Chat.NavState)
			l, ok := prm.(*list.List)
			if !ok {
				panic("Not list")
			}
			l.InputHandlerRaw(key, nil)
		}
	}
}

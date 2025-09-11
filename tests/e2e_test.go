package e2e_test

import (
	"bufio"
	"context"
	"slices"

	"strconv"
	"time"

	"log"
	"testing"

	"github.com/gdamore/tcell/v2"

	// "log"
	"os"

	"github.com/retinotopic/GoChat/app"
	"github.com/retinotopic/GoChat/app/list"
)

type ChatInfo struct {
	ChatClient *app.Chat
	User       string
	Errch      <-chan error
}

func Test_e2e(t *testing.T) {
	users := []string{"u1", "u2", "u3"}
	usertochat := map[string]ChatInfo{}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.Open(dir + "/testlogs")
	if err != nil {
		panic(err)
	}
	logsoftest, err := os.OpenFile(dir+"/logsoftest", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	lg := log.New(logsoftest, ":", 0)
	wsstr := "ws"
	if os.Getenv("SSL_ENABLE") == "true" {
		wsstr = "wss"
	}
	apphost := os.Getenv("APP_HOST")
	if len(apphost) == 0 {
		apphost = "localhost"
	}
	appport := os.Getenv("APP_PORT")
	if len(appport) == 0 {
		appport = "8080"
	}
	wsUrl := wsstr + "://" + "localhost" + ":" + "80" + "/connect"
	dflog := log.New(os.Stdout, "app log: ", 0)
	for i := range users {
		chat := app.NewChat(users[i], wsUrl, 20, true, true, dflog, dflog)
		errch := chat.TryConnect()
		go chat.ProcessEvents()
		for range 2 {
			chat.TestCh <- struct{}{} // skipping initial events
		}
		usertochat[users[i]] = ChatInfo{User: users[i], ChatClient: chat, Errch: errch}
	}
	scan := bufio.NewScanner(f)
	currChatInfo := usertochat["u1"]
	var tcev *tcell.EventKey
	for scan.Scan() {
		select {
		case _, ok := <-usertochat["u1"].Errch:
			if !ok {
				t.Fatal("test failed")
			}
		case _, ok := <-usertochat["u2"].Errch:
			if !ok {
				t.Fatal("test failed")
			}
		case _, ok := <-usertochat["u3"].Errch:
			if !ok {
				t.Fatal("test failed")
			}
		default:
			ev := scan.Text()
			runeev := []rune(ev)
			user := string(runeev[:2])
			if currChatInfo.User != user {
				v, ok := usertochat[user]
				if !ok {
					continue
				}
				currChatInfo = v
			}
			typekey := string(runeev[3:6])
			switch string(typekey) {
			case "Sgn":
				expected := string(runeev[8:])
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				WaitForEvent(ctx, cancel, t, currChatInfo.ChatClient, expected, lg, usertochat)
			case "Key", "Run":
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
			}
			if tcev != nil {
				// tcell.NewEventKey(tcell.KeyRune,)
				currChatInfo.ChatClient.Mtx.Lock()
				key := currChatInfo.ChatClient.MainFlexNavigation(tcev)
				prm := currChatInfo.ChatClient.MainFlex.GetItem(currChatInfo.ChatClient.NavState)
				l, ok := prm.(*list.List)
				if !ok {
					t.Fatal(" Not list (somehow) ")
				}
				l.InputHandlerRaw(key, nil)
				currChatInfo.ChatClient.Mtx.Unlock()
			}
			tcev = nil
		}
	}
	if scan.Err() != nil {
		t.Fatal(scan.Err())
	}
}

func WaitForEvent(ctx context.Context, cancel context.CancelFunc, t *testing.T, app *app.Chat,
	expected string, lg *log.Logger, usertochat map[string]ChatInfo) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			// lg.Println(usertochat["u1"].ChatClient.Checksums)
			// lg.Println(usertochat["u2"].ChatClient.Checksums)
			// lg.Println(usertochat["u3"].ChatClient.Checksums)
			t.Fatal("context timeout")
		case app.TestCh <- struct{}{}:
			if !slices.Contains(app.Checksums, expected) {
				lg.Println(usertochat["u1"].ChatClient.Checksums)
				lg.Println(usertochat["u2"].ChatClient.Checksums)
				lg.Println(usertochat["u3"].ChatClient.Checksums)
				t.Fatal("unexpected checksum")
			}
			return
		}
	}
}

package app

import (
	"reflect"

	"github.com/retinotopic/GoChat/app/list"
)

type EventExecer interface {
	ExecEvent()
}

func Lst(cnt ...string) []list.Content {
	listcnt := make([]list.Content, 0)
	txt := "SampleText"
	mstxt := txt
	for i := range cnt {
		switch cnt[i] {
		case "true":
			listcnt = append(listcnt, list.Content{MainText: mstxt})
			mstxt = txt
		case "false":
			listcnt = append(listcnt, list.Content{SecondaryText: mstxt})
			mstxt = txt
		default:
			mstxt = cnt[i]
		}
	}
	return listcnt
}

func (c *Chat) UI(cnt []list.Content, lists ...int) UIEvent {
	return UIEvent{
		C:         c,
		Content:   cnt,
		ShowLists: lists,
	}
}

func (c *Chat) Send(prep any, funcname string, eventname string, tp int, trg ...int) SendEvent {
	switch prep.(type) {
	case *Room, *Message, *User:
		break
	default:
		panic("This type is not implemented")
	}
	raw := reflect.ValueOf(prep).MethodByName(funcname).Interface()
	fn := raw.(func(*EventInfo, []list.Content) error)
	return SendEvent{
		InitEvent: EventInfo{Type: tp, Event: eventname},
		ExecFn:    fn,
		C:         c,
		TargetList: func() int {
			if len(trg) > 0 {
				return trg[0]
			}
			return -1
		}(),
	}
}
func Key(maintxt string, sectext string) list.Content {
	return list.Content{
		MainText:      maintxt,
		SecondaryText: sectext,
	}
}

func (c *Chat) InitEvents() {
	room := &Room{}
	message := &Message{C: c}
	user := &User{}
	c.EventMap = map[list.Content]EventExecer{
		// UI events
		Key("Events", ""): c.UI(Lst(), 4),

		Key("Menu", ""): c.UI(Lst("Create Duo Room", "true", "Create Group Room", "true",
			"Unblock Users", "true", "Change Username", "true", "Change Privacy", "true",
			"Find Users", "true", "Block Users", "true"), 2),

		Key("This Group Room(Admin)", ""): c.UI(Lst("Delete Users From Room", "true",
			"Add Users To Room", "true", "Change Room Name", "true", "Show Users", "true"), 2),

		Key("This Group Room", ""):        c.UI(Lst("Show Users", "true", "Delete Users From Room", "true"), 2),
		Key("This Duo Room", ""):          c.UI(Lst("Show Users", "true"), 2),
		Key("Create Duo Room", ""):        c.UI(Lst("Create Duo Room", "false"), 2, 5),
		Key("Create Group Room", ""):      c.UI(Lst("Create Group Room", "false"), 2, 6, 3),
		Key("Unblock Users", ""):          c.UI(Lst("Unblock User", "false", "Get Blocked Users", "false"), 2, 7),
		Key("Block Users", ""):            c.UI(Lst("Block User", "false"), 2, 6),
		Key("Add Users To Room", ""):      c.UI(Lst("Add Users", "false"), 2, 6),
		Key("Delete Users From Room", ""): c.UI(Lst("Delete Users", "false"), 2, 8),
		Key("Change Room Name", ""):       c.UI(Lst("Change Room Name", "false"), 2, 3),
		Key("Change Username", ""):        c.UI(Lst("Change Username", "false"), 2, 3),
		Key("Find Users", ""):             c.UI(Lst("Find Users", "false"), 2, 5, 3),

		Key("Change Privacy", ""): c.UI(Lst("Change Duo Room Policy", "false",
			"Change Group Room Policy", "false"), 2, 9),

		Key("Show Users", ""): c.UI(Lst(), 8),
		// Send events
		Key("", "Send Message"):             c.Send(message, "SendMessage", "Send Message", 1),
		Key("", "Add Users"):                c.Send(room, "AddDeleteUsersInRoom", "Add Users To Room", 2, 6),
		Key("", "Delete Users"):             c.Send(room, "AddDeleteUsersInRoom", "Delete Users From Room", 2, 8),
		Key("", "Get Blocked Users"):        c.Send(user, "GetBlockedUsers", "Get Blocked Users", 3),
		Key("", "Unblock User"):             c.Send(room, "BlockUnblockUser", "Unblock User", 2, 7),
		Key("", "Block User"):               c.Send(room, "BlockUnblockUser", "Block User", 2, 6),
		Key("", "Change Duo Room Policy"):   c.Send(user, "ChangePrivacy", "Change Privacy Direct", 3, 9),
		Key("", "Change Group Room Policy"): c.Send(user, "ChangePrivacy", "Change Privacy Group", 3, 9),
		Key("", "Change Username"):          c.Send(user, "ChangeUsernameFindUsers", "Change Username", 3),
		Key("", "Find Users"):               c.Send(user, "ChangeUsernameFindUsers", "Find Users", 3),
		Key("", "Get Messages From Room"):   c.Send(message, "GetMessagesFromRoom", "Get Messages From Room", 1),
		Key("", "Change Room Name"):         c.Send(room, "ChangeRoomName", "Change RoomName", 2),
		Key("", "Create Duo Room"):          c.Send(room, "CreateDuoRoom", "Create Duo Room", 2, 5),
		Key("", "Create Group Room"):        c.Send(room, "CreateGroupRoom", "Create Group Room", 2, 6),
	}
	c.SendEventCh = make(chan EventInfo, 100)
	c.TestCh = make(chan struct{})
	c.Checksums = []string{}
	c.RoomMsgs = make(map[uint64]*RoomInfo)
	c.DuoUsers = make(map[uint64]User)
	c.UserBuf = make([]User, 10)
}

var DebugKey = "2ooP5g11hBa62nV06sFjf4j24M1c9vRDY" // checksums for testing

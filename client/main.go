package main

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/marcusolsson/tui-go"
)

type client struct {
	conn net.Conn
}

func (c *client) connect() error {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *client) receive(box *tui.Box) {
	for {
		b := make([]byte, 1024)
		_, err := c.conn.Read(b)
		if err == io.EOF {
			// サーバが終了
			log.Fatal("server down")
		} else if err != nil {
			log.Print(err)
			continue
		}

		box.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(string(b))),
			tui.NewSpacer(),
		))
		ui.Repaint()
	}
}

func (c *client) send(msg string) error {
	_, err := c.conn.Write([]byte(msg))
	return err
}

var ui tui.UI

func main() {
	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	c := new(client)

	// サーバに接続
	err := c.connect()
	if err != nil {
		log.Fatal("failed to connect server", err)
	}

	// サーバからのメッセージ受信待ち
	go c.receive(history)

	// サーバへメッセージ送信
	input.OnSubmit(func(e *tui.Entry) {
		text := e.Text()
		err := c.send(text)
		if err != nil {
			log.Print("send error", err)
		}
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(text)),
			tui.NewSpacer()),
		)
		input.SetText("")
	})

	ui, err = tui.New(chat)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

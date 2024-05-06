package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"net"
	"os"
	"strings"
)

const (
	MAX_MESSAGE_SIZE = 2024
	DEFAULT_PORT     = "127.0.0.1:8080"
	HELP_MESSAGE     = `
Help menu
    /quit               Shuts down the client and closes the process
    /addr               Displays the socket address
    /limit              Displays the character limit of a message you can send to the server
    /help               Displays this message
    /file_send <FILE>   Attempts to read FILE's contents then send it to the server `
)

var history []string

func Login() string {
    var userName string
    fmt.Print("Enter your username: ")
    _, err := fmt.Scanln(&userName)
    if err != nil {
        log.Fatal(err)
    }
    return userName
}

func IsCommand(text string) bool {
	return strings.HasPrefix(text, "/")
}

func HandleCommand(command string, conn *net.TCPConn, textView *tview.TextView) {
	commands := strings.Fields(command)
	switch len(commands) {
	case 1:
		switch commands[0] {
		case "/help":
			history = append(history, HELP_MESSAGE)
		case "/limit":
			history = append(history, fmt.Sprintf("Max message size is %d characters\n", MAX_MESSAGE_SIZE))
		case "/addr":
			history = append(history, fmt.Sprintf("Server address: %s\n", conn.RemoteAddr()))
		default:
			history = append(history, "Error: Didn't recognize the command, try using /help\n")
		}

	case 2:
		switch commands[0] {
		case "/file_send":
			body, err := os.ReadFile(commands[1])
			if err != nil {
				history = append(history, err.Error())
			} else {
				_, err := conn.Write(body)
				if err != nil {
					log.Fatal(err)
				}
				history = append(history, fmt.Sprintf("Successfully sent %s to the server\n", commands[1]))
			}
		default:
			history = append(history, "Error: Didn't recognize the command, try using /help\n")
		}
	default:
		history = append(history, "Error: Too many args passed, try using the help command\n")
	}
}

func HandleRead(conn *net.TCPConn, textView *tview.TextView, app *tview.Application) {
	for {
		reply := make([]byte, MAX_MESSAGE_SIZE)
		_, err := conn.Read(reply)
		if err != nil {
			log.Fatal(err)
		}
		history = append(history, string(reply))
		textView.SetText(strings.Join(history, "\n"))
		app.ForceDraw()
	}
}

func main() {
	var serverAddress string
	if len(os.Args) < 2 {
		serverAddress = DEFAULT_PORT
	} else {
		serverAddress = os.Args[1]
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	textView := tview.NewTextView().SetScrollable(true)
	textView.SetBackgroundColor(tcell.ColorDefault)
	inputField := tview.NewInputField().SetPlaceholder("Send a message").SetLabel("> ")
	inputField.SetBackgroundColor(tcell.ColorDefault)

	app := tview.NewApplication()

    userName := Login()

	inputField.SetDoneFunc(func(key tcell.Key) {
		input := inputField.GetText()
		if IsCommand(input) {
			if input == "/quit" {
				app.Stop()
			} else {
				HandleCommand(input, conn, textView)
			}
		} else {
			history = append(history, "you: "+input)
			_, err := conn.Write([]byte(fmt.Sprintf("%s: ", userName)+input))
			if err != nil {
				log.Fatal(err)
			}
		}
		inputField.SetText("")
		textView.SetText(strings.Join(history, "\n"))
	})

	go HandleRead(conn, textView, app)

	flex := tview.NewFlex().
		AddItem(textView, 0, 1, false).
		AddItem(inputField, 1, 0, true).SetDirection(tview.FlexRow)

	app.SetRoot(flex, true).EnableMouse(true).Run()
}

package main

import (
	"fmt"
	"os"
	"telegram-api/core"
)

func usage() {
	fmt.Print("Telegram is a simple MTProto tool.\n\nUsage:\n\n")
	fmt.Print("    ./telegram <command> [arguments]\n\n")
	fmt.Print("The commands are:\n\n")
	fmt.Print("    auth  <phone_number>            auth connection by code\n")
	fmt.Print("    msg   <user_id> <msgtext>       send message to user\n")
	fmt.Print("    list                            get contact list\n")
	fmt.Println()
}

func main() {
	var err error

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	commands := map[string]int{"auth": 1, "msg": 2, "dialogs": 0, "list": 0}
	valid := false
	for k, v := range commands {
		if os.Args[1] == k {
			if len(os.Args) < v+2 {
				usage()
				os.Exit(1)
			}
			valid = true
			break
		}
	}

	if !valid {
		usage()
		os.Exit(1)
	}
	appConfig, err := core.NewConfig(41994,
					"269069e15c81241f5670c397941016a2",
					"0.0.1",
					"",
					"",
					"")
	if err != nil {
		fmt.Printf("Create failed: %s\n", err)
		os.Exit(2)
	}
	m, err := core.NewMTProto(os.Getenv("HOME") + "/.telegram_go", appConfig)
	if err != nil {
		fmt.Printf("Create failed: %s\n", err)
		os.Exit(2)
	}

	err = m.Connect()
	if err != nil {
		fmt.Printf("Connect failed: %s\n", err)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "auth":
		err = m.Auth(os.Args[2])
	//case "msg":
	//	user_id, _ := strconv.Atoi(os.Args[2])
	//	err = m.SendMessage(int32(user_id), os.Args[3])
	//
	//case "list":
	//	err = m.GetContacts()
	//case "dialogs": {
	//	dialogs, users, err := m.GetDialogs(int32(0), int32(0), int32(100))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println("Users: ", len(users))
	//	fmt.Println("Dialogs: ", dialogs)
	//	err = m.GetChats(dialogs)
	//	// err = m.GetRecentGeochats(0, 1000)
	//	// err = m.GetUsers(users)
	//}
	default:
		err = fmt.Errorf("Unknown command %s\n", os.Args[1])
	}


	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
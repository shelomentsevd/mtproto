package main

import (
	"fmt"
	"os"
	"mtproto"
	"reflect"
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
	const telegramAddress = "149.154.167.40:443"
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
	appConfig, err := mtproto.NewConfiguration(41994,
		"269069e15c81241f5670c397941016a2",
		"0.0.1",
		"",
		"",
		"")
	if err != nil {
		fmt.Printf("Create failed: %s\n", err)
		os.Exit(2)
	}
	m, err := mtproto.NewMTProto(false, telegramAddress, os.Getenv("HOME")+"/.telegram_go", *appConfig)
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
		phonenumber := os.Args[2]
		err, authSentCode := m.AuthSendCode(phonenumber)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("Phone code hash: ", authSentCode.Phone_code_hash)
		fmt.Println("Phone registered: ", authSentCode.Phone_registered)
		fmt.Println("Flags: ", authSentCode.Flags)
		fmt.Println("Code type: ", reflect.TypeOf(authSentCode.Code_type))
		fmt.Println("Next type: ", reflect.TypeOf(authSentCode.Next_type))
		fmt.Println("Timeout: ", authSentCode.Timeout)
		if !authSentCode.Phone_registered {
			fmt.Println("Cannot sign in: Phone isn't registered")
			os.Exit(2)
		}
		var code string
		fmt.Printf("Enter code: ")
		fmt.Scanf("%s", &code)
		err, auth := m.AuthSignIn(phonenumber, code, authSentCode.Phone_code_hash)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		userSelf := auth.User.(mtproto.TL_user)
		fmt.Printf("Signed in: Id %d name <%s %s>\n", userSelf.Id, userSelf.First_name, userSelf.Last_name)
	//case "msg":
	//	user_id, _ := strconv.Atoi(os.Args[2])
	//	err = m.SendMessage(int32(user_id), os.Args[3])
	//
	//case "peers":
	//	err, topPeers = m.GetTopPeers()
	//case "dialogs": {
	//	dialogs, users, err := m.MessagesGetDialogs(int32(0), int32(0), int32(100))
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

	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(2)
	//}
}

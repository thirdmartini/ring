package main

import (
	"fmt"
	"flag"
	"time"

	"github.com/thirdmartini/ring"
)

func main() {
	username := flag.String("username", "", "e-mail address used as username")
	password := flag.String("password", "", "password")
	saveRecCount := flag.Int("save-recordings", 0, "number of recordings to save")
	flag.Parse()

	if *username == "" || *password == "" {
		panic("Need to provide --username and --password")
	}

	r, err := ring.New(*username, *password)
	if err != nil {
		fmt.Println("1:")
		panic(err)
	}

	prof, err := r.Profile()
	if err != nil {
		fmt.Println("2:")
		panic(err)
	}

	fmt.Println("Account")
	fmt.Println("    Name:", prof.FirstName, prof.LastName)
	fmt.Println("  E-Mail:", prof.EMail)

	devs, err := r.Devices()
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Println("Devices:")
	for _, bot := range devs.Doorbots {
		fmt.Println("    ", bot.Description, bot.Address)
	}

	history, err := r.History(10)
	if err != nil {
		panic(err)
	}
	fmt.Println("")
	fmt.Println("History:")
	for _, h := range history {
		fmt.Println("    ", h.String())
	}

	if *saveRecCount > len(history) {
		*saveRecCount = len(history)
	}

	for idx :=int(0); idx < *saveRecCount; idx++ {
		h := history[idx]
		fmt.Println("Saving Recording:", fmt.Sprintf("saved-recording-%d.mp4", h.Id))
		err = r.Recording(h.Id, fmt.Sprintf("saved-recording-%d.mp4", h.Id))
		if err!= nil {
			panic(err)
		}
	}

	// Listen for doorbell events  checking every 20 seconds
	//   note that there seems to be a rate limit applied if you poll < ~17 sec
	//
	fmt.Println("")
	fmt.Println("Listening for Doorbot events:")
	err = r.Listen(time.Second*20, func(d *ring.Ding) {
		fmt.Println("   ", d.DoorbotDescription, d.SipServerAddress)
	})

	if err != nil {
		panic(err)
	}

}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	evdev "github.com/gvalkov/golang-evdev"
	"github.com/lotib/uinput"
)

var (
	flagVirtualDevice = flag.String("virtual-device", "testkeyboard", "virtual device injects received event")
	flagPort          = flag.Int("port", 36666, "UDP port")
	flagAddr          = flag.String("addr", "", "address to listen")
)

func main() {
	flag.Parse()

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte(*flagVirtualDevice))
	if err != nil {
		return
	}
	// always do this after the initialization in order to guarantee that the device will be properly closed
	defer keyboard.Close()
	fmt.Println("Virtual Keyboard created")

	listenAddr := fmt.Sprintf("%s:%d", *flagAddr, *flagPort)

	s, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	fmt.Println("Server started on ", listenAddr)

	for {
		// Read a msg at a time
		buf := make([]byte, 24)
		n, addr, err := s.ReadFrom(buf)
		if err != nil {
			continue
		}

		fmt.Printf("Received from %s %v (%d)", addr, buf, n)

		//var toto *event.InputEvent

		event := (*(*evdev.InputEvent)(unsafe.Pointer(&buf[0])))

		println("Where am I")

		fmt.Printf("%d %d %d \n", event.Code, event.Type, event.Value)
		fmt.Printf("sqdsqdqsds  %d %d %d \n", event.Code, event.Type, event.Value)

		fmt.Printf("%v \n", event)
		fmt.Printf("%s \n", event.String())

		println("Where am I")

		keyboard.SendBufferEvent(buf)

		// if n != 24 {
		// }

	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("received interrupt: %s", <-sig)
}

// u, err := evdev.NewUserInput(
// 	0644,
// 	evdev.WithID(evdev.ID{
// 		BusType: evdev.BusUSB,
// 		Vendor:  0x01,
// 		Product: 0x02,
// 		Version: 0x0a0b,
// 	}),
// 	evdev.WithName(*flagOutput),
// 	evdev.WithPath(*flagOutput),
// )
// if err != nil {
// 	log.Fatal(err)
// }
// defer u.Close()

// // create context
// ctxt, cancel := context.WithCancel(context.Background())
// defer cancel()

// // start polling and relaying
// in := d.Poll(ctxt)
// go func() {
// 	for {
// 		select {
// 		case <-ctxt.Done():
// 			return

// 		case event := <-in:
// 			if event == nil {
// 				return
// 			}
// 			log.Printf("<- %+v", event)
// 			go u.Send(*event)
// 		}
// 	}
// }()
// log.Printf("created: %s", u.Path())

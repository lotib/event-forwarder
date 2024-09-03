package main

import (
	"bytes"
	"encoding/binary"
	"flag"

	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	evdev "github.com/gvalkov/golang-evdev"
)

var (
	flagListDevices = flag.Bool("list", false, "list input devices that can be attached")

	flagGrabDevice  = flag.Bool("grab", false, "grab events from device exclusively")
	flagInputDevice = flag.String("input-device", "", "input device to forward")

	flagForwardEvent = flag.Bool("forward", false, "forward events over network")
	flagAddr         = flag.String("addr", "127.0.0.1", "destination address / host")
	flagPort         = flag.Int("port", 36666, "destination UDP port")

	flagDebug = flag.Bool("debug", false, "add verbosity")

	conn *net.UDPConn = nil
)

func main() {
	flag.Parse()

	// device, _ := Open
	// fmt.Println(device)

	if *flagListDevices {
		devices, err := evdev.ListInputDevices()

		if err != nil {
			fmt.Printf("Failed to list input devices %v", err)
			os.Exit(1)
		}
		for _, dev := range devices {
			fmt.Printf("%s %s %s \n", dev.Fn, dev.Name, dev.Phys)
		}
	}

	if len(*flagInputDevice) > 0 {
		device, err := evdev.Open(*flagInputDevice)

		if err != nil {
			fmt.Printf("Failed to open input device %v", err)
			os.Exit(1)
		}

		fmt.Println(device)
		fmt.Println(device.Capabilities)
		// fmt.Println(device.ResolveCapabilities())

		if *flagGrabDevice {
			device.Grab()
			defer device.Release()
			fmt.Printf("Device %s attached exclusively\n", *flagInputDevice)
		}

		if *flagForwardEvent {

			addr := fmt.Sprintf("%s:%d", *flagAddr, *flagPort)

			// Define the UDP server destination
			serverAddr, err := net.ResolveUDPAddr("udp", addr)
			if err != nil {
				fmt.Println("Error resolving UDP address:", err)
				os.Exit(1)
			}

			// Create a UDP connection
			conn, err = net.DialUDP("udp", nil, serverAddr)
			if err != nil {
				fmt.Println("Error dialing UDP:", err)
				os.Exit(1)
			}
			defer conn.Close()

			fmt.Printf("%s events will be forwarded to %s", *flagInputDevice, addr)
		}

		for {
			event, err := device.ReadOne()
			if err != nil {
				fmt.Printf("Failed to read input event %v", err)
				os.Exit(1)
			}

			if *flagDebug {
				fmt.Println(event)
			}

			if conn != nil {

				// buf := new(bytes.Buffer)
				buf := bytes.NewBuffer(make([]byte, 0, 24))

				err := binary.Write(buf, binary.LittleEndian, event)

				if err != nil {
					fmt.Println("Error serializing struct:", err)
					os.Exit(1)
				}
				data := buf.Bytes()
				_, err = conn.Write(data)
				if err != nil {
					fmt.Println("Error sending message:", err)
					os.Exit(1)
				}

				if *flagDebug {
					fmt.Println(*buf)
				}
			}
		}
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		log.Printf("received interrupt: %s", <-sig)
	}
}

func sendEvent() {

}

// d, err := evdev.OpenFile(*flagInput)
// if err != nil {
// 	log.Fatal(err)
// }
// defer d.Close()

//  // Define the UDP server destination
// serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", *flagDestAddr, *flagDestPort))
// if err != nil {
// 	fmt.Println("Error resolving UDP address:", err)
// 	os.Exit(1)
// }

// // Create a UDP connection
// conn, err := net.DialUDP("udp", nil, serverAddr)
// if err != nil {
// 	fmt.Println("Error dialing UDP:", err)
// 	os.Exit(1)
// }
// defer conn.Close()

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

// 			buf := new(bytes.Buffer)
// 			err := binary.Write(buf, binary.LittleEndian, event.Event)
// 			if err != nil {
// 				fmt.Println("Error serializing struct:", err)
// 				os.Exit(1)
// 			}
// 			data := buf.Bytes()

// 			// Send the message
// 			_, err = conn.Write(data)
// 			if err != nil {
// 				fmt.Println("Error sending message:", err)
// 				os.Exit(1)
// 			}
// 			fmt.Println("Message sent to", serverAddr)
// 		}
// 	}
// }()

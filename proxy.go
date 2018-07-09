package mcpeproxy

import (
	"net"
	"os"
	"strconv"
)

type Proxy struct {
	*net.UDPConn
	Config Config
}

const (
	Prefix  = "GoProxy > "
	Author  = "xBeastMode"
	Version = "1.2"
)

var LogoType = []string{
	"  _____       _____",
	" / ____|     |  __ \\",
	"| |  __  ___ | |__) | __ _____  ___   _ ",
	"| | |_ |/ _ \\|  ___/ '__/ _ \\ \\/ / | | |",
	"| |__| | (_) | |   | | | (_) >  <| |_| |",
	" \\_____|\\___/|_|   |_|  \\___/_/\\_\\\\__, |",
	"                                   __/ |",
	"                                  |___/ ",
}

// this function starts the Proxy
// it will start "sniffing" packets
func StartProxy() {
	var config = NewConfig()
	var proxy = Proxy{Config:config}
	var err error

	// forward the Proxy to a custom port
	proxy.UDPConn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(config.BindAddr), Port: config.BindPort})

	if err != nil {
		Panic(err.Error())
		os.Exit(1)
	}

	// just typical console logs
	for _, v := range LogoType {
		Info(AnsiBrightRed + v)
	}
	Info(AnsiGreen + Author + "'s Proxy version " + Version)
	Info("Starting Proxy on " + config.BindAddr + ":" + strconv.Itoa(config.BindPort))

	if err != nil {
		Panic(err.Error())
		os.Exit(1)
	}

	// looks up the Server domain and returns it address
	addrs, err := net.LookupHost(config.ServerAddr)

	if err != nil {
		Panic(err.Error())
		os.Exit(1)
	}

	// resolve Server address
	addr, err := net.ResolveUDPAddr("udp", addrs[0] + ":" + strconv.Itoa(config.ServerPort))

	if err != nil {
		Panic(err.Error())
		os.Exit(1)
	}

	// start a new Connection between Server and Client
	conn := NewConnection(&proxy)
	conn.Server.SetAddress(addr)

	Info("Listening to Server address: " + addrs[0] + ":" + strconv.Itoa(config.ServerPort))

	// start the packet "sniffing" from Client
	conn.HandleIncomingPackets()
}
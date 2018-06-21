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
	Version = "1.1"
)

// this function starts the proxy
// it will start "sniffing" packets
func StartProxy() {
	var config = NewConfig()
	var proxy = Proxy{Config:config}
	var err error

	// forward the proxy to a custom port
	proxy.UDPConn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(config.BindAddr), Port: config.BindPort})

	if err != nil {
		Panic(err.Error())
		os.Exit(1)
	}

	// just typical console logs
	Info(AnsiGreen + Author + "'s proxy version " + Version)
	Info("Starting proxy on " + config.BindAddr + ":" + strconv.Itoa(config.BindPort))

	if err != nil {
		Panic(err.Error())
		os.Exit(1)
	}

	// looks up the server domain and returns it address
	addrs, err := net.LookupHost(config.ServerAddr)

	if err != nil {
		Panic(err.Error())
		os.Exit(1)
	}

	// resolve server address
	addr, err := net.ResolveUDPAddr("udp", addrs[0] + ":" + strconv.Itoa(config.ServerPort))

	if err != nil {
		Panic(err.Error())
		os.Exit(1)
	}

	// start a new connection between server and client
	conn := NewConnection(&proxy)
	conn.server.SetAddress(addr)

	Info("Listening to server address: " + addrs[0] + ":" + strconv.Itoa(config.ServerPort))

	// start the packet "sniffing" from client
	conn.HandleIncomingPackets()
}
package aprsis

import (
	"bufio"
	"fmt"
	"net"

	"github.com/so5dz/utils/misc"
)

const _Software = "SO5DZ_Station_Monitor"
const _Version = "v0.0.0"

type PacketCallback func(tnc2Bytes []byte)

type APRSIS struct {
	connection     net.Conn
	packetCallback PacketCallback
}

func (a *APRSIS) OnPacket(packetCallback PacketCallback) {
	a.packetCallback = packetCallback
}

func (a *APRSIS) Start(host string, port int, user string, pass int, filter string) error {
	err := a.connect(host, port)
	if err != nil {
		return misc.WrapError("establishing connection to server", err)
	}
	err = a.login(user, pass, filter)
	if err != nil {
		return misc.WrapError("authenticating with server", err)
	}
	go a.readLoop()
	return nil
}

func (a *APRSIS) connect(host string, port int) error {
	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return misc.WrapError("connecting to address", err)
	}
	a.connection = connection
	return nil
}

func (a *APRSIS) login(user string, pass int, filter string) error {
	loginLine := fmt.Sprintf("user %s pass %d vers %s %s filter %s\r\n",
		user, pass, _Software, _Version, filter,
	)
	_, err := a.connection.Write([]byte(loginLine))
	if err != nil {
		return misc.WrapError("writing to socket", err)
	}
	return nil
}

func (a *APRSIS) readLoop() {
	scanner := bufio.NewScanner(a.connection)
	for {
		if ok := scanner.Scan(); !ok {
			break
		}
		packetBytes := scanner.Bytes()
		if len(packetBytes) > 0 && packetBytes[0] != '#' && a.packetCallback != nil {
			a.packetCallback(packetBytes)
		}
	}
}

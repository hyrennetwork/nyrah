package connection

import (
	"fmt"
	"io"
	"net"
	"net/hyren/nyrah/minecraft/protocol"

	ProxyApp "net/hyren/nyrah/applications"
)

func copy(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func SendToProxy(connection *protocol.Connection, name string) {
	var inetSocketAddress = ProxyApp.GetProxyAddress(name)

	ds, err := net.Dial("tcp", fmt.Sprintf("%s:%d", inetSocketAddress.GetHostAddress(), inetSocketAddress.GetPort()))

	if err != nil {
		_ = connection.Close()
	}

	us := connection.Handle

	go copy(ds, us)

	bg := protocol.NewConnection(ds)

	for _, item := range connection.PacketQueue {
		_, _ = bg.Write(item)
	}

	go copy(us, ds)
}

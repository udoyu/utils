package thriftst

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"os"
)

type ThriftSt struct {
	tsocket          *thrift.TSocket
	ttransport       thrift.TTransport
	tprotocolFactory thrift.TProtocolFactory
}

func (this *ThriftSt) TTransport() thrift.TTransport {
	return this.ttransport
}

func (this *ThriftSt) TProtocolFactory() thrift.TProtocolFactory {
	return this.tprotocolFactory
}

func NewThriftSt(addr, port string) (*ThriftSt, error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	tprotocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	tsocket, err := thrift.NewTSocket(net.JoinHostPort(addr, port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		return nil, err
	}
	ttransport := transportFactory.GetTransport(tsocket)
	p := &ThriftSt{}
	p.tsocket = tsocket
	p.tprotocolFactory = tprotocolFactory
	p.ttransport = ttransport
	return p, err
}

func (this *ThriftSt) Start() error {
	err := this.Open()
	if err != nil {
		return err
	}
	defer this.Close()
	return nil
}

func (this *ThriftSt) Open() error {
	//client = user.NewUserServiceClientFactory(clientTransport, protocolFactory)
	return this.tsocket.Open()
}

func (this *ThriftSt) Close() {
	this.tsocket.Close()
}

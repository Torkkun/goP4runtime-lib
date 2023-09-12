package simpleswitch

import (
	"context"
	"fmt"
	"log"

	v1 "github.com/p4lang/p4runtime/go/p4/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const MSG_LOG_MAX_LEN = 1024

var connections []*SwitchConnection

func ShutdownAllSwitchConnections() {
	for _, c := range connections {
		c.Shutdown()
	}
}

type SwitchConnection struct {
	Name          string
	Address       string
	DeviceId      uint64
	ProtoDumpFile string
	Channel       *grpc.ClientConn
	Client        *v1.P4RuntimeClient
	RequestStream v1.P4Runtime_StreamChannelClient
	StreamMsgResp *v1.StreamMessageResponse
}

func NewSwitchConnection(name string, address string, deviceid uint64, protodumpfile string) SwitchConnection {
	newswcon := new(SwitchConnection)
	newswcon.Name = name
	newswcon.Address = address
	newswcon.DeviceId = deviceid
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("NewConnection failed:%v", err)
	}
	defer conn.Close()

	newswcon.Channel = conn
	/*if protodumpfile != "" {
		//TODO:Intercept Client
	}*/
	client := v1.NewP4RuntimeClient(conn)
	newswcon.Client = &client
	requestStream, err := client.StreamChannel(context.Background())
	if err != nil {
		log.Fatalf("Request Stream Channel faital error: %v", err)
	}
	newswcon.RequestStream = requestStream
	newswcon.ProtoDumpFile = protodumpfile
	if connections == nil {
		connections = make([]*SwitchConnection, 0)
	}
	connections = append(connections, newswcon)
	return *newswcon
}

func (swcon *SwitchConnection) Shutdown() {
	if err := swcon.Channel.Close(); err != nil {
		log.Printf("channel close shutdownfunc error : %v", err)
	}
	if err := swcon.RequestStream.CloseSend(); err != nil {
		log.Printf("close send shutdownfunc error: %v", err)
	}
}

func (swcon *SwitchConnection) MasterArbitrationUpdate(dryrun bool, opt ...interface{}) {
	request := v1.StreamMessageRequest{
		Update: &v1.StreamMessageRequest_Arbitration{
			Arbitration: &v1.MasterArbitrationUpdate{
				DeviceId: swcon.DeviceId,
				ElectionId: &v1.Uint128{
					High: 0,
					Low:  1,
				},
			},
		},
	}

	if dryrun {
		fmt.Printf("P4Runtime MasterArbitrationUpdate: %d", request)
	} else {
		// TODO: return value
		if err := swcon.RequestStream.Send(&request); err != nil {
			log.Fatalln("MasterArbitrationUpdate channel send error: %v", err)
		}
		streamMsgResp, err := swcon.RequestStream.Recv()
		if err != nil {
			log.Fatalln("MasterArbitrationUpdate channel recv error: %v", err)
		}
		swcon.StreamMsgResp = streamMsgResp
	}
}

func (swcon *SwitchConnection) SetForwardingPipelineConfig() {

}

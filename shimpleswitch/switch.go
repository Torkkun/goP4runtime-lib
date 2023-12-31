package simpleswitch

import (
	"context"
	"log"

	v1conf "github.com/p4lang/p4runtime/go/p4/config/v1"
	v1 "github.com/p4lang/p4runtime/go/p4/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const MSG_LOG_MAX_LEN = 1024

var connections []*SwitchConnection

func ShutdownAllSwitchConnections() {
	for _, c := range connections {
		c.shutdown()
	}
}

type SwitchConnection struct {
	Name          string
	Address       string
	DeviceId      uint64
	ProtoDumpFile string
	Channel       *grpc.ClientConn
	Client        v1.P4RuntimeClient
	RequestStream v1.P4Runtime_StreamChannelClient
	StreamMsgResp *v1.StreamMessageResponse
}

func NewSwitchConnection(name string, address string, deviceid uint64, protodumpfile string) *SwitchConnection {
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
	//defer conn.Close()

	newswcon.Channel = conn
	/*if protodumpfile != "" {
		//TODO:Intercept Client
	}*/
	client := v1.NewP4RuntimeClient(conn)
	newswcon.Client = client
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
	return newswcon
}

type BuildDeviceConfig interface {
	BuildDeviceConfig(string) []byte
}

func NewBuildDeviceConfig(builddevconf BuildDeviceConfig, config string) []byte {
	return builddevconf.BuildDeviceConfig(config)
}

func (swcon *SwitchConnection) shutdown() {
	if err := swcon.Channel.Close(); err != nil {
		log.Printf("channel close shutdownfunc error : %v", err)
	}
	if err := swcon.RequestStream.CloseSend(); err != nil {
		log.Printf("close send shutdownfunc error: %v", err)
	}
}

func (swcon *SwitchConnection) MasterArbitrationUpdate() {
	request := &v1.StreamMessageRequest{
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

	// RequestStreamがNil
	if err := swcon.RequestStream.Send(request); err != nil {
		log.Fatalf("MasterArbitrationUpdate channel send error: %v", err)
	}
	streamMsgResp, err := swcon.RequestStream.Recv()
	if err != nil {
		log.Fatalf("MasterArbitrationUpdate channel recv error: %v", err)
	}
	swcon.StreamMsgResp = streamMsgResp
}

func (swcon *SwitchConnection) SetForwardingPipelineConfig(p4info *v1conf.P4Info, p4devconf []byte) error {
	devconf := new(v1.ForwardingPipelineConfig)
	devconf.P4Info = p4info
	devconf.P4DeviceConfig = p4devconf
	request := &v1.SetForwardingPipelineConfigRequest{
		DeviceId: swcon.DeviceId,
		ElectionId: &v1.Uint128{
			Low: 1,
		},
		Config: devconf,
		Action: v1.SetForwardingPipelineConfigRequest_VERIFY_AND_COMMIT,
	}
	_, err := swcon.Client.SetForwardingPipelineConfig(context.Background(), request)
	if err != nil {
		return err
	}
	return nil
}

func (swcon *SwitchConnection) WriteTableEntry(te *v1.TableEntry) {
	newupdate := new(v1.Update)
	if te.IsDefaultAction {
		newupdate.Type = v1.Update_MODIFY
	} else {
		newupdate.Type = v1.Update_INSERT
	}
	newupdate.Entity = &v1.Entity{
		Entity: &v1.Entity_TableEntry{
			TableEntry: te,
		},
	}
	request := &v1.WriteRequest{
		DeviceId: swcon.DeviceId,
		ElectionId: &v1.Uint128{
			Low: 1,
		},
		Updates: []*v1.Update{newupdate},
	}
	swcon.Client.Write(context.Background(), request)
}

func (swcon *SwitchConnection) ReadTableEntry(tid uint32) (*v1.ReadResponse, error) {
	te := new(v1.TableEntry)
	te.TableId = tid
	request := new(v1.ReadRequest)
	request.DeviceId = swcon.DeviceId
	request.Entities =
		[]*v1.Entity{
			{
				Entity: &v1.Entity_TableEntry{
					TableEntry: te,
				},
			},
		}
	cl, err := swcon.Client.Read(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return cl.Recv()
}

func (swcon *SwitchConnection) ReadCouter(cnterid uint32, index int64) (*v1.ReadResponse, error) {
	cnte := new(v1.CounterEntry)
	cnte.CounterId = cnterid
	cnte.Index = &v1.Index{
		Index: index,
	}
	request := new(v1.ReadRequest)
	request.DeviceId = swcon.DeviceId
	request.Entities = []*v1.Entity{
		{
			Entity: &v1.Entity_CounterEntry{
				CounterEntry: cnte,
			},
		},
	}
	cl, err := swcon.Client.Read(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return cl.Recv()
}

// yet unimplement
func WritePREEntry() {}

type GrpcRequestLogger struct{}

func NewGrpcRequestLogger() {}

func (GrpcRequestLogger) LogMessage() {}

func (GrpcRequestLogger) InterceptUnaryUnary() {}

func (GrpcRequestLogger) InterceptUnaryStream() {}

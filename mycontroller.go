package main

import (
	"flag"
	"fmt"
	"gop4runtimelib/architecture"
	"gop4runtimelib/helper"
	simpleswitch "gop4runtimelib/shimpleswitch"
	"log"
	"os"
	"os/signal"
	"time"
)

var SWITCH_TO_HOST_PORT = 1
var SWITCH_TO_SWITCH_PORT = 2

func main() {
	p4ipath := flag.String(
		"p4info",
		"./build/advanced_tunnel.p4.p4info.txt",
		"p4info proto in text format from p4c")

	bmv2jpath := flag.String(
		"bmv2json",
		"./build/advanced_tunnel.json",
		"BMv2 JSON file from p4c")
	flag.Parse()

	if _, err := os.Stat(*p4ipath); err != nil {
		log.Fatalf("File %s does not exist!: %v", *p4ipath, err)
	}

	if _, err := os.Stat(*bmv2jpath); err != nil {
		log.Fatalf("File %s does not exist!: %v", *bmv2jpath, err)
	}

	// Instantiate a P4Runtime helper from the p4info file
	p4ih := helper.NewP4InfoHelper(*p4ipath)

	// Create a switch connection object for s1 and s2;
	// this is backed by a P4Runtime gRPC connection.
	// Also, dump all P4Runtime messages sent to switch to given txt files.
	s1 := architecture.Bmv2SwitchConnection{
		SwitchConnection: simpleswitch.NewSwitchConnection(
			"s1",
			"127.0.0.1:50051",
			0,
			"logs/s1-p4runtime-requests.txt",
		),
	}
	s2 := architecture.Bmv2SwitchConnection{
		SwitchConnection: simpleswitch.NewSwitchConnection(
			"s2",
			"127.0.0.1:50052",
			1,
			"logs/s2-p4runtime-requests.txt",
		),
	}

	// Send master arbitration update message to establish this controller as
	// master (required by P4Runtime before performing any other write operation)
	s1.MasterArbitrationUpdate()
	s2.MasterArbitrationUpdate()

	// Install the P4 program on the switches
	p4devconf := simpleswitch.NewBuildDeviceConfig(s1, *bmv2jpath)
	s1.SetForwardingPipelineConfig(p4ih.P4Info, p4devconf)
	fmt.Println("Installed P4 Program using SetForwardingPipelineConfig on s1")
	p4devconf = simpleswitch.NewBuildDeviceConfig(s2, *bmv2jpath)
	s2.SetForwardingPipelineConfig(p4ih.P4Info, p4devconf)
	fmt.Println("Installed P4 Program using SetForwardingPipelineConfig on s2")

	// Write the rules that tunnel traffic from h1 to h2
	writeTunnelRules(p4ih, s1.SwitchConnection, s2.SwitchConnection, 100, "08:00:00:00:02:22", "10.0.2.2")

	// Write the rules that tunnel traffic from h2 to h1
	writeTunnelRules(p4ih, s1.SwitchConnection, s2.SwitchConnection, 200, "08:00:00:00:01:11", "10.0.1.1")

	readTablesRules(p4ih, s1)
	readTablesRules(p4ih, s2)

	go func() {
		for {
			time.Sleep(time.Duration(2))
			fmt.Println("\n----- Reading tunnel counters -----")
			if err := printCounter(p4ih, s1, "MyIngress.ingressTunnelCounter", 100); err != nil {
				log.Println(err)
			}
			if err := printCounter(p4ih, s2, "MyIngress.egressTunnelCounter", 100); err != nil {
				log.Println(err)
			}
			if err := printCounter(p4ih, s2, "MyIngress.ingressTunnelCounter", 200); err != nil {
				log.Println(err)
			}
			if err := printCounter(p4ih, s1, "MyIngress.egressTunnelCounter", 200); err != nil {
				log.Println(err)
			}
		}
	}()
	// Ctrl + Cを検知
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println(" Shut down. ")
	simpleswitch.ShutdownAllSwitchConnections()
}

func writeTunnelRules(
	p4ih *helper.P4InfoHelper,
	ingresssw *simpleswitch.SwitchConnection,
	egresssw *simpleswitch.SwitchConnection,
	tunid uint32, dstethaddr,
	dstipaddr string,
) {
	//Tunnel Ingress Rule
	tableentry, err := p4ih.BuildTableEntry(
		"MyIngress.ipv4_lpm",
		&helper.IsEntryMatchField{
			Name: "hdr.ipv4.dstAddr",
			EntryMatchField: &helper.EntryMatchFieldLpm{
				Value0: dstethaddr,
				Value1: 32,
			},
		},
		helper.Action(
			"MyIngress.myTunnel_ingress",
			map[string]interface{}{
				"dst_id": tunid,
			},
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	ingresssw.WriteTableEntry(tableentry)
	fmt.Printf("Installed ingress tunnel rule on %s\n", ingresssw.Name)
	//Tunnel Transit Rule

	tableentry, err = p4ih.BuildTableEntry(
		"MyIngress.myTunnel_exact",
		&helper.IsEntryMatchField{
			Name: "hdr.myTunnel.dst_id",
			EntryMatchField: &helper.EntryMatchFieldExact{
				Value0: tunid,
			},
		},
		helper.Action(
			"MyIngress.myTunnel_forward",
			map[string]interface{}{
				"dstAddr": dstethaddr,
				"port":    uint32(SWITCH_TO_SWITCH_PORT),
			},
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	ingresssw.WriteTableEntry(tableentry)
	fmt.Println("TODO Install transit tunnel rule")

	tableentry, err = p4ih.BuildTableEntry(
		"MyIngress.myTunnel_exact",
		&helper.IsEntryMatchField{
			Name: "hdr.myTunnel.dst_id",
			EntryMatchField: &helper.EntryMatchFieldExact{
				Value0: tunid,
			},
		},
		helper.Action(
			"MyIngress.myTunnel_egress",
			map[string]interface{}{
				"dstAddr": dstethaddr,
				"port":    uint32(SWITCH_TO_HOST_PORT),
			},
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	egresssw.WriteTableEntry(tableentry)
	fmt.Printf("Install egress tunnel rule on %s\n", egresssw.Name)
}

func readTablesRules(
	p4ih *helper.P4InfoHelper,
	sw architecture.Bmv2SwitchConnection,
) error {
	fmt.Printf("\n----- Reading tables rules for %s -----", sw.Name)
	resp, err := sw.ReadTableEntry(0)
	if err != nil {
		return err
	}
	for _, entt := range resp.Entities {
		te := entt.GetTableEntry()
		tname, err := p4ih.GetName(helper.TableEntity, te.TableId)
		if err != nil {
			return err
		}
		fmt.Printf("%s: ", tname)
		for _, m := range te.Match {
			mf, err := p4ih.GetMatchFieldWithId(tname, m.FieldId)
			if err != nil {
				return err
			}
			fmt.Printf("%s ", mf.Name)
			v, err := p4ih.GetMatchFieldValue(mf, m)
			if err != nil {
				return err
			}
			fmt.Printf("%v ", v)
		}
		action := te.GetAction()
		actionname, err := p4ih.GetName(helper.ActionEntity, action.GetAction().ActionId)
		if err != nil {
			return err
		}
		fmt.Printf("->%s ", actionname)
		for _, p := range action.GetAction().Params {
			ap, err := p4ih.GetActionParamWithId(actionname, p.ParamId)
			if err != nil {
				return err
			}
			fmt.Printf("%s ", ap.Name)
			fmt.Printf("%v ", p.Value)
		}
		fmt.Println()
	}
	return nil
}
func printCounter(
	p4ih *helper.P4InfoHelper,
	sw architecture.Bmv2SwitchConnection,
	cntername string,
	index int64,
) error {
	id, err := p4ih.GetID(helper.CounterEntity, cntername)
	if err != nil {
		return err
	}
	resp, err := sw.ReadCouter(
		id,
		index)
	if err != nil {
		return err
	}
	for _, entt := range resp.Entities {
		cnte := entt.GetCounterEntry()
		fmt.Printf("%s %s %d: %d packets (%d bytes)",
			sw.Name, cntername, index,
			cnte.Data.PacketCount, cnte.Data.ByteCount)
	}
	return nil
}

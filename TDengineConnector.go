package main

import (
	eiicfgmgr "ConfigMgr/eiiconfigmgr"
	eiimsgbus "EIIMessageBus/eiimsgbus"
	"bytes"
	"database/sql/driver"
	"flag"
	"fmt"
	"io"
	"json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	taos "github.com/taosdata/driver-go/v2/af"
)

var cfgMgr *eiicfgmgr.ConfigMgr
var pubclient *eiimsgbus.MsgbusClient
var publisher *eiimsgbus.Publisher
var subclient *eiimsgbus.MsgbusClient
var subscriber *eiimsgbus.Subscriber
var done = make(chan bool)

func startEIIPublisher() {
	pubCtx, err := cfgMgr.GetPublisherByIndex(0)
	if err != nil {
		glog.Errorf("Error occured with error:%v", err)
		return
	}
	topics, err := pubCtx.GetTopics()
	if err != nil {
		glog.Errorf("Failed to fetch topics : %v", err)
		return
	}
	topic := topics[0]
	glog.Infof("Publisher topic is : %s", topic)
	config, err := pubCtx.GetMsgbusConfig()
	if err != nil {
		glog.Error("Failed to get message bus config :%v", err)
		return
	}
	pubclient, err = eiimsgbus.NewMsgbusClient(config)
	if err != nil {
		glog.Errorf("-- Error creating context: %v\n", err)
		return
	}
	publisher, err = pubclient.NewPublisher(topic)
	if err != nil {
		glog.Errorf("-- Error creating publisher: %v\n", err)
		return
	}
	ps := strings.Split(topic, ".") // [dbName, stableName]
	subscribeToTDengine(ps[0], ps[1])
}

func subscribeToTDengine(dbName string, stableName string) {
	conn, err := taos.Open("", "", "", dbName, 0)
	if err != nil {
		glog.Errorf("-- Error connecting tdengine: %v\n", err)
		done <- false
		return
	}
	defer conn.Close()

	// only subscribe new comming data from 10s ago.
	topic, err := conn.Subscribe(false, stableName, "select * from "+stableName+" where _c0 > now - 10s", time.Second)
	if err != nil {
		glog.Errorf("-- Error subscribing tdengine: %v\n", err)
		done <- false
		return
	}
	defer topic.Unsubscribe(false) // not keep progress

	for {
		func() {
			rows, err := topic.Consume()
			if err != nil {
				glog.Errorf("-- Error consuming topic: %v\n", err)
				return
			}
			defer func() { rows.Close(); time.Sleep(time.Second) }()
			cols := rows.Columns()
			for {
				values := make([]driver.Value, len(cols))
				err := rows.Next(values)
				if err == io.EOF {
					break
				} else if err != nil {
					glog.Errorf("%v\n", err)
					break
				}
				var buf bytes.Buffer
				ts, ok := values[0].(time.Time)
				if !ok {
					continue
				}
				buf.WriteString("ts=")
				ms := toMicroseconds(ts)
				buf.WriteString(strconv.FormatInt(ms, 10))
				buf.WriteString(" ")
				for i, col := range cols[1:] {
					s := fmt.Sprintf("%s=%v ", col, values[i+1])
					buf.WriteString(s)
				}
				msg := map[string]interface{}{"data": buf.String()}
				publisher.Publish(msg)
				glog.Infof("Published message: %v", msg)
			}
		}()
	}

}

func startEIISubscriber() {
	subCtx, err := cfgMgr.ConfigMgr.GetSubscriberByIndex(0)
	if err != nil {
		glog.Errorf("Error occured with error:%v", err)
		return
	}
	topics, err := subCtx.GetTopics()
	if err != nil {
		glog.Errorf("Failed to fetch topics : %v", err)
		return
	}
	topic := topics[0]
	glog.Infof("Subscriber topic is : %v", topic)
	config, err := subCtx.GetMsgbusConfig()
	if err != nil {
		glog.Error("Failed to get message bus config :%v", err)
		return
	}
	subclient, err = eiimsgbus.NewMsgbusClient(config)
	if err != nil {
		glog.Errorf("-- Error creating context: %v\n", err)
		return
	}
	subscriber, err = subclient.NewSubscriber(topic)
	if err != nil {
		glog.Errorf("-- Error creating subscriber: %v\n", err)
		return
	}
	for {
		msg := <-subscriber.MessageChannel
		bytemsg, err := json.Marshal(msg.Data)
		if err != nil {
			glog.Errorf("error: %s", err)
		}
		fmt.Println(bytemsg)
		glog.Infof("Subscribe data received from topic: %s  Data: %v", msg.Name, msg.Data)
	}
}

func writeToTDengine(msgData string) {

}

func cleanup() {
	if pubclient != nil {
		pubclient.Close()
	}
	if publisher != nil {
		publisher.Close()
	}
}

func toMicroseconds(t time.Time) int64 {
	return t.Unix()*1e3 + int64(t.Nanosecond()/1000000)
}

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", os.Getenv("GO_LOG_LEVEL"))
	flag.Set("v", os.Getenv("GO_VERBOSE"))
	cfgMgr, _ = eiicfgmgr.ConfigManager()
	appName, err := cfgMgr.GetAppName()
	if err != nil {
		glog.Fatalf("Not able to read appname from etcd")
		return
	}
	glog.Infof("Start %s", appName)

	go startEIIPublisher()
	go startEIISubscriber()
	<-done
	glog.Info("do cleanup...")
	cleanup()
}

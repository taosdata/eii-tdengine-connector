package main

import (
	eiicfgmgr "ConfigMgr/eiiconfigmgr"
	eiimsgbus "EIIMessageBus/eiimsgbus"
	"time"
	"flag"
	"github.com/golang/glog"
	"os"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", os.Getenv("GO_LOG_LEVEL"))
	flag.Set("v", os.Getenv("GO_VERBOSE"))
	glog.Info("Start TDengineConnector")
	cfgMgr, _ := eiicfgmgr.ConfigManager()
	appName, err := cfgMgr.GetAppName()
	if err != nil {
		glog.Fatalf("Not able to read appname from etcd")
		return
	}
	glog.Infof("AppName=%v", appName)
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
	msgbusclient, err  := eiimsgbus.NewMsgbusClient(config)
	if err != nil {
		glog.Errorf("-- Error creating context: %v\n", err)
		return
	}
	publisher, err := msgbusclient.NewPublisher(topic)
	if err != nil {
		glog.Errorf("-- Error creating publisher: %v\n", err)
		return
	}
	for {
		msg := map[string]interface{}{"data": "111"}
		publisher.Publish(msg)
		time.Sleep(time.Duration(2) * time.Second)
	}
}

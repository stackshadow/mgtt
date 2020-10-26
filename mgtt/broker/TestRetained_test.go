package broker

import (
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"gitlab.com/mgtt/cli"
	"gitlab.com/mgtt/plugin"
)

/*
	cmd := exec.Command(
		"mosquitto_pub",
		"-L",
		"mqtts://admin:admin@127.0.0.1:1234/$SYS/broker/cr",
		"-m",
		"2",
		"-d",
		"-q",
		"0",
	)
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
*/

func TestTimeout(t *testing.T) {

	cli.CLI.URL = "tcp://127.0.0.1:1237"
	cli.CLI.CLICommon.Terminal = true
	cli.CLI.CLICommon.Terminal.AfterApply()
	cli.CLI.CLICommon.Debug = true
	cli.CLI.CLICommon.Debug.AfterApply()
	cli.CLI.DBFilename = "test1.db"
	cli.CLI.ConnectTimeout = 1

	// the broker
	os.Remove(cli.CLI.DBFilename)
	server, _ := New()
	go server.Serve(
		Config{
			URL:      cli.CLI.URL,
			CertFile: cli.CLI.CertFile,
			KeyFile:  cli.CLI.KeyFile,
		},
	)
	time.Sleep(time.Second * 2)

	connection, err := net.Dial("tcp", "127.0.0.1:1237")
	if err != nil {
		t.FailNow()
	}

	msg := make([]byte, 4000)

	_, err = connection.Read(msg)
	if err == nil {
		t.FailNow()
	}

}

func TestRetained(t *testing.T) {

	cli.CLI.URL = "tcp://127.0.0.1:1236"
	cli.CLI.CLICommon.Terminal = true
	cli.CLI.CLICommon.Terminal.AfterApply()
	cli.CLI.CLICommon.Debug = true
	cli.CLI.CLICommon.Debug.AfterApply() // ensure debugger is setuped
	cli.CLI.DBFilename = "test1.db"

	// the broker
	os.Remove(cli.CLI.DBFilename)
	server, _ := New()
	go server.Serve(
		Config{
			URL:      cli.CLI.URL,
			CertFile: cli.CLI.CertFile,
			KeyFile:  cli.CLI.KeyFile,
		},
	)

	// register an plugin for test-purpose
	var sendRetained sync.Mutex
	sendRetained.Lock()

	cli.CLI.DBFilename = "test2.db"
	os.Remove(cli.CLI.DBFilename)
	client, _ := New()

	clientPlugin := plugin.V1{
		OnConnack: func(clientID string) {
			client.Publish("/test/retained", []byte{0, 0}, true, 0)
			return
		},
		OnPublishRecvRequest: func(clientID string, topic string, payload string) bool {
			if "/test/retained" == topic {
				sendRetained.Unlock()
			}
			return true
		},
	}
	plugin.Register("clientPlugin", &clientPlugin)

	go client.Connect(
		Config{
			URL:      cli.CLI.URL,
			CertFile: cli.CLI.CertFile,
			KeyFile:  cli.CLI.KeyFile,
		},
		"",
		"",
	)
	sendRetained.Lock()
	for _, client := range server.clients {
		client.SendPingreq()
	}

}

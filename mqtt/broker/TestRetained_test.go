package broker

import (
	"sync"
	"testing"

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

func TestRetained(t *testing.T) {

	cli.CLI.URL = "tcp://127.0.0.1:1236"
	cli.CLI.CLICommon.Terminal = true
	cli.CLI.CLICommon.Terminal.AfterApply()
	cli.CLI.CLICommon.Debug = true
	cli.CLI.CLICommon.Debug.AfterApply() // ensure debugger is setuped
	cli.CLI.DBFilename = "test1.db"

	// the broker
	server, _ := New()
	go server.Serve(
		Config{
			URL:      cli.CLI.URL,
			CertFile: cli.CLI.CertFile,
			KeyFile:  cli.CLI.KeyFile,
		},
	)
	//time.Sleep(time.Millisecond * 500)

	// register an plugin for test-purpose
	var sendRetained sync.Mutex
	sendRetained.Lock()

	cli.CLI.DBFilename = "test2.db"
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

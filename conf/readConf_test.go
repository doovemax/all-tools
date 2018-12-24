package conf

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func Test_ReadConf(t *testing.T) {
	name := "./test.conf"
	test, _ := ReadConf(name)
	for _, host := range test {
		client, err := host.SshClient()
		if err != nil {
			logrus.Fatalln(err)

		}

		client.Stderr = os.Stderr
		client.Stdout = os.Stdout
		err = client.Run("/usr/bin/touch /tmp/hahaha")
		fmt.Println(err)
		client.Close()
		time.Sleep(time.Second * 20)
	}
}

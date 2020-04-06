package main

import (
	"bufio"
	"fmt"
	// "io/ioutil"
	"os"

	"github.com/aristanetworks/goeapi"
	// "github.com/aristanetworks/goeapi/module"
)

type Device struct {
	Transport string
	Hostname  string
	Username  string
	Password  string
	Port      int
}

type Output struct {
	text      string
	timestamp string
}

func RunCommand(cmd string, dut Device) (Output, error) {
	node, err := goeapi.Connect(dut.Transport, dut.Hostname, dut.Username, dut.Password, dut.Port)
	if err != nil {
		return Output{}, err
	}
	result, err := node.Enable([]string{"show version"})
	if err != nil {
		return Output{}, err
	}
	return Output{result[0]["result"], ""}, nil
}

func WriteFile(output Output, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)

	n4, err := w.WriteString(output.text)
	if err != nil {
		return err
	}
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()
	return nil
}

func main() {
	dut := Device{"https", "dmz-lf11", "fredlhsu", "arista", 443}
	// Take in list of devices + creds?
	devices := []Device{dut}
	// Take in command to run
	command := "show version"
	// connect to each device and run the command, save to file
	for _, device := range devices {
		output, err := RunCommand(command, device)
		if err != nil {
			fmt.Errorf("%s", err)
		}
		err = WriteFile(output, "output/show-version.output")
	}
}

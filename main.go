package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	//	"io"

	"log"
	"os"

	// "strings"

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

type DeviceList struct {
	Hosts []string
	Vars  Vars
}
type Vars struct {
	Username  string
	Password  string
	Transport string
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
	// Read in env vars for username/password
	fmt.Println("EAPI_USERNAME:", os.Getenv("EAPI_USERNAME"))
	fmt.Println("EAPI_PASSWORD:", os.Getenv("EAPI_PASSWORD"))
	// var dl DeviceList

	// reader := bufio.Reader{}
	// b := []byte{}
	// Check if there is something on stdin
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		// There is no piped input
		fmt.Println("no pipe :(")
		viper.SetConfigName("devices")
		viper.AddConfigPath(".")
		// TODO add err checking for below
		viper.ReadInConfig()
	} else {
		// There is piped input
		log.Println("hi pipe!")
		//reader = *bufio.NewReader(os.Stdin)

		// reader, _ := ioutil.ReadAll(os.Stdin)
		viper.SetConfigType("json")
		viper.ReadConfig(os.Stdin)
		//_, err := reader.Read(b)
		//if err != nil {
		//		fmt.Println(err)
		//}
		// fmt.Printf("%+v", dl)
		// err = json.Unmarshal(reader, &dl)
		if err != nil {

			log.Fatalf("badly formed JSON input %s", err)
		}
	}

	fmt.Println(viper.Get("hosts"))
	// Load file if nothing on stdin

	// TODO Read in and parse json file of devices
	directory := "output"
	dut := Device{"https", "dmz-lf11", "fredlhsu", "arista", 443}
	// Take in list of devices + creds?
	devices := []Device{dut}
	// Take in command to run
	command := "show version"
	// connect to each device and run the command, save to file
	for _, device := range devices {
		output, err := RunCommand(command, device)
		if err != nil {
			log.Printf("%s\n", err)
		}
		filename := fmt.Sprintf("%s--%s", device.Hostname, strings.ReplaceAll(command, " ", "_"))
		// TODO Add trailing / if not present in directory
		err = WriteFile(output, directory+"/"+filename)
	}
}

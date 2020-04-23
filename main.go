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

// Device that will be contacted by eapi
type Device struct {
	Transport string
	Hostname  string
	Username  string
	Password  string
	Port      int
}

// DeviceList is a list of devices to connect to that is parsed from the Ansible based JSON file format
type DeviceList struct {
	Hosts []string
	Vars  Vars
}

// Vars holds the login information for a device, based on Ansible
type Vars struct {
	Username  string
	Password  string
	Transport string
	Port      int
}

// Output holds the data to be written to a file after sending the commands
type Output struct {
	device    string
	text      string
	timestamp string
}

// RunCommand takes a command to run, and a device to run it against.  It generates output of the command execution.
func RunCommand(cmd string, dut Device) (Output, error) {
	node, err := goeapi.Connect(dut.Transport, dut.Hostname, dut.Username, dut.Password, dut.Port)
	if err != nil {
		return Output{}, err
	}
	result, err := node.Enable([]string{"show version"})
	if err != nil {
		return Output{}, err
	}
	// TODO Add timestamp
	return Output{dut.Hostname, result[0]["result"], ""}, nil
}

// WriteFile saves the command output to a path specified at `path`
func WriteFile(output Output, path string) error {
	// TODO Add devicename and timestamp to the filename
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
	// Read in env vars for username/password if present
	fmt.Println("EAPI_USERNAME:", os.Getenv("EAPI_USERNAME"))
	fmt.Println("EAPI_PASSWORD:", os.Getenv("EAPI_PASSWORD"))
	// var dl DeviceList

	// TODO: Options for output, if writing to file send path/filename to stdout
	
	// reader := bufio.Reader{}
	// b := []byte{}
	// Check if there is something on stdin - this is useful if piping a list of devices into eoscmd
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	// TODO add cli flag to choose if pulling from device list, stdin, or specified on CLI

	// Viper is used below to parse a device list to execute the command against
	if fi.Mode()&os.ModeNamedPipe == 0 {
		// There is no piped input, check for a file with a list of devices
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

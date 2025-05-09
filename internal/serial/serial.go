package serial

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/tarm/serial"
)

func OpenSerialPort() (*serial.Port, error) {
	config := &serial.Config{
		Name: "COM9",
		Baud: 9600,
	}
	return serial.OpenPort(config)
}

func ListenToSerial(port *serial.Port, callback func(string, int)) {
	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			message := parts[0]
			increment := 0
			fmt.Sscanf(parts[1], "%d", &increment)
			callback(message, increment)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading from serial port:", err)
	}
}

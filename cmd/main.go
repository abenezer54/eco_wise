package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"github.com/sfreiberg/gotwilio"
	"github.com/tarm/serial"
)

var (
	mu          sync.Mutex
	bottleCount int
	serialPort  *serial.Port
	useSerial   bool
)

func main() {
	flag.BoolVar(&useSerial, "useSerial", true, "Enable or disable serial communication")
	flag.Parse()
	// Start TCP server
	go startTCPServer()

	// Set up HTTP server (Gin API routes)
	go func() {
		r := gin.Default()
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal("Failed to set trusted proxies: ", err)
		}
		corsHandler := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Origin", "Content-Type", "Accept"},
			ExposedHeaders:   []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           int(12 * time.Hour / time.Second),
		})

		r.Use(func(c *gin.Context) {
			corsHandler.HandlerFunc(c.Writer, c.Request)
			c.Next()
		})

		r.POST("/payment", handlePayment)
		r.Static("/assets", "./frontend/dist/assets")
		r.StaticFile("/", "./frontend/dist/index.html")

		r.NoRoute(func(c *gin.Context) {
			c.File("./frontend/dist/index.html")
		})

		r.Run(":8080")
	}()

	if useSerial {
		// Set up Serial communication with Arduino
		var err error
		serialPort, err = openSerialPort()
		if err != nil {
			log.Fatal("Failed to open serial port: ", err)
		}
		defer serialPort.Close()

		// Start listening to the serial port in a separate goroutine
		go listenToSerial(serialPort, handleSerialMessage)
	}
	// Run the server indefinitely
	select {}
}

func startTCPServer() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Failed to start TCP server: ", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection: ", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Received TCP message: %s", line)
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			message := parts[0]
			increment := 0
			fmt.Sscanf(parts[1], "%d", &increment)
			handleSerialMessage(message, increment)

			// Send message to Arduino
			if useSerial {
				_, err := serialPort.Write([]byte(message + "\n"))
				if err != nil {
					log.Println("Error sending message to Arduino:", err)
				} else {
					log.Printf("Sent message to Arduino: %s", message)
				}
			}
		} else {
			log.Printf("Invalid message format: %s", line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading from TCP connection:", err)
	}
}

func openSerialPort() (*serial.Port, error) {
	config := &serial.Config{
		Name: "COM9", // Replace with your serial port name
		Baud: 9600,
	}
	return serial.OpenPort(config)
}

func listenToSerial(port *serial.Port, callback func(string, int)) {
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

func handleSerialMessage(message string, increment int) {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("Received serial message: %s, increment: %d", message, increment)

	if message == "accept" {
		bottleCount += increment
		log.Printf("Bottle count incremented to: %d", bottleCount)
	} else {
		log.Printf("Received unrecognized message: %s", message)
	}
}

func handlePayment(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate phone number
	phoneRegex := regexp.MustCompile(`^\+251\d{9}$`)
	if !phoneRegex.MatchString(request.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}

	mu.Lock()
	amount := bottleCount * 1
	bottleCount = 0
	mu.Unlock()

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)
	from := "+12317517378"
	to := request.PhoneNumber

	var message string
	if amount > 0 {
		message = "\nኢኮ-ዋይዝ \n እንኳን ደስ አለዎት እና " + strconv.Itoa(amount) + " ብር ተቀበሉ። \nእናንተን በኢኮ-ዋይዝ ስለተጠቀሙ እናመሰግናለን!"
	} else {
		message = "\nኢኮ-ዋይዝ \n የተሰበሰበው ብዛት የበቃ አይደለም። እባክዎን በድጋሚ ይሞክሩ።"
	}

	_, _, err := twilio.SendSMS(from, to, message, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful and SMS sent"})
}

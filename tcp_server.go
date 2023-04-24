package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/kristianvv/is105sem03/mycrypt"
	"github.com/kristianvv/minyr/yr"
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.2:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("mottatt tilkobling fra %s", conn.RemoteAddr().String())

			wg.Add(1)
			go func(conn net.Conn) {
				defer wg.Done()
				defer conn.Close()

				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Printf("feil ved mottak: %v", err)
						}
						break
					}

					command := string(buf[:n])
					log.Printf("mottatt kommando '%s'", command)

					var response string

					switch command {
					case "conv":
						yr.ConvTemperature()
						response = "Temperature conversion complete"
					case "avg":
						avgTemp := yr.AverageTemp()
						response = fmt.Sprintf("Average temperature for the period is: %.2f°C\n", avgTemp)
					case "crypt":
						response = "Enter message to encrypt:"
						conn.Write([]byte(response))

						n, err := conn.Read(buf)
						if err != nil {
							log.Printf("feil ved mottak: %v", err)
							break
						}

						message := string(buf[:n])
						log.Printf("mottatt melding '%s'", message)

						encryptedMessage := string(mycrypt.Krypter([]rune(message), mycrypt.ALF_SEM03, 5))
						response = fmt.Sprintf("Encrypted message: %s", encryptedMessage)
					default:
						response = "Invalid command"
					}

					conn.Write([]byte(response))
				}
			}(conn)
		}
	}()

	wg.Wait()
}

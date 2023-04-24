package main

import (
	"io"
	"log"
	"net"
	"sync"

	"github.com/kristianvv/is105sem03/mycrypt"
        "github.com/kristianvv/funtemps/conv"
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
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					decrypted := mycrypt.Decrypt(buf[:n])
					switch msg := string(decrypted); msg {
					case "ping":
						encrypted := mycrypt.Encrypt([]byte("pong"))
						_, err = c.Write(encrypted)
					case string(mycrypt.Decrypt([]byte("Kjevik"))):
						temperature := yr.GetTemperature()
						converted := conv.ConvertTemperature(temperature, "C", "F")
						encrypted := mycrypt.Encrypt([]byte(converted))
						_, err = c.Write(encrypted)
					default:
						encrypted := mycrypt.Encrypt(decrypted)
						_, err = c.Write(encrypted)
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()

	wg.Wait()
}

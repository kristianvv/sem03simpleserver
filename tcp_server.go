package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/kristianvv/funtemps/conv"
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.4:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	log.Printf("bundet til %s", server.Addr().String())

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			log.Println("før server.Accept() kallet")

			conn, err := server.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			wg.Add(1)
			go func(c net.Conn) {
				defer wg.Done()
				defer conn.Close()

				for {
					buf := make([]byte, 2048)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return
					}

					dekryptertMelding := Krypter([]rune(string(buf[:n])), ALF_SEM03, len(ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))

					msgString := string(dekryptertMelding)

					switch msgString {
					case "ping":
						kryptertMelding := Krypter([]rune("pong"), ALF_SEM03, 4)
						log.Println("Kryptert melding: ", string(kryptertMelding))
						_, err = c.Write([]byte(string(kryptertMelding)))

					default:
						if strings.HasPrefix(msgString, "Kjevik") {
							newString, err := CelsiusToFarenheitLine("Kjevik;SN39040;18.03.2022 01:50;6")
							if err != nil {
								log.Fatal(err)
							}

							kryptertMelding := Krypter([]rune(newString), ALF_SEM03, 4)
							_, err = conn.Write([]byte(string(kryptertMelding)))
						} else {
							kryptertMelding := Krypter([]rune(string(buf[:n])), ALF_SEM03, 4)
							_, err = c.Write([]byte(string(kryptertMelding)))

						}
					}

					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return
					}
				}
			}(conn)
		}
	}()

	wg.Wait()
}

func CelsiusToFarenheitLine(line string) (string, error) {

	dividedString := strings.Split(line, ";")
	var err error

	if len(dividedString) == 4 {
		dividedString[3], err = CelsiusToFarenheitString(dividedString[3])
		if err != nil {
			return "", err
		}
	} else {
		return "", errors.New("linje har ikke forventet format")
	}
	return strings.Join(dividedString, ";"), nil
}

func CelsiusToFarenheitString(celsius string) (string, error) {
	var fahrFloat float64
	var err error
	if celsiusFloat, err := strconv.ParseFloat(celsius, 64); err == nil {
		fahrFloat = conv.CelsiusToFahrenheit(celsiusFloat)
	}
	fahrString := fmt.Sprintf("%.1f", fahrFloat)
	return fahrString, err
}

var ALF_SEM03 []rune = []rune("abcdefghijklmnopqrstuvwxyzæøå0123456789.,:; KSN") //ABCDEFGHIJKLMNOPQRSTUVWXYZÆØÅ

func Krypter(melding []rune, alphabet []rune, chiffer int) []rune {
	kryptertMelding := make([]rune, len(melding))
	for i := 0; i < len(melding); i++ {
		indeks := sokIAlfabetet(melding[i], alphabet)
		if indeks+chiffer >= len(alphabet) {
			kryptertMelding[i] = alphabet[indeks+chiffer-len(alphabet)]
		} else {
			kryptertMelding[i] = alphabet[indeks+chiffer]
		}

	}
	return kryptertMelding
}

func sokIAlfabetet(symbol rune, alfabet []rune) int {
	for i := 0; i < len(alfabet); i++ {
		if symbol == alfabet[i] {
			return i
		}
	}
	return -1
}

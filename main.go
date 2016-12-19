package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type AppConfig struct {
	Port int
}

var cfg AppConfig

func init() {
	cfg = AppConfig{}
	flag.IntVar(&cfg.Port, "port", 8000, "-port NUMBER")
	flag.Parse()
}

func main() {
	log.Printf(time.Now().String())
	log.Println("Starting on port ", cfg.Port)
	addr := "localhost:" + strconv.Itoa(cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Cannot start listener: ", err)
		return
	}

	defer listener.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	go acceptReq(listener)
	select {}
}

func acceptReq(listener net.Listener) {
	defer listener.Close()

	for {
		log.Println("Waiting for connections...")
		con, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Client acceped")

		go handleMessages(con)
	}
}

func handleMessages(con net.Conn) {
	defer con.Close()

	for {
		reader := bufio.NewReader(con)
		command, err := reader.ReadString('\n')

		if err != nil {
			switch t := err.(type) {
			default:
				log.Println(err)
			case *net.OpError:
				log.Println("Connection closed")
				if t.Temporary() {
					log.Println("temporary")
				}
			}

			break
		}

		loc, err := time.LoadLocation(command)
		if err != nil {
			log.Println(err)
			con.Write([]byte(time.Now().UTC().String()))

		} else {
			readTime := time.Now().In(loc).String()
			con.Write([]byte(readTime))
		}

		log.Printf("Command: %s", command)
	}
}

package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	"gopkg.in/sevlyar/go-daemon.v0"
)

const appname = "daemon-producer"

var (
	signal = flag.String("s", "", `Send signal to the daemon:
          quit — graceful shutdown
          stop — fast shutdown
          reload — reloading the configuration file`)

	queueUrl string
)

func main() {

	flag.Parse()

	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	cntxt := &daemon.Context{
		PidFileName: appname + ".pid",
		PidFilePerm: 0644,
		LogFileName: appname + ".log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[" + appname + "]"},
	}

    queueUrl = os.Getenv("QUEUE_URL")
    if queueUrl == "" {
        log.Fatalln("QUEUE_URL must be specified");
    }


	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	go worker()

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker() {
LOOP:
	for {
		time.Sleep(2 * time.Second) // this is work to be done by worker.
		select {
		case <-stop:
			break LOOP
		default:
            if err := Send(queueUrl); err != nil {
                panic(err.Error())
            }
		}
	}
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("configuration reloaded")
	return nil
}
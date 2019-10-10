package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elwinar/rcoredump"
)

func main() {
	var cfg struct {
		dest string
		src  string
		log  string
	}
	flag.StringVar(&cfg.dest, "dest", "localhost:1105", "address of the destination host")
	flag.StringVar(&cfg.src, "src", "-", "path of the coredump to send to the host ('-' for stdin)")
	flag.StringVar(&cfg.log, "log", "/var/log/rcoredump.log", "path of the log file for rcoredump")
	flag.Parse()

	// Open the log file.
	l, err := os.OpenFile(cfg.log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer l.Close()
	log.SetOutput(l)

	// Gather a few variables.
	// Args from the command line should be, in order:
	// - %E, pathname of executable
	// - %t, time of dump
	// - %h, hostname
	args := flag.Args()
	if len(args) != 3 {
		log.Println("unexpected number of arguments on command-line")
		return
	}

	// Pathname of the executable comes up with ! instead of /.
	executable := strings.Replace(args[0], "!", "/", -1)
	timestamp, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		log.Println("invalid timestamp format")
		return
	}
	hostname := args[2]

	// Open the connection to the backend.
	conn, err := net.Dial("tcp", cfg.dest)
	if err != nil {
		log.Println("dialing tcp:", err)
		return
	}
	defer conn.Close()

	// Compress the data before sending it.
	w := gzip.NewWriter(conn)
	defer w.Close()

	// Send the header first.
	err = json.NewEncoder(w).Encode(rcoredump.Header{
		Executable: executable,
		Date:       time.Unix(timestamp, 0),
		Hostname:   hostname,
	})
	if err != nil {
		log.Println("writing header:", err)
		return
	}
	err = w.Close()
	if err != nil {
		log.Println("closing header stream:", err)
		return
	}
	w.Reset(conn)

	// Then the core itself.
	var in io.ReadCloser
	if cfg.src == "-" {
		in = os.Stdin
	} else {
		in, err = os.Open(cfg.src)
		if err != nil {
			log.Println("opening core file:", err)
			return
		}
		defer in.Close()
	}
	_, err = io.Copy(w, in)
	if err != nil {
		log.Println("writing core:", err)
		return
	}
	err = w.Close()
	if err != nil {
		log.Println("closing core stream:", err)
		return
	}
	w.Reset(conn)

	// Then the binary.
	bin, err := os.Open(executable)
	if err != nil {
		log.Println("opening bin file:", err)
		return
	}
	defer in.Close()
	_, err = io.Copy(w, bin)
	if err != nil {
		log.Println("writing bin:", err)
		return
	}
	err = w.Close()
	if err != nil {
		log.Println("closing bin stream:", err)
		return
	}

	// All done, k-thx-bye.
}
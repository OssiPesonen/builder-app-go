package utils

import (
	"bytes"
	"io"
	"log"
	"strings"
)

type Logstreamer struct {
	Logger *log.Logger
	buf    *bytes.Buffer
	// "stdout", "stderr"
	prefix string
	// if true, saves output in memory
	record  bool
	persist string
}

func NewLogstreamer(logger *log.Logger, prefix string, record bool) *Logstreamer {
	streamer := &Logstreamer{
		Logger:  logger,
		buf:     bytes.NewBuffer([]byte("")),
		prefix:  prefix,
		record:  record,
		persist: "",
	}

	return streamer
}

func (l *Logstreamer) Write(p []byte) (n int, err error) {
	if n, err = l.buf.Write(p); err != nil {
		return
	}

	err = l.OutputLines()
	return
}

func (l *Logstreamer) Close() error {
	if err := l.Flush(); err != nil {
		return err
	}
	l.buf = bytes.NewBuffer([]byte(""))
	return nil
}

func (l *Logstreamer) Flush() error {
	p := make([]byte, l.buf.Len())
	if _, err := l.buf.Read(p); err != nil {
		return err
	}

	l.out(string(p))
	return nil
}

func (l *Logstreamer) OutputLines() error {
	for {
		line, err := l.buf.ReadString('\n')

		if len(line) > 0 {
			if strings.HasSuffix(line, "\n") {
				l.out(line)
			} else {
				// put back into buffer, it's not a complete line yet
				//  Close() or Flush() have to be used to flush out
				//  the last remaining line if it does not end with a newline
				if _, err := l.buf.WriteString(line); err != nil {
					return err
				}
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Logstreamer) FlushRecord() string {
	buffer := l.persist
	l.persist = ""
	return buffer
}

func (l *Logstreamer) out(str string) {
	if len(str) < 1 {
		return
	}

	if l.record == true {
		l.persist = l.persist + str
	}

	if l.prefix == "stdout" {
		str = "-- DATA: " + str
	} else if l.prefix == "stderr" {
		str = "-- ERROR: " + str
	} else {
		str = "-- INFO: " + str
	}

	l.Logger.Print(str)
}

package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	config := Config{
		SampleRate:       44100,
		NumChannels:      2,
		BitRate:          16,
		ByteOrder:        binary.LittleEndian,
		CompLevel:        8,
		SilenceThreshold: 30 * 44100 / 1000,
	}
	err := runMain(config, os.Stdin)
	if err != nil {
		panic(err)
	}
}

func runMain(cfg Config, stdin io.Reader) error {
	sessId := time.Now().Unix()
	streamId := 0
	var filename string
	for {
		streamId++
		filename = fmt.Sprintf("stream_%d-%d.flac", sessId, streamId)
		fmt.Println(filename)
		err := processPcmStream(cfg, stdin, filename)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("processing stdin: %w", err)
		}
	}
	// remove the last stream since it's garbage (either silence or not finished)
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("performing cleanup: %w", err)
	}
	return nil
}

func processPcmStream(cfg Config, stdin io.Reader, filename string) error {
	cmd := cfg.FlacCmd(filename)
	output, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("obtaining cmd stdin pipe: %w", err)
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("starting cmd: %w", err)
	}
	defer withErrLog(cmd.Wait, "waiting on cmd")
	defer withErrLog(output.Close, "closing cmd stdin")
	err = copyPcmSubstream(cfg, stdin, output)
	if err == io.EOF {
		return err
	}
	if err != nil {
		return fmt.Errorf("processing substream: %w", err)
	}
	return nil
}

func copyPcmSubstream(
	cfg Config,
	stdin io.Reader,
	output io.Writer,
) error {
	streamStarted := false
	var buf []int16
	for {
		s1, s2, err := readTwoSamples(stdin, cfg.ByteOrder)
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			return fmt.Errorf("reading samples: %w", err)
		}
		buf = append(buf, s1, s2)
		if s1 == 0 && s2 == 0 {
			if len(buf) >= cfg.SilenceThreshold*cfg.NumChannels && streamStarted {
				return nil // stream ended
			}
		} else {
			if !streamStarted {
				fmt.Println("Stream started")
				buf = []int16{s1, s2}
				streamStarted = true
			}
			err = binary.Write(output, cfg.ByteOrder, buf[:])
			if err != nil {
				return fmt.Errorf("writing to output: %w", err)
			}
			buf = nil
		}
	}
}

func readTwoSamples(r io.Reader, byteOrder binary.ByteOrder) (int16, int16, error) {
	var s1, s2 int16
	err := binary.Read(r, byteOrder, &s1)
	if err != nil {
		return s1, s2, err
	}
	err = binary.Read(r, byteOrder, &s2)
	if err != nil {
		return s1, s2, err
	}
	return s1, s2, nil
}

func withErrLog(f func() error, msg string) {
	err := f()
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", msg, err)
	} else {
		fmt.Printf("OK: %s\n", msg)
	}
}

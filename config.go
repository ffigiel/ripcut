package main

import (
	"encoding/binary"
	"os/exec"
	"strconv"
)

type Config struct {
	SampleRate       int
	NumChannels      int
	BitRate          int
	ByteOrder        binary.ByteOrder
	CompLevel        int
	SilenceThreshold int // number of samples
}

func (c Config) FlacCmd(outputName string) *exec.Cmd {
	endianStr := "little"
	if c.ByteOrder != binary.LittleEndian {
		endianStr = "big"
	}
	return exec.Command(
		"flac",
		"-",
		intFlag("-", c.CompLevel),
		"--output-name="+outputName,
		"--endian="+endianStr,
		"--sign=signed",
		intFlag("--channels=", c.NumChannels),
		intFlag("--sample-rate=", c.SampleRate),
		intFlag("--bps=", c.BitRate),
		"--force-raw-format",
	)
}

func intFlag(flag string, value int) string {
	return flag + strconv.Itoa(value)
}

// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ocp implements OCP/Facebook-specific IPMI client functions.
package ocp

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// Set RW_VPD key-value
func Set(key string, value []byte) error {
	file, err := ioutil.TempFile("/tmp", "rwvpd*.bin")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	cmd := exec.Command("flashrom", "-p", "internal:ich_spi_mode=hwseq", "-c", "Opaque flash chip", "--fmap", "-i", "RW_VPD", "-r", file.Name())
	cmd.Stdin, cmd.Stdout = os.Stdin, os.Stdout
	if err = cmd.Run(); err != nil {
		log.Printf("flashrom failed to read RW_VPD: %v", err)
		return err
	}
	cmd = exec.Command("vpd", "-f", file.Name(), "-s", key+"="+string(value[:]))
	if err = cmd.Run(); err != nil {
		log.Printf("vpd failed to set key: %v value: %v, err: %v", key, string(value[:]), err)
		return err
	}
	cmd = exec.Command("flashrom", "-p", "internal:ich_spi_mode=hwseq", "-c", "Opaque flash chip", "--fmap", "-i", "RW_VPD", "--noverify-all", "-w", file.Name())
	if err = cmd.Run(); err != nil {
		log.Printf("flashrom failed to write RW_VPD: %v", err)
		return err
	}
	return nil
}

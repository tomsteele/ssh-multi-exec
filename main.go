package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ok(err error, location string) {
	fmt.Fprintf(os.Stderr, "There was an error at %s\n", location)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err.Error())
		os.Exit(1)
	}
}

func grabLines(scanner *bufio.Scanner, count int, joinWith string) string {
	lines := []string{}
	scanner.Scan()

	for i := 0; i < count; i++ {
		if scanner.Scan() {
			lines = append(lines, scanner.Text())
		} else {
			break
		}
	}

	return strings.Join(lines, joinWith)
}

func main() {
	sshFile := flag.String("ssh-server-file", "", "File containing newline separated server list. Format user@host:port. Example: root@192.168.1.1:22")
	keyFile := flag.String("key", "", "Path to SSH private key file.")
	inputFile := flag.String("input", "", "File containing chunks of newline separated inputs. For nmap/rumble this would be targets.")
	chunkSize := flag.Int("chunk-size", 1, "How many lines to combine per run.")
	joinWith := flag.String("join-with", " ", "String used to join chunks. Default is a space, and probably what you want.")
	command := flag.String("c", "rumble --probes connect --text --output-raw - -o disable --nowait %s", "Format string to pass to SSH.")
	dryRun := flag.Bool("dry-run", false, "Perform a dry run. Output only the commands to execute.")
	flag.Parse()
	inputFH, err := os.Open(*inputFile)
	ok(err, "opening input file")

	sshServerList := []string{}
	sshFH, err := os.Open(*sshFile)
	ok(err, "opening ssh-server-file")
	sshScanner := bufio.NewScanner(sshFH)
	for sshScanner.Scan() {
		sshServerList = append(sshServerList, sshScanner.Text())
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	scanner := bufio.NewScanner(inputFH)
	for {
		hosts := grabLines(scanner, *chunkSize, *joinWith)
		if hosts == "" {
			break
		}
		sshServer := sshServerList[r.Intn(len(sshServerList))]
		commandToExecute := fmt.Sprintf(*command, hosts)
		fmt.Fprintf(os.Stderr, "\n[+] Using ssh server %s to execute: %s\n", sshServer, commandToExecute)
		if !*dryRun {
			cmd := exec.Command("ssh", "-i", *keyFile, sshServer, "-o", "StrictHostKeyChecking=no", commandToExecute)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			ok(cmd.Start(), "starting ssh command")
			ok(cmd.Wait(), "waiting for ssh command to complete")
		}
	}
}

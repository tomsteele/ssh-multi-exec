# ssh-multi-exec
Execute tasks across SSH hosts using random selection. Can be used to execute tools such as Rumble or Nmap across many SSH hosts.

## Install
No binaries available (yet). If you have a Go environment configured:
```
go get github.com/tomsteele/ssh-multi-exec
```

## Usage
```
Usage of ./ssh-multi-exec:
  -c string
        Format string to pass to SSH. (default "rumble --probes connect --text --output-raw - -o disable --nowait %s")
  -chunk-size int
        How many lines to combine per run. (default 1)
  -dry-run
        Perform a dry run. Output only the commands to execute.
  -input string
        File containing chunks of newline separated inputs. For nmap/rumble this would be targets.
  -join-with string
        String used to join chunks. Default is a space, and probably what you want. (default " ")
  -key string
        Path to SSH private key file.
  -ssh-server-file string
        File containing newline separated server list. Format user@host:port. Example: root@192.168.1.1:22
```

Example Run. This will scan 5 hosts at a time using a random SSH server from a file provided. Can tee the results into a single file for later review.
```
# ./ssh-multi-exec -chunk-size 5 \
-input ips.txt \
-key ~/.ssh/id_rsa \
-ssh-server-file ssh_list.txt \
-c "rumble -S false --probes connect --text --output-raw - -o disable --nowait -p 80,443,8443,8080 %s" | tee -a rumble_web.jsonl
```

Rumble makes it easy to generate a scan from this file with `--import`:
```
# rumble --import rumble_web.jsonl --text
```

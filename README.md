# ssh-multi-exec
Execute tasks across SSH hosts using random selection. Can be used to execute tools such as Rumble or Nmap across many SSH hosts.

## Usage
```
# ./ssh-multi-exec -chunk-size 5 \
-input ips.txt \
-key ~/.ssh/id_rsa \
-ssh-server-file ssh_list.txt \
-c "rumble -S false --probes connect --text --output-raw - -o disable --nowait -p 80,443,8443,8080 %s" | tee -a rumble_web.jsonl
```

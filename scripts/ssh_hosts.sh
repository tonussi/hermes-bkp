#!/bin/sh

emulabuser=$1
known_hosts_file="ansible/known_hosts"

ssh-keyscan -f $known_hosts_file >> ~/.ssh/known_hosts

for line in $(cat $known_hosts_file); do
    ssh -p 22 $emulabuser@$line "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
done

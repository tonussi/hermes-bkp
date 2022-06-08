#!/bin/sh

# ssh-keyscan -f ansible/known_hosts >> ~/.ssh/known_hosts

gnome-terminal -e "ssh -p 22 LUCAS@pc443.emulab.net"
gnome-terminal -e "ssh -p 22 LUCAS@pc484.emulab.net"
gnome-terminal -e "ssh -p 22 LUCAS@pc515.emulab.net"
gnome-terminal -e "ssh -p 22 LUCAS@pc441.emulab.net"
gnome-terminal -e "ssh -p 22 LUCAS@pc445.emulab.net"
gnome-terminal -e "ssh -p 22 LUCAS@pc494.emulab.net"
gnome-terminal -e "ssh -p 22 LUCAS@pc447.emulab.net"
gnome-terminal -e "ssh -p 22 LUCAS@pc460.emulab.net"
gnome-terminal -e "ssh -p 22 LUCAS@pc444.emulab.net"
gnome-terminal -e "ssh -p 22 LUCAS@pc448.emulab.net"

ssh -p 22 LUCAS@pc443.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
ssh -p 22 LUCAS@pc484.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
ssh -p 22 LUCAS@pc515.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
ssh -p 22 LUCAS@pc441.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
ssh -p 22 LUCAS@pc445.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
ssh -p 22 LUCAS@pc494.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
ssh -p 22 LUCAS@pc447.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
ssh -p 22 LUCAS@pc460.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
ssh -p 22 LUCAS@pc444.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"
ssh -p 22 LUCAS@pc448.emulab.net "ifconfig | tr '\n' '~' | sed 's/~~/\n/g' | sed 's/~/ /g' | sed -E 's/^\s*(\S+)\s+.*\s+inet\s+(\S+)\s.*$/\1 \2/g'"

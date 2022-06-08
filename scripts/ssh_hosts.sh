#!/bin/sh

ssh-keyscan -f ansible/known_hosts >> ~/.ssh/known_hosts

gnome-terminal -e "ssh -p 22 pc443.emulab.net"
gnome-terminal -e "ssh -p 22 pc484.emulab.net"
gnome-terminal -e "ssh -p 22 pc515.emulab.net"
gnome-terminal -e "ssh -p 22 pc441.emulab.net"
gnome-terminal -e "ssh -p 22 pc445.emulab.net"
gnome-terminal -e "ssh -p 22 pc494.emulab.net"
gnome-terminal -e "ssh -p 22 pc447.emulab.net"
gnome-terminal -e "ssh -p 22 pc460.emulab.net"
gnome-terminal -e "ssh -p 22 pc444.emulab.net"
gnome-terminal -e "ssh -p 22 pc448.emulab.net"

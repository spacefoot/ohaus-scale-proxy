#!/bin/sh
IP="$(zenity --entry --title 'OHAUS scale proxy' --text 'Entrez la nouvelle adresse IP de la balance')"
echo "address: $IP" > ~/.config/ohaus-scale-proxy.yml
systemctl --user restart ohaus-scale-proxy

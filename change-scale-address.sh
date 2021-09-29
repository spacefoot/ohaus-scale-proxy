#!/bin/sh

set -x

systemctl --user status ohaus-scale-proxy > /dev/null
RUNNING=$?

TITLE='OHAUS scale proxy'
if [ "$RUNNING" = 0 ]; then
    TITLE="$TITLE (running)"
fi

LAST_IP="$(grep 'address:' ~/.config/ohaus-scale-proxy.yml | awk '{print $2}')"
DATA="$(zenity --entry --title "$TITLE" --text 'Entrez la nouvelle adresse IP de la balance' --entry-text=$LAST_IP --extra-button=Arrêter --ok-label=Démarrer --cancel-label=Annuler)"
STATUS=$?

if [ "$STATUS" = 0 ]; then
    echo "address: $DATA" > ~/.config/ohaus-scale-proxy.yml
    systemctl --user restart ohaus-scale-proxy
elif [ "$DATA" = Arrêter ]; then
    systemctl --user stop ohaus-scale-proxy
fi


#!/bin/bash

if [ $(docker ps -q -f name=selenium) ]; then
    echo "Le conteneur Selenium est déjà en cours d'exécution."
else
    echo "Démarrage du conteneur Selenium..."
    docker run -d -p 4444:4444 -p 7900:7900 --shm-size="2g" --name selenium --platform=linux/amd64 selenium/standalone-chrome
fi
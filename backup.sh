#!/bin/bash
REPO_PATH="/home/evgeniyfimushkin/git/minecraft-server"

cd "$REPO_PATH" || exit

git add .

git commit -m "Auto commit: $(date '+%Y-%m-%d %H:%M:%S')"

git push origin master

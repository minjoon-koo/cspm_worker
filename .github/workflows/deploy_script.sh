#!/bin/bash

cd /app
git pull origin main

cd /app/worker

# Install steampipe
if [ ! -d "steampipe" ]; then
  git clone https://github.com/turbot/steampipe
fi
cd steampipe/
sudo make

# Install azure plugin
steampipe plugin install azure
steampipe plugin install azuread

# Install powerpipe and azure_mod
sudo /bin/sh -c "$(curl -fsSL https://powerpipe.io/install/powerpipe.sh)"
mkdir -p dashboards
cd dashboards
powerpipe mod init
powerpipe mod install github.com/turbot/steampipe-mod-azure-insights
powerpipe mod install github.com/turbot/steampipe-mod-azure-compliance

cd /app/worker
go build -o main /app/worker/main.go
nohup ./main > /dev/null 2>&1 &
disown

echo "Deployment script executed successfully"

#install steampipe
git clone https://github.com/turbot/steampipe
cd steampipe/
sudo make

#install azure plugin
steampipe plugin install azure
steampipe plugin install azuread

#install powerpipe, azure_mod
sudo /bin/sh -c "$(curl -fsSL https://powerpipe.io/install/powerpipe.sh)"
mkdir dashboards
cd dashboards
powerpipe mod init
powerpipe mod install github.com/turbot/steampipe-mod-azure-insights
powerpipe mod install github.com/turbot/steampipe-mod-azure-compliance

chmod +x /path/to/your/app/deploy_script.sh
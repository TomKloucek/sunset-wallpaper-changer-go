#! /bin/bash
#write out current crontab
crontab -l > mycron
#echo new cron into cron file
current_dir=$(pwd)
echo "@reboot $current_dir/login.sh" >> mycron
echo "@reboot sleep 60 && $current_dir/login.sh" >> mycron
#install new cron file
crontab mycron
rm mycron
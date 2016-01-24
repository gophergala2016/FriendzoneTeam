#!/bin/bash
sudo apt-get install mysql-server php5-mysql
sudo mysql_install_db
sudo apt-get -y install nginx php5-fpm
sudo systemctl restart php5-fpm
sudo systemctl restart mysql
sudo systemctl restart nginx

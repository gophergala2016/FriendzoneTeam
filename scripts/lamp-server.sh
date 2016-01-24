#!/bin/bash
sudo apt-get -y install apache2 apache2-doc mysql-server mysql-client php5-mysql libapache2-mod-php5 perl libapache2-mod-perl2
sudo mysql_install_db
sudo systemctl restart apache2
sudo systemctl restart mysql
sudo a2enmod userdir

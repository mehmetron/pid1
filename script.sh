#!/bin/bash

apt-get -y update
apt-get -y install golang
apt-get -y install nodejs
apt-get -y install python3
#apt-get -y install nano
#apt-get -y install lsof
#apt-get -y install curl
#apt-get -y install net-tools

# Node
#curl "https://nodejs.org/dist/v16.3.0/node-v16.3.0-linux-x64.tar.xz" -o node.tar.xz
#tar xf node.tar.xz --strip-components=1
#rm node.tar.xz

# Go
#curl -LO https://golang.org/dl/go1.16.2.linux-amd64.tar.gz
#tar -xzf go1.16.2.linux-amd64.tar.gz
#rm go1.16.2.linux-amd64.tar.gz

## Python
#curl "https://www.python.org/ftp/python/3.9.4/Python-3.9.4.tgz" -o python.tar.gz
#tar xzf python.tar.gz --strip-components=1
#rm python.tar.gz
#./configure


#wget https://www.python.org/ftp/python/3.8.3/Python-3.8.3.tgz
#tar -xf Python-3.8.3.tgz



#---------------------------Testing languages---------------------------
cd home
touch main.js
echo "console.log('OK javascript');" > main.js

cd home
touch main.go
cat << EOF > main.go
package main

import "fmt"

func main() {
	fmt.Println("OK go")
}
EOF

cd home
touch main.py
echo "print('OK python')" > main.py



node main.js
go run main.go
python3 main.py
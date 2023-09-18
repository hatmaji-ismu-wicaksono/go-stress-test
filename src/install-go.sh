wget https://go.dev/dl/go1.18.10.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
sudo cp ./golang.sh /etc/profile.d/
export PATH=$PATH:/usr/local/go/bin
go version
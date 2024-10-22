cd /root/home/nightflavor/httpserver
git reset --hard
git pull
go build -o /root/home/nightflavor/httpserver/main.go
sudo /bin/systemctl restart lukasweb.service


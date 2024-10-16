cd ~/httpserver
systemctl stop lukasweb.service
rm pull.sh
git pull
chmod +x pull.sh
go build main.go
systemctl start lukasweb.service

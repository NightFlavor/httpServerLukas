echo "bitchass"
cd ~/httpserver
sudo systemctl stop lukasweb.service
rm pull.sh
git pull
chmod +x pull.sh
go build main.go
sudo systemctl start lukasweb.service

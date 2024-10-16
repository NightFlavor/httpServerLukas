cd ~/httpserver
systemctl stop httpweb.service
git pull
go build main.go
systemctl start httpweb.service
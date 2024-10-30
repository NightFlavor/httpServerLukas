git reset --hard
git pull

go build -o main main.go

chmod +x pull.sh

sudo /bin/systemctl restart lukasweb.service

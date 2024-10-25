git reset --hard
git pull

go build -o main main.go

sudo /bin/systemctl restart lukasweb.service

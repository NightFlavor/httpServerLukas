cd ~/httpserver
#!/bin/bash
git pull
go build -o /path/to/output
sudo /bin/systemctl restart lukasweb.service


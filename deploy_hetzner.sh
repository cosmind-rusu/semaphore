export GOOS=linux
export GOARCH=amd64

task build

scp ./bin/semaphore hetzner:~/

#ssh hetzner << EOF
#sudo systemctl stop demo2.semaphore
#sudo cp ~/semaphore /usr/bin/semaphore
#sudo systemctl start demo2.semaphore
#EOF
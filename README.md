# monitor-agent
An agent used to discover processes and monitor their resource utilization

Steps to build:
1. Clone this repo
2. Install go on your system
3. Build for your architecture

Win: GOOS=windows GOARCH=amd64 go build -o monitoragent.exe
Linux: GOOS=linux GOARCH=amd64 go build -o monitoragent

4. Run the executable as root/administrator
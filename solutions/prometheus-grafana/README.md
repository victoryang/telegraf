# Solution

## Download repo
    git clone git@gitlab.elite-robot.org:RobotControllerMonitor/solution-prometheus-grafana.git

## Install docker-compose
    sudo apt install -y docker-compose

## Configure
    vi prometheus.yml
    modify target to telegraf's address, like: 192.168.1.166:9273

## Run containers
    docker-compose up -d

## Stop containers
    docker-compose down

## Access and login into grafana
    http://${containers_host_ip}:3000
    username: admin
    password: admin

## Create Data Source
- Click 'Data Sources' button in 'Configuration' section on left pannel
- Click '+Add data source'
- Choose prometheus
- Fill the name with telegraf's address, same with the target's address in prometheus.yml
- Fill the URL section with http://${ip}:9090
  - ip is the address which containers are running on
  - 9090 is the default prometheus' port
- 'Save and Test'

## Import the dashboard
- Click 'Manage' button in 'Dashboards' section on left pannel
- Click '+Import'
- Import dashboard.json

## time sync
- ntpdate  0.cn.pool.ntp.org
- systemctl start ntpd
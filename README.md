# An agent for monitor system

## Docker image
    docker build -t telegraf:telegraf .
    docker run -it --rm -v ${PWD}:/root/go/src/github.com/influxdata/telegraf/ telegraf:telegraf

## Build
    make

## Run
   telegraf --config telegraf.conf

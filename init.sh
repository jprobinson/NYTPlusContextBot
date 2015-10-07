#!/bin/bash

# chkconfig: 3 99 05
# description: The NYTPlusContext bot daemon
# processname: bot

BIN_PATH="/opt/nyt/bin"
SERVICE="$BIN_PATH/bot"
PID_FILE="/var/run/bot.pid"
BOT_LOG="/var/nyt/logs/bot/bot.log"

function launchDaemon() {
    local id
    checkRunning
    id=$?
    if [ $id -eq 0 ]
    then
        touch $PID_FILE
        chown bot $PID_FILE
        cd $BIN_PATH
        su bot -c "$SERVICE  >> $BOT_LOG 2>&1 & echo \$! > $PID_FILE"
        echo "bot started"
    else
        echo "bot is running already"
    fi
}

function stopDaemon() {
    local is_running
    local id
    checkRunning
    is_running=$?
    if [ $is_running -eq 1 ]
    then
        id=`cat $PID_FILE`
        kill $id
        if [ $? -eq 0 ]
        then 
            echo "bot stopped"
        else
            echo "unable to stop bot"
        fi
    else
        echo "bot is stopped already"
    fi
}

function checkRunning() {
    local id
    id=`cat $PID_FILE`
    if [ -z $id ] 
    then
        return 0
    elif [ -d "/proc/$id" ] 
    then
        return 1
    else
        return 0
    fi
}

function main {
    local is_running
    case "$1" in
    start)
        launchDaemon
    ;;
    stop)
        stopDaemon
    ;;
    restart)
        stopDaemon
        sleep 5
        launchDaemon
    ;;
    status)
        checkRunning
        is_running=$?
        if [ $is_running -eq 1 ] 
        then
            echo "bot is running..."
        else
            echo "bot is not running"
        fi
    ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        exit 1
    ;;
    esac
}

main $1

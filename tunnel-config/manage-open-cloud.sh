#!/bin/bash

# Configuration
CONFIG_PATH="/home/joe/cloudflare-tunnels/open-cloud-config.yml"
TUNNEL_NAME="open-cloud"
LOG_FILE="/home/joe/cloudflare-tunnels/open-cloud.log"
PID_FILE="/home/joe/cloudflare-tunnels/open-cloud.pid"

start() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "Tunnel is already running with PID $(cat $PID_FILE)"
        return
    fi
    
    echo "Starting $TUNNEL_NAME tunnel in background..."
    nohup cloudflared tunnel --config "$CONFIG_PATH" run "$TUNNEL_NAME" > "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    echo "Tunnel started. PID: $(cat $PID_FILE)"
    echo "Logs are being written to: $LOG_FILE"
}

stop() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        echo "Stopping tunnel with PID $PID..."
        kill "$PID" 2>/dev/null || pkill -f "cloudflared tunnel --config $CONFIG_PATH"
        rm "$PID_FILE" 2>/dev/null
        echo "Tunnel stopped."
    else
        echo "PID file not found. Finding processes manually..."
        PIDS=$(pgrep -f "cloudflared tunnel --config $CONFIG_PATH")
        if [ -z "$PIDS" ]; then
            echo "No running tunnel found for $CONFIG_PATH"
        else
            echo "Killing processes: $PIDS"
            kill $PIDS
            echo "Done."
        fi
    fi
}

status() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "Status: RUNNING (PID: $(cat $PID_FILE))"
    else
        # Double check with pgrep
        if pgrep -f "cloudflared tunnel --config $CONFIG_PATH" > /dev/null; then
            echo "Status: RUNNING (PID file missing, see pgrep)"
        else
            echo "Status: NOT RUNNING"
        fi
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        stop
        sleep 2
        start
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        echo "Quick kill command: pkill -f \"cloudflared tunnel --config $CONFIG_PATH\""
        exit 1
esac

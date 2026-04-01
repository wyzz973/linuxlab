#!/bin/bash
cat > /usr/local/bin/safe_backup.sh << SBEOF
#!/bin/bash
set -euo pipefail

# Configuration
LOCKFILE="/var/run/safe_backup.lock"
PIDFILE="/var/run/safe_backup.pid"
LOGFILE="/var/log/safe_backup.log"

# Logging function
log() {
    echo "\$(date "+%Y-%m-%d %H:%M:%S") [\$\$] \$1" >> "\$LOGFILE"
}

# Cleanup on exit
cleanup() {
    rm -f "\$LOCKFILE" "\$PIDFILE"
    log "Script finished, lock released"
}
trap cleanup EXIT INT TERM

# Check for existing lock
if [ -f "\$LOCKFILE" ]; then
    OLD_PID=\$(cat "\$PIDFILE" 2>/dev/null)
    if [ -n "\$OLD_PID" ] && kill -0 "\$OLD_PID" 2>/dev/null; then
        log "ERROR: Another instance is running (PID: \$OLD_PID)"
        exit 1
    else
        log "WARN: Stale lock found, removing"
        rm -f "\$LOCKFILE" "\$PIDFILE"
    fi
fi

# Create lock
touch "\$LOCKFILE"
echo \$\$ > "\$PIDFILE"
log "Script started with PID \$\$"

# Main backup logic
log "Starting backup..."
sleep 2
log "Backup completed successfully"
SBEOF
chmod +x /usr/local/bin/safe_backup.sh

cat > /tmp/test_lock.sh << TLEOF
#!/bin/bash
# Test lock mechanism
/usr/local/bin/safe_backup.sh &
PID1=\$!
sleep 0.5
/usr/local/bin/safe_backup.sh &
PID2=\$!
wait \$PID1
wait \$PID2
echo "Lock test completed"
TLEOF
chmod +x /tmp/test_lock.sh

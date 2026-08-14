<!-- hub-reference-banner -->
> **Reference file — part of the `devops-linux-admin` hub.** Formerly the standalone `linux-sysadmin` skill.
> Sibling topics in this family are now reference files under the hubs (`devops-linux-internals`, `devops-linux-admin`, `devops-containers-cicd`, `devops-observability`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: linux-sysadmin
title: Linux Systems Administration & Networking
description: >
  Diagnose and resolve Linux and macOS systems issues using authoritative command-line tools. Covers network diagnostics (ip, ss, dig, tcpdump, nmap, iptables/nftables, macOS equivalents), system monitoring (top/htop/btop, vmstat, iostat, sar, perf, /proc and /sys), process management (ps, systemd, journalctl, kill signals, cron, launchd), log analysis (journalctl, /var/log, dmesg, OOM killer, grep/awk/sed), disk and filesystem management (fdisk/parted, LVM, ext4/xfs/btrfs/APFS, smartctl, iotop), user and permission management (chmod/chown, ACLs, sudo, PAM, SSH keys, ufw/firewalld/pf), and structured troubleshooting playbooks for high CPU, OOM, disk full, DNS failure, port conflicts, service failures, and certificate errors.
  TRIGGER: user diagnoses network connectivity, DNS, or firewall issues on Linux or macOS; investigates high CPU, memory, or I/O problems; manages processes or systemd services; parses logs; manages disk partitions or LVM; works through a sysadmin troubleshooting playbook.
  SKIP: application-level debugging (use the language skill); Kubernetes/container orchestration internals (use kubernetes-networking); cloud-provider networking (AWS VPC, Azure vNets); Windows Server administration.
version: "1.2.1"
category: developer
updated: "2026-06-12"
tags:
  - linux
  - sysadmin
  - networking
  - monitoring
  - troubleshooting
  - macos
  - devops
  - security
related_skills:
  - devops-infra
  - programming-languages
---

# Linux Systems Administration & Networking

## Quick Reference

| Section | Topics |
|---------|--------|
| [1 — Network Diagnostics](#section-1--network-diagnostics) | ip, ss, ping, traceroute, dig, curl, nmap, tcpdump, ip netns, iptables/nftables, pf |
| [2 — System Monitoring](#section-2--system-monitoring) | top/htop/btop, vmstat, iostat, sar, perf, /proc, /sys, macOS equivalents |
| [3 — Process Management](#section-3--process-management) | ps, signals, systemd, journalctl, cron, launchd |
| [4 — Log Analysis](#section-4--log-analysis--debugging) | journalctl, /var/log, dmesg, OOM killer, grep/awk/sed, logrotate |
| [5 — Disk & Filesystem](#section-5--disk--filesystem) | fdisk/parted, mount/fstab, LVM, ext4/xfs/btrfs/APFS, smartctl, iotop |
| [6 — Users & Permissions](#section-6--user--permission-management) | chmod/chown, ACLs, sudo, PAM, SSH keys, firewall |
| [7 — Troubleshooting Playbooks](#section-7--troubleshooting-playbooks) | High CPU, OOM, disk full, network unreachable, DNS, port conflict, cert errors |

> **Distro note:** Commands marked `(Debian/Ubuntu)` or `(RHEL/CentOS/Fedora)` differ by distribution. macOS equivalents are noted inline throughout.

## When to use this skill

- Diagnosing network connectivity, DNS, or firewall issues on Linux or macOS
- Investigating high CPU, high memory, or I/O performance problems
- Managing processes, services, and systemd units
- Parsing logs to find root cause of failures
- Managing disk partitions, filesystems, and LVM volumes
- Hardening file permissions, SSH access, and firewall rules
- Working through a structured troubleshooting playbook for common failure modes
- Diagnosing container networking at the Linux layer (docker network inspect, ip netns)

## When NOT to use this skill

- Application-level debugging (use language-specific debuggers or profilers)
- Kubernetes/container orchestration internals (kubectl, crictl — use kubernetes-networking skill)
- Cloud provider-specific networking (AWS VPC, Azure vNets — use provider docs)
- Windows Server administration

---

## Section 1 — Network Diagnostics

### 1.1 Interface and Address Management

```bash
# Modern iproute2 (preferred over net-tools)
ip addr show                      # list all interfaces and addresses
ip addr show eth0                 # specific interface
ip link set eth0 up               # bring interface up
ip route show                     # routing table
ip route add default via 192.168.1.1  # add default gateway
ip -s link show eth0              # interface statistics (rx/tx bytes, errors)

# Legacy (still available on many systems)
ifconfig -a
ifconfig eth0 192.168.1.100 netmask 255.255.255.0

# macOS equivalents
ifconfig -a                       # macOS still uses ifconfig
networksetup -listallnetworkservices
networksetup -getinfo "Wi-Fi"
networksetup -setmanual "Wi-Fi" 192.168.1.50 255.255.255.0 192.168.1.1
networksetup -setdnsservers "Wi-Fi" 8.8.8.8 8.8.4.4
scutil --nwi                      # network interface information
scutil --dns                      # DNS configuration
```

### 1.2 Socket and Connection Analysis

```bash
# ss — replacement for netstat, much faster on systems with many connections
ss -tunlp                         # TCP+UDP, listening, numeric, with process
ss -s                             # summary statistics
ss -t state established           # established TCP connections
ss -o state TIME-WAIT             # TIME-WAIT connections
ss 'dport = :443'                 # connections to port 443
ss -tnp dst 10.0.0.1              # connections to specific host

# netstat (legacy, part of net-tools)
netstat -tunlp                    # same as ss -tunlp
netstat -rn                       # routing table
netstat -i                        # interface statistics

# macOS: netstat works, lsof is often more useful
netstat -anv | grep LISTEN
lsof -iTCP -sTCP:LISTEN -P        # listening TCP ports with PIDs
lsof -i :8080                     # what is using port 8080
```

### 1.3 Connectivity Testing

```bash
# ping
ping -c 4 google.com              # 4 packets
ping -i 0.2 -c 100 host           # flood-style with 0.2s interval
ping6 ::1                         # IPv6 loopback
ping -W 1 192.168.1.1             # 1-second timeout

# traceroute / mtr
traceroute -n google.com          # numeric, skip DNS reverse lookup
traceroute -T -p 443 host         # TCP traceroute (bypasses ICMP blocks)
mtr --report --report-cycles 20 google.com   # combined ping+traceroute report
mtr -n --report host              # numeric mtr report

# macOS equivalents (traceroute works; mtr needs homebrew)
traceroute google.com
sudo traceroute -T -p 443 host    # TCP mode
brew install mtr && sudo mtr host
```

### 1.4 DNS Resolution

```bash
# dig — most informative
dig google.com                    # A record
dig google.com MX                 # mail exchangers
dig @8.8.8.8 google.com           # query specific resolver
dig +trace google.com             # full delegation chain
dig +short google.com AAAA        # IPv6, terse output
dig -x 8.8.8.8                    # reverse DNS lookup

# nslookup (interactive or one-shot)
nslookup google.com
nslookup -type=TXT _dmarc.google.com

# host (concise)
host google.com
host -t MX google.com

# macOS: all three work; also (depth: references/macos-networking.md):
sudo dscacheutil -flushcache       # clears the Directory Services cache (DNS cache lives in mDNSResponder)
sudo killall -HUP mDNSResponder    # drop mDNSResponder's DNS cache + re-read config (not a restart)
scutil --dns                       # view current resolver configuration
```

### 1.5 HTTP and Application-Layer Testing

```bash
# curl
curl -v https://example.com                     # verbose: headers + body
curl -I https://example.com                     # HEAD request only
curl -o /dev/null -w "%{http_code}\n" URL       # just the status code
curl -L -k URL                                  # follow redirects, skip TLS verify
curl --resolve example.com:443:1.2.3.4 https://example.com  # override DNS
curl -x socks5://proxy:1080 https://example.com # through SOCKS proxy
curl -w "@curl-format.txt" -s -o /dev/null URL  # detailed timing breakdown

# wget
wget -q --spider https://example.com            # check URL existence

# openssl — TLS certificate inspection
openssl s_client -connect host:443 </dev/null 2>/dev/null | openssl x509 -noout -dates
echo | openssl s_client -servername example.com -connect example.com:443 2>/dev/null \
  | openssl x509 -noout -text
```

### 1.6 Port Scanning and Discovery

```bash
# nmap
nmap -sn 192.168.1.0/24           # ping sweep (host discovery)
nmap -p 22,80,443 host            # specific ports
nmap -sV -p- host                 # service version detection, all ports
nmap -sU -p 53,161 host           # UDP scan
nmap --script vuln host           # NSE vulnerability scripts
nmap -A -T4 host                  # aggressive scan (OS detect, version, scripts)
```

### 1.7 Packet Capture

```bash
# tcpdump
tcpdump -i eth0                   # capture on interface
tcpdump -i any -n port 443        # all interfaces, numeric, port filter
tcpdump -w /tmp/capture.pcap -i eth0  # write to file
tcpdump -r /tmp/capture.pcap      # read capture file
tcpdump 'tcp[tcpflags] & tcp-syn != 0'   # SYN packets only
tcpdump -i eth0 host 10.0.0.1 and port 80

# tshark (wireshark CLI)
tshark -i eth0 -Y "http.request"  # HTTP requests
tshark -r file.pcap -T fields -e ip.src -e http.request.uri
```

### 1.8 Network Namespaces (Containers / VMs)

```bash
# ip netns — manage network namespaces (used by Docker, LXC, systemd-nspawn)
ip netns list                     # list named namespaces
ip netns exec <ns> ip addr        # run command inside a namespace
ip netns exec <ns> ping 8.8.8.8  # test connectivity from inside namespace

# Docker container networking
docker network ls                 # list docker networks
docker network inspect bridge     # inspect bridge details (subnets, containers)
docker exec <ctr> ip addr         # interface inside container
docker exec <ctr> ss -tlnp        # listening ports inside container

# Find the veth peer for a container interface
ip link show type veth            # list all veth pairs on host
# Match ifindex in container's /sys/class/net/<iface>/iflink to host interface

# Inspect a container's network namespace directly
pid=$(docker inspect -f '{{.State.Pid}}' <ctr>)
nsenter -t $pid -n ip addr        # enter container netns without docker exec
nsenter -t $pid -n ss -tlnp
```

### 1.9 Firewall: iptables / nftables / ufw / firewalld

```bash
# iptables (legacy; still widely used)
iptables -L -n -v                 # list rules with counts
iptables -A INPUT -p tcp --dport 22 -j ACCEPT
iptables -A INPUT -j DROP         # default deny (append to end)
iptables-save > /etc/iptables/rules.v4
iptables-restore < /etc/iptables/rules.v4

# nftables (modern replacement)
nft list ruleset                  # show all rules
nft add table inet filter
nft add chain inet filter input { type filter hook input priority 0 \; policy drop \; }
nft add rule inet filter input tcp dport 22 accept
nft add rule inet filter input ct state established,related accept
nft list table inet filter

# ufw (Ubuntu/Debian simplified frontend)
ufw status verbose
ufw allow 22/tcp
ufw allow from 10.0.0.0/8 to any port 3306
ufw deny 23
ufw enable
ufw delete allow 22/tcp

# firewalld (RHEL/CentOS/Fedora)
firewall-cmd --state
firewall-cmd --list-all
firewall-cmd --add-service=http --permanent
firewall-cmd --add-port=8080/tcp --permanent
firewall-cmd --reload

# macOS pf (packet filter)
sudo pfctl -sr                    # show current rules
sudo pfctl -f /etc/pf.conf        # load ruleset
sudo pfctl -e                     # enable pf
sudo pfctl -d                     # disable pf
# /etc/pf.conf example:
# pass in on en0 proto tcp to any port {80 443}
# block in quick from <bruteforce>
```

---

## Section 2 — System Monitoring

### 2.1 CPU and Process Overview

```bash
# top — interactive process viewer
top -b -n 1                       # batch mode, one snapshot
top -p 1234                       # monitor specific PID
# Inside top: P=sort CPU, M=sort memory, k=kill, r=renice, q=quit

# htop — colorful interactive viewer (install separately)
htop -p 1234,5678                 # specific PIDs
htop -u www-data                  # filter by user
# htop keys: F5=tree, F6=sort, F9=kill, space=tag

# btop — modern alternative with graphs
btop                              # interactive with bandwidth graphs

# macOS Activity Monitor CLI
top -l 1 -stats pid,command,cpu,mem | head -20  # one snapshot
```

### 2.2 Memory Statistics

```bash
# free
free -h                           # human-readable (MB/GB)
free -h -s 2                      # refresh every 2 seconds
# "available" column is what matters, not "free"
# Linux aggressively uses RAM for page cache; high cache usage is healthy

# /proc filesystem
grep -E 'MemTotal|MemFree|MemAvailable|Buffers|Cached|SwapUsed' /proc/meminfo
cat /proc/vmstat                  # virtual memory statistics
```

### 2.2b dstat — Combined Real-time Overview

```bash
# dstat — combines vmstat + iostat + ifstat in one view (install: apt/yum install dstat)
dstat                             # default: cpu, disk, net, paging, system
dstat -cdngy 1                    # cpu, disk, net, page, sys — 1s interval
dstat --top-cpu --top-mem         # show top CPU and memory consumers per interval
dstat -D sda,sdb                  # specific disks only
dstat --output /tmp/dstat.csv 1 60  # save 60s of stats to CSV
# Note: dstat is deprecated upstream; modern alternative is dool (drop-in replacement)
```

### 2.3 CPU and I/O Performance Deep Dive

```bash
# vmstat — combined CPU + memory + swap + I/O overview
vmstat 1 10                       # 10 reports, 1-second interval
# CPU columns: us=user, sy=system, id=idle, wa=iowait, st=steal (VM)
# Memory/swap: si/so=swap-in/out — any nonzero value warrants investigation
# r=run queue length; b=processes blocked on I/O
# Rule: si/so > 0 → swapping → check free -h immediately

# iostat (from sysstat package)
iostat -xz 1                      # extended, suppress zero-activity, 1s interval
# Key: await=avg I/O wait ms, %util=device saturation, r/s and w/s=throughput

# sar — historical data (sysstat)
sar -u 1 5                        # CPU 5 samples, 1s interval
sar -r 1 5                        # memory
sar -d 1 5                        # disk
sar -n DEV 1 5                    # network interfaces
sar -f /var/log/sa/sa01           # read saved data file (date 01)

# For programmable in-kernel tracing beyond perf/strace/ftrace — bpftrace one-liners,
# bcc tools (execsnoop/opensnoop/biolatency), kprobe/uprobe/tracepoint/fentry, CO-RE,
# and continuous profiling (Parca/Pixie) — load references/ebpf-observability.md in this hub.
# For the perf/ftrace/USE-method/PSI methodology in depth (sampling, flame graphs,
# function_graph, trace-cmd, /proc/pressure) — load references/linux-perf-tracing.md.
# For nftables rulesets, netns/veth/bridge topologies, tc qdiscs/QoS, and XDP in depth —
# load references/linux-networking-stack.md.

# perf (Linux performance counters)
perf stat ls                      # hardware counters for command
perf top                          # live top-like CPU profiler
perf record -g ./app && perf report   # flame-graph style call chain
perf stat -e cache-misses,cache-references ./app

# macOS equivalents
vm_stat                           # virtual memory stats
fs_usage                          # per-process filesystem activity (sudo required)
sudo fs_usage -e -f filesystem    # filter filesystem events
dtrace -n 'syscall:::entry { @[execname] = count(); }' -c "sleep 5"
```

### 2.4 Disk Space

```bash
# df — filesystem space usage
df -h                             # human-readable
df -hi                            # include inode usage
df -T                             # show filesystem type

# du — directory size
du -sh /var/log/*                 # size of each item in /var/log
du -sh --max-depth=1 /            # top-level directory sizes
du -ah | sort -rh | head -20      # find 20 largest files/dirs

# find large files
find / -type f -size +100M -exec ls -lh {} \; 2>/dev/null | sort -k5 -rh | head -20
```

### 2.5 Open Files and Descriptors

```bash
# lsof — list open files
lsof -p 1234                      # all files opened by PID 1234
lsof /var/log/app.log             # processes using a specific file
lsof -i :8080                     # process listening on port 8080
lsof -i TCP@192.168.1.1           # connections to specific host
lsof -u username                  # all files by user
lsof +D /mnt/data                 # everything under a directory — use to find what blocks umount

# /proc per-process
ls -la /proc/1234/fd/             # file descriptors for PID 1234
cat /proc/1234/status             # process status including memory maps
cat /proc/1234/cmdline | tr '\0' ' '  # full command line
```

### 2.6 System Call Tracing

```bash
# strace
strace ls                         # trace all syscalls
strace -p 1234                    # attach to running process
strace -e trace=open,read,write ls  # filter specific syscalls
strace -c ls                      # summary statistics
strace -f -o /tmp/trace.txt ./app  # follow forks, write to file
strace -e trace=network curl google.com  # network syscalls only

# ltrace (library calls)
ltrace ls

# macOS: dtruss (dtrace-based strace equivalent)
sudo dtruss ls                    # requires SIP disabled or specific entitlements
```

### 2.7 /proc and /sys Filesystem

```bash
# Key /proc entries
/proc/cpuinfo                     # CPU information
/proc/loadavg                     # 1, 5, 15-minute load averages + runnable/total
/proc/net/tcp                     # raw TCP socket table
/proc/net/dev                     # interface statistics
/proc/sys/vm/swappiness           # kernel swappiness (0-100)
/proc/sys/net/ipv4/ip_forward     # IP forwarding state

# Tuning at runtime (not persistent)
echo 10 > /proc/sys/vm/swappiness
sysctl -w vm.swappiness=10
sysctl -w net.ipv4.ip_forward=1

# Persistent sysctl settings
echo "vm.swappiness=10" >> /etc/sysctl.conf
sysctl -p                         # reload /etc/sysctl.conf

# /sys
/sys/class/net/eth0/statistics/   # NIC stats
/sys/block/sda/queue/             # disk queue settings
echo deadline > /sys/block/sda/queue/scheduler   # change I/O scheduler
```

---

## Section 3 — Process Management

### 3.1 Process Inspection

```bash
# ps
ps aux                            # all processes, user-oriented
ps auxf                           # process tree (forest)
ps -eo pid,ppid,user,%cpu,%mem,cmd --sort=-%cpu | head -20
ps -p 1234 -o pid,ppid,stat,command

# pgrep / pkill — search by name
pgrep nginx                       # list PIDs matching "nginx"
pgrep -a nginx                    # show full command line
pkill -SIGTERM nginx              # send signal to matched processes
pkill -u www-data                 # kill all processes by user
pgrep -x "exact-name"            # exact match

# Process tree
pstree -p                         # show PIDs in tree
pstree -u username                # tree for user
```

### 3.2 Signals

| Signal | Number | Meaning | Behavior |
|--------|--------|---------|----------|
| SIGTERM | 15 | Graceful termination | Can be caught; process should clean up |
| SIGKILL | 9 | Forced termination | Cannot be caught; kernel kills immediately |
| SIGHUP | 1 | Hangup / reload | Many daemons re-read config on SIGHUP |
| SIGUSR1 | 10 | User-defined 1 | App-specific; often trigger log rotation |
| SIGUSR2 | 12 | User-defined 2 | App-specific |
| SIGSTOP | 19 | Pause process | Cannot be caught |
| SIGCONT | 18 | Resume process | Resume paused process |
| SIGINT | 2 | Interrupt (Ctrl-C) | Can be caught |

```bash
kill -15 1234                     # SIGTERM — always try this first
kill -9 1234                      # SIGKILL — last resort
kill -1 $(pgrep nginx)            # reload nginx config
kill -SIGUSR1 $(pgrep -x rsyslogd)  # rotate rsyslog logs

# NOTE: SIGKILL cannot reach a process in D state (uninterruptible sleep)
# D-state processes are stuck in kernel waiting for I/O; only rebooting helps
```

### 3.3 systemd — Service Management

> **Going deeper than ops basics?** This section covers day-to-day service control.
> For systemd internals — the unit dependency/ordering model (Wants/Requires/After,
> Type=, socket activation, drop-ins, generators), journald storage/FSS/namespaces,
> cgroup-v2 resource control (slices, MemoryHigh/Max, CPUWeight/Quota, IOWeight,
> TasksMax, Delegate), systemd-oomd (PSI-driven OOM), and image-based extensions
> (sysext/confext, portable services, sandboxing) — **Read `references/systemd.md`**
> in this hub.

```bash
# Service lifecycle
systemctl start nginx
systemctl stop nginx
systemctl restart nginx           # stop + start
systemctl reload nginx            # send SIGHUP / graceful reload
systemctl enable nginx            # enable at boot
systemctl disable nginx
systemctl is-active nginx
systemctl is-enabled nginx
systemctl status nginx            # status + recent journal lines

# System-wide state
systemctl list-units --type=service --state=failed
systemctl list-units --type=service --all
systemctl list-timers --all       # systemd timers

# Dependency inspection
systemctl list-dependencies nginx
systemctl show -p After nginx

# Masking (prevent a service from being started by anything)
systemctl mask nginx
systemctl unmask nginx
```

### 3.4 systemd Unit Files

```ini
# /etc/systemd/system/myapp.service
[Unit]
Description=My Application
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=myapp
Group=myapp
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/bin/server --config /etc/myapp/config.yaml
ExecReload=/bin/kill -HUP $MAINPID
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal
SyslogIdentifier=myapp
# Security hardening
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ReadWritePaths=/var/lib/myapp /var/log/myapp

[Install]
WantedBy=multi-user.target
```

```bash
# After editing unit files
systemctl daemon-reload
systemctl restart myapp
```

### 3.5 journalctl — Log Querying

```bash
journalctl -u nginx               # logs for nginx unit
journalctl -u nginx -f            # follow (live tail)
journalctl -u nginx --since "1 hour ago"
journalctl -u nginx --since "2024-01-01" --until "2024-01-02"
journalctl -p err -b              # errors since last boot
journalctl -b -1                  # previous boot
journalctl -k                     # kernel messages only (like dmesg)
journalctl --disk-usage           # how much journal space is used
journalctl --vacuum-time=7d       # delete journals older than 7 days
journalctl -o json-pretty -n 10   # JSON output
journalctl _PID=1234              # by PID
journalctl CONTAINER_NAME=web     # by container (Docker journal driver)
```

### 3.6 Scheduling: cron and systemd timers

```bash
# cron
crontab -e                        # edit current user's crontab
crontab -l                        # list current user's crontab
```

**Cron field reference:**

| Field | Values | Example |
|-------|--------|---------|
| minute | 0–59 | `*/5` = every 5 min |
| hour | 0–23 | `2` = 2 AM |
| day-of-month | 1–31 | `1` = 1st of month |
| month | 1–12 | `*` = every month |
| day-of-week | 0–7 (0=Sun) | `1-5` = weekdays |

```
# Examples
*/5 * * * * /usr/local/bin/check-health.sh
0 2 * * * /usr/bin/backup.sh >> /var/log/backup.log 2>&1
@reboot /usr/local/bin/startup.sh

# System-wide drop-ins: /etc/cron.d/, /etc/cron.daily/, /etc/cron.hourly/
```

```ini
# systemd timers (preferred over cron for system services)
# /etc/systemd/system/myapp-backup.timer
[Timer]
OnCalendar=*-*-* 02:00:00
Persistent=true                   # catch up missed runs after reboot
RandomizedDelaySec=900            # spread load across a 15-min window
```

```bash
systemctl enable --now myapp-backup.timer
systemctl list-timers --all
```

### 3.7 macOS launchd

```bash
# launchctl — macOS equivalent of systemctl
launchctl list                    # list all loaded agents/daemons
launchctl list com.example.myapp  # specific service
launchctl load ~/Library/LaunchAgents/com.example.myapp.plist
launchctl unload ~/Library/LaunchAgents/com.example.myapp.plist
launchctl start com.example.myapp
launchctl stop com.example.myapp

# plist locations
# ~/Library/LaunchAgents/         — per-user, run as user
# /Library/LaunchDaemons/         — system-wide, run as root
# /System/Library/LaunchDaemons/  — Apple daemons (don't edit)

# macOS process management
launchctl kickstart gui/501/com.example.myapp   # force start
launchctl bootout gui/501/com.example.myapp     # force unload
```

---

## Section 4 — Log Analysis & Debugging

### 4.1 /var/log Hierarchy

| Path | Contents |
|------|----------|
| `/var/log/syslog` | General system messages (Debian/Ubuntu) |
| `/var/log/messages` | General system messages (RHEL/CentOS) |
| `/var/log/auth.log` | Authentication events (Debian/Ubuntu) |
| `/var/log/secure` | Authentication events (RHEL/CentOS) |
| `/var/log/kern.log` | Kernel messages |
| `/var/log/dmesg` | Boot-time kernel ring buffer |
| `/var/log/nginx/` | Nginx access and error logs |
| `/var/log/apache2/` | Apache access and error logs |
| `/var/log/mysql/` | MySQL/MariaDB error log |
| `/var/log/postgresql/` | PostgreSQL logs |
| `/var/log/dpkg.log` | Package install history (Debian) |
| `/var/log/yum.log` | Package install history (RHEL older) |
| `/var/log/dnf.log` | Package install history (RHEL newer) |

```bash
# dmesg — kernel ring buffer
dmesg -T                          # human-readable timestamps
dmesg -T -l err,crit,alert,emerg  # errors and above
dmesg -T --follow                 # live kernel messages
dmesg | grep -i "error\|fail\|warn" | tail -30

# macOS
log show --last 1h --predicate 'process == "kernel"'
log show --last 30m --level error
```

### 4.2 grep / awk / sed Log Parsing

```bash
# grep patterns
grep -E "ERROR|FATAL|CRITICAL" /var/log/app.log
grep -c "404" /var/log/nginx/access.log    # count occurrences
grep -n "OutOfMemoryError" app.log          # show line numbers
grep -B 5 -A 10 "FATAL" app.log             # 5 lines before, 10 after
grep -r "Connection refused" /var/log/       # recursive search

# awk — structured field extraction
# Apache/nginx access log: extract status codes and count them
awk '{print $9}' /var/log/nginx/access.log | sort | uniq -c | sort -rn

# Extract IPs hitting 404s
awk '$9 == 404 {print $1}' /var/log/nginx/access.log | sort | uniq -c | sort -rn | head

# Average response time from nginx logs (if logged)
awk '{sum += $NF; count++} END {print sum/count " ms avg"}' /var/log/nginx/access.log

# sed — transform log output
sed -n '/2024-01-15 14:/,/2024-01-15 15:/p' app.log   # time range slice
sed 's/password=[^ ]*/password=REDACTED/g' auth.log    # redact passwords

# Combined pipeline — errors in last hour with timestamps
journalctl -u myapp --since "1 hour ago" | grep -i error | awk '{print $1, $2, $3, $NF}'
```

### 4.3 OOM Killer

```bash
# Detect OOM events
dmesg -T | grep -i "out of memory\|oom-killer\|killed process"
grep -i "oom-killer\|out of memory" /var/log/syslog
journalctl -k | grep -i "oom\|killed"

# OOM log shows: victim process, oom_score, per-process memory map
# Example:
# [12345.67] Out of memory: Killed process 4589 (postgres) total-vm:8388608kB
#            anon-rss:4194304kB, file-rss:0kB, shmem-rss:0kB

# View OOM scores (higher = more likely to be killed)
for pid in /proc/[0-9]*/; do
  pid=${pid%/}; pid=${pid##*/}
  printf "%d\t%s\t%s\n" "$pid" \
    "$(cat /proc/$pid/oom_score 2>/dev/null)" \
    "$(cat /proc/$pid/comm 2>/dev/null)"
done | sort -t$'\t' -k2 -rn | head -20

# Protect a process from OOM killer
echo -1000 > /proc/$(pgrep -x postgres)/oom_score_adj

# Check current memory pressure
grep -E 'MemAvailable|SwapFree|Dirty' /proc/meminfo
```

### 4.4 Common Error Patterns

```bash
# Connection refused — service not listening
grep "Connection refused" /var/log/app.log
ss -tlnp | grep :PORT              # verify service is listening

# Disk full — ENOSPC
dmesg | grep -i "no space\|ENOSPC"
df -h                              # find full filesystem
du -sh /var/log/* | sort -rh       # find log directory culprit

# DNS failures
grep -i "NXDOMAIN\|SERVFAIL\|timeout" /var/log/syslog
resolvectl status                  # check systemd-resolved state (replaces deprecated systemd-resolve)

# Permission denied
grep "permission denied\|EACCES" /var/log/app.log
# Check with strace: strace -e trace=open,openat app 2>&1 | grep EACCES

# Certificate errors
openssl s_client -connect host:443 </dev/null 2>&1 | grep -E "verify|Verify|error"
# Check expiry
echo | openssl s_client -connect host:443 2>/dev/null | openssl x509 -noout -enddate

# Failed SSH auth
grep "Failed password\|Invalid user\|authentication failure" /var/log/auth.log
journalctl -u sshd | grep "Failed\|Invalid"
```

### 4.5 logrotate

```bash
# Test logrotate config
logrotate -d /etc/logrotate.conf   # dry run
logrotate -f /etc/logrotate.conf   # force rotation now

# /etc/logrotate.d/myapp
# /var/log/myapp/*.log {
#     daily
#     rotate 14
#     compress
#     delaycompress
#     missingok
#     notifempty
#     sharedscripts
#     postrotate
#         systemctl reload myapp
#     endscript
# }
```

---

## Section 5 — Disk & Filesystem

### 5.1 Partitioning

```bash
# List block devices
lsblk -f                          # with filesystem types and UUIDs
lsblk -o NAME,SIZE,TYPE,FSTYPE,MOUNTPOINT,UUID
blkid                             # UUID and filesystem of all block devices

# fdisk (MBR, good for < 2TB)
fdisk /dev/sdb
# Interactive: n=new, d=delete, p=print, w=write, q=quit

# parted (GPT, required for > 2TB)
parted /dev/sdb print
parted /dev/sdb mklabel gpt
parted /dev/sdb mkpart primary ext4 0% 100%
parted /dev/sdb --script mkpart primary xfs 0% 50%

# Create filesystems
mkfs.ext4 -L mydata /dev/sdb1
mkfs.xfs -L mydata /dev/sdb1
mkfs.btrfs -L mydata /dev/sdb1

# macOS
diskutil list
diskutil info /dev/disk2
diskutil eraseDisk APFS MyDisk GPTFormat /dev/disk2
diskutil partitionDisk /dev/disk2 GPT APFS MyVol 100%
```

### 5.2 Mounting

```bash
# Mount / unmount
mount /dev/sdb1 /mnt/data
mount -t ext4 -o noatime,nodiratime /dev/sdb1 /mnt/data
umount /mnt/data
umount -l /mnt/data               # lazy unmount (when busy)

# Show mounts
mount | grep "^/dev"
findmnt --real                    # tree view
cat /proc/mounts

# /etc/fstab — persistent mounts
# UUID=xxxx /mnt/data ext4 defaults,noatime 0 2
# Device    Mountpoint  Type   Options   Dump  Pass
# 0=no dump; Pass: 0=no check, 1=root, 2=other

# Mount with UUID (always use UUID, not device name)
blkid /dev/sdb1                   # get UUID
echo "UUID=$(blkid -s UUID -o value /dev/sdb1) /mnt/data ext4 defaults,noatime 0 2" >> /etc/fstab
mount -a                          # mount all in fstab
```

### 5.3 Filesystem Types

| Filesystem | Use Case | Key Features |
|-----------|----------|--------------|
| ext4 | General-purpose Linux | Stable, journaled, broad support |
| xfs | Large files, high throughput | RHEL default, excellent scalability, can't shrink |
| btrfs | Snapshots, compression, RAID | CoW, subvolumes, checksums; less battle-tested |
| APFS | macOS (SSD-optimized) | CoW, snapshots, encryption; macOS only |
| tmpfs | RAM-backed temp storage | Survives restart as memory allows |
| nfs | Network filesystem | Shared storage across hosts |

```bash
# btrfs operations
btrfs filesystem show
btrfs subvolume list /
btrfs subvolume create /mnt/data/subvol1
btrfs subvolume snapshot /mnt/data /mnt/snapshots/snap-$(date +%Y%m%d)
btrfs filesystem df /mnt/data

# ext4 operations
e2fsck -f /dev/sdb1               # filesystem check (unmounted)
resize2fs /dev/sdb1 50G           # resize (can grow or shrink with ext4)
tune2fs -l /dev/sdb1              # filesystem parameters

# xfs operations
xfs_repair /dev/sdb1              # filesystem repair (unmounted)
xfs_growfs /mnt/data              # grow to fill partition (mounted)
xfs_info /mnt/data                # filesystem information
```

### 5.4 LVM (Logical Volume Manager)

```bash
# Inspect existing LVM
pvs                               # physical volumes
vgs                               # volume groups
lvs                               # logical volumes
pvdisplay /dev/sdb                # detailed PV info

# Create LVM stack
pvcreate /dev/sdb /dev/sdc        # initialize PVs
vgcreate myvg /dev/sdb /dev/sdc   # create VG from PVs
lvcreate -L 100G -n mydata myvg   # create 100G LV
mkfs.ext4 /dev/myvg/mydata        # format
mount /dev/myvg/mydata /mnt/data

# Extend LV and filesystem
lvextend -L +50G /dev/myvg/mydata      # extend LV by 50G
resize2fs /dev/myvg/mydata             # extend ext4 to fill LV
lvextend -L +50G -r /dev/myvg/mydata   # extend + resize filesystem in one step

# LVM snapshot
lvcreate -s -L 10G -n mydata-snap /dev/myvg/mydata
mount -o ro /dev/myvg/mydata-snap /mnt/snapshot
```

### 5.5 SMART Monitoring and Disk I/O

```bash
# smartctl (smartmontools package)
smartctl -a /dev/sda              # full SMART report
smartctl -H /dev/sda              # health check only
smartctl -t short /dev/sda        # start short self-test
smartctl -l selftest /dev/sda     # view self-test history
# Key attributes to watch: Reallocated_Sector_Ct, Uncorrectable_Error_Cnt,
#   Spin_Retry_Count, Current_Pending_Sector

# iotop — I/O per process (requires root)
iotop -o                          # only show active I/O
iotop -a                          # accumulated I/O

# iostat for disk analysis
iostat -xz 1                      # 1-second interval, extended stats
# await > 10ms: I/O latency concern
# %util approaching 100%: disk saturation

# blktrace — block-level I/O tracing (advanced)
blktrace -d /dev/sda -o /tmp/trace
blkparse /tmp/trace.sda.blktrace.* | head -50

# macOS
diskutil info /dev/disk0          # includes S.M.A.R.T. status
iostat -d 1                       # disk I/O stats (macOS)
```

---

## Section 6 — User & Permission Management

### 6.1 File Permissions

```bash
# chmod — change permissions
chmod 644 file.txt                # rw-r--r-- (owner rw, group r, other r)
chmod 755 script.sh               # rwxr-xr-x (owner rwx, group rx, other rx)
chmod 700 ~/.ssh                  # rwx------ (owner only)
chmod u+x script.sh               # add execute for owner
chmod g-w file                    # remove group write
chmod o= file                     # remove all other permissions
chmod -R 755 /var/www/html        # recursive

# chown — change ownership
chown alice file.txt
chown alice:developers file.txt   # owner and group
chown -R www-data:www-data /var/www/html

# Special permission bits
chmod u+s /usr/local/bin/app      # SUID: run as owner
chmod g+s /shared/dir             # SGID: new files inherit group
chmod +t /tmp                     # sticky bit: only owner can delete

# Verify permissions
ls -la file.txt
stat file.txt                     # detailed numeric permissions
```

### 6.2 ACLs (Access Control Lists)

```bash
# getfacl / setfacl (acl package)
getfacl /var/www/html             # view current ACLs
setfacl -m u:alice:rwx /var/www/html     # grant user alice rwx
setfacl -m g:developers:rx /var/www/html # grant group read+execute
setfacl -x u:alice /var/www/html  # remove alice's ACL
setfacl -R -m u:alice:rwx /opt/app  # recursive ACL
setfacl -d -m g:developers:rx /shared  # default ACL for new files

# Check if filesystem supports ACLs (must be mounted with acl option)
tune2fs -l /dev/sda1 | grep "Default mount options"
# Add to /etc/fstab: UUID=xxx /mnt/data ext4 defaults,acl 0 2
```

### 6.3 sudo Configuration

```bash
# Edit with visudo (validates syntax before saving)
visudo
visudo -f /etc/sudoers.d/myapp    # edit supplemental file

# /etc/sudoers.d/myapp examples:
# alice ALL=(ALL:ALL) ALL                    # full sudo
# alice ALL=(ALL) NOPASSWD: /usr/bin/systemctl restart nginx
# %developers ALL=(ALL) /usr/local/bin/deploy.sh  # group sudo
# alice ALL=(www-data) NOPASSWD: ALL         # run as www-data only

# List sudo privileges
sudo -l                           # current user's sudo rules
sudo -l -U alice                  # another user's rules (as root)
```

### 6.4 SSH Key Management

```bash
# Generate keys
ssh-keygen -t ed25519 -C "user@host"          # modern, recommended
ssh-keygen -t rsa -b 4096 -C "user@host"      # RSA fallback
ssh-keygen -t ecdsa -b 256 -C "user@host"     # ECDSA

# Key file permissions — SSH will refuse wrong permissions
chmod 700 ~/.ssh
chmod 600 ~/.ssh/id_ed25519       # private key: owner read only
chmod 644 ~/.ssh/id_ed25519.pub   # public key: any read
chmod 600 ~/.ssh/authorized_keys  # authorized keys
chmod 600 ~/.ssh/config           # SSH client config

# Deploy public key
ssh-copy-id -i ~/.ssh/id_ed25519.pub user@host
# Or manually:
cat ~/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys

# ~/.ssh/config
# Host bastion
#   HostName 10.0.0.5
#   User ec2-user
#   IdentityFile ~/.ssh/prod-key
#   ForwardAgent yes
#
# Host internal-server
#   HostName 10.0.1.10
#   User ubuntu
#   ProxyJump bastion

# ssh-agent
eval $(ssh-agent)
ssh-add ~/.ssh/id_ed25519
ssh-add -l                        # list loaded keys

# Audit authorized_keys
for user in $(cut -f1 -d: /etc/passwd); do
  home=$(getent passwd "$user" | cut -f6 -d:)
  if [ -f "$home/.ssh/authorized_keys" ]; then
    echo "=== $user ==="; cat "$home/.ssh/authorized_keys"
  fi
done
```

### 6.5 PAM Basics

```bash
# /etc/pam.d/ — PAM configuration files
# /etc/pam.d/sshd — PAM config for SSH
# /etc/pam.d/sudo — PAM config for sudo
# /etc/pam.d/common-auth — shared auth config (Debian)

# PAM module types: auth, account, session, password
# Control flags: required, requisite, sufficient, optional

# Common modules
# pam_unix.so     — standard Unix password auth
# pam_faillock.so — lock accounts after N failures
# pam_limits.so   — enforce /etc/security/limits.conf
# pam_env.so      — set environment variables
# pam_mkhomedir.so — auto-create home directories
# pam_google_authenticator.so — TOTP 2FA

# Check PAM is enforcing limits
cat /etc/security/limits.conf     # resource limits per user
ulimit -a                         # current session limits
```

---

## Section 7 — Troubleshooting Playbooks

### 7.1 High CPU

```bash
# Step 1 — Identify the hog
top -b -n 1 -o %CPU | head -20
ps aux --sort=-%cpu | head -10

# Step 2 — Dig deeper into the top offender (replace <PID>)
perf top -p <PID>                 # what function is hot?
strace -p <PID> -c                # what syscalls dominate?
cat /proc/<PID>/wchan             # what kernel function is it waiting on?

# Step 3 — Check process state
ps -o pid,ppid,%cpu,stat,cmd -p <PID>
# stat=R → running; D → blocked on I/O (see §7.2); Z → zombie

# Step 4 — Load average vs CPU count
uptime && nproc
# load > nproc × 2 → run-queue congestion

# Step 5 — If sy% (kernel CPU) is high
dmesg -T | tail -20
perf record -g -p <PID> && perf report

# Step 6 — Per-process CPU samples over time (sysstat)
pidstat -u 1 5                    # 5 samples, 1s interval — install: apt/yum install sysstat

# Resolution: renice → cgroup throttle → kill -15 → kill -9 as last resort
```

### 7.2 High Memory / OOM

```bash
# Step 1 — Current pressure
free -h
# "available" < 10% of MemTotal → at risk of OOM

# Step 2 — Find memory hogs
ps aux --sort=-%mem | head -10
smem -r -s rss | head -10         # real RSS (no shared inflation) — install: apt/yum install smem

# Step 3 — Check if OOM already fired
dmesg -T | grep -i "oom-killer\|killed process"
journalctl -k | grep -i oom

# Step 4 — Check swap activity
swapon -s
vmstat 1 5 | awk 'NR>1 {print "si="$7" so="$8}'

# Step 5 — Resolution options
# a) Kill the hog gracefully first
kill -15 <hog-pid>
# b) Add swap (emergency)
dd if=/dev/zero of=/swapfile bs=1M count=2048 && mkswap /swapfile && swapon /swapfile
# c) Protect a critical process from OOM killer
echo -1000 > /proc/$(pgrep -x postgres)/oom_score_adj
# d) Long-term: reduce workload, tune app heap, or add RAM
```

### 7.3 Disk Full

```bash
# Step 1 — Find full filesystem (>80% used)
df -h | awk 'NR>1 && $5+0 > 80 {print}'

# Step 2 — Find space consumers
du -sh /var/log/* | sort -rh | head -10
find / -xdev -type f -size +500M 2>/dev/null | xargs ls -lh | sort -k5 -rh | head

# Step 3 — Check for deleted-but-still-open files (space not freed until restart)
lsof +L1 | grep deleted
# → restart the process holding the deleted fd to reclaim space

# Step 4 — Quick wins
journalctl --vacuum-time=3d        # trim systemd journals
apt-get clean                      # (Debian/Ubuntu) package cache
yum clean all                      # (RHEL/CentOS/Fedora) package cache
find /tmp -type f -atime +7 -delete
docker system prune                # if Docker present

# Step 5 — Prevent recurrence
# • logrotate config (see §4.5)
# • monitoring alert at 80% usage
```

### 7.4 Network Unreachable

```bash
# Layer-by-layer: work from L1 up until the break is found

# Step 1 — Physical/Link layer
ip link show                      # is the interface UP? errors/drops?
ethtool eth0                      # NIC link detected? speed/duplex?
# ethtool install: apt install ethtool  /  yum install ethtool

# Step 2 — IP/Routing layer
ip addr show                      # does the interface have an IP address?
ip route show                     # is there a default route (0.0.0.0/0)?
ping -c 2 8.8.8.8                 # can reach internet by IP (bypasses DNS)?

# Step 3 — DNS layer (only if step 2 passes)
ping -c 2 google.com              # fails but 8.8.8.8 pings → DNS issue
dig @8.8.8.8 google.com           # bypass local resolver
cat /etc/resolv.conf              # check configured nameserver

# Step 4 — Firewall
iptables -L -n | grep DROP        # check for blocking rules
curl -v telnet://host:port        # test specific port reachability

# Step 5 — Remote service
ss -tlnp | grep :<port>           # is the service actually listening?
systemctl status <service>        # is it running?
```

### 7.5 DNS Resolution Failure

```bash
# Step 1: Check resolv.conf
cat /etc/resolv.conf
# Step 2: Test with explicit server
dig @8.8.8.8 example.com         # bypass local resolver
dig @127.0.0.1 example.com       # test local resolver
# Step 3: Check systemd-resolved
systemctl status systemd-resolved
resolvectl status
resolvectl query example.com
# Step 4: Check /etc/nsswitch.conf
cat /etc/nsswitch.conf | grep hosts
# Should be: hosts: files dns myhostname
# Step 5: Flush DNS cache
sudo resolvectl flush-caches     # systemd-resolved (systemd ≥239, current standard)
# sudo systemd-resolve --flush-caches  # deprecated alias — avoid on modern systems
# macOS (the HUP does the DNS work; depth: references/macos-networking.md):
sudo dscacheutil -flushcache && sudo killall -HUP mDNSResponder
```

### 7.6 Port Already in Use

```bash
# Find what is using the port
ss -tlnp | grep :8080
lsof -i :8080
fuser 8080/tcp                    # show PID
fuser -k 8080/tcp                 # kill whatever is using it

# TIME_WAIT prevents immediate reuse — check
ss -tn state time-wait | grep :8080
# Fix: SO_REUSEADDR in application, or wait 2 × MSL (60-120s)
# Force reuse (kernel tuning, use with care):
sysctl -w net.ipv4.tcp_tw_reuse=1
```

### 7.7 Service Won't Start

```bash
# Step 1: Check status and last lines
systemctl status myservice -l     # -l for full output
# Step 2: Check journal for errors
journalctl -u myservice -b --no-pager | tail -50
# Step 3: Common causes:
#   - Configuration syntax error: nginx -t, apachectl configtest
#   - Port in use: ss -tlnp | grep :<configured-port>
#   - Missing files: verify ExecStart path exists and is executable
#   - Permission denied: check User= in unit file, file permissions
#   - Missing dependency: systemctl list-dependencies myservice
# Step 4: Run ExecStart manually as same user
sudo -u www-data /usr/local/bin/myapp --config /etc/myapp/config.yaml
# Step 5: Check SELinux/AppArmor
getenforce                        # SELinux mode
ausearch -m avc -ts recent        # SELinux denials
aa-status                         # AppArmor status
```

### 7.8 Slow SSH

```bash
# On client — add timing verbosity
ssh -v user@host
ssh -o ConnectTimeout=5 user@host

# Common causes and fixes:
# 1. Reverse DNS lookup: UseDNS no  (in /etc/ssh/sshd_config)
# 2. GSSAPI probing: GSSAPIAuthentication no (client: ~/.ssh/config)
# 3. Waiting for entropy: haveged or rng-tools on server
# 4. Key exchange: try ssh -o KexAlgorithms=ecdh-sha2-nistp256 user@host

# On server
grep -E "UseDNS|GSSAPIAuthentication" /etc/ssh/sshd_config
journalctl -u sshd -f &           # watch SSH logs while connecting
```

### 7.9 Certificate Errors

```bash
# Inspect certificate
openssl s_client -connect host:443 </dev/null 2>&1
echo | openssl s_client -servername host -connect host:443 2>/dev/null | openssl x509 -text -noout

# Check expiry
echo | openssl s_client -connect host:443 2>/dev/null \
  | openssl x509 -noout -dates

# Check chain
openssl s_client -connect host:443 -showcerts </dev/null 2>/dev/null

# Common issues:
# - Expired: check notAfter date
# - Self-signed: -k flag in curl to skip verify; or install CA cert
# - Hostname mismatch: verify CN and SANs vs hostname used
# - Incomplete chain: intermediate certificate missing from server config
# - Wrong CA bundle: curl --cacert /path/to/ca-bundle.crt URL

# Verify local cert file
openssl verify -CAfile /etc/ssl/certs/ca-certificates.crt /path/to/cert.pem
openssl x509 -noout -subject -issuer -dates -in /path/to/cert.pem
```

---

## Anti-Patterns

- **Using SIGKILL as first response** — always try SIGTERM first; SIGKILL can leave resources locked or corrupt state.
- **Editing /etc/resolv.conf directly** on systemd-resolved systems — it gets overwritten; use `resolvectl` or NetworkManager.
- **Using device names (/dev/sda1) in /etc/fstab** — device assignment can change on reboot; always use UUID.
- **Running `rm -rf` to free disk space** without checking `lsof +L1` first — deleted files may still be held open; space won't be freed until the process restarts.
- **Disabling SELinux/AppArmor** to "fix" permission problems — diagnose and write a policy instead.
- **`chmod 777` on directories** — use group permissions and ACLs for sharing; 777 removes all security boundaries.
- **Modifying /etc/sudoers directly** instead of `visudo` — a syntax error locks out all sudo access.
- **Ignoring `vmstat si/so`** — swap I/O is a critical memory pressure signal; even small values indicate exhaustion.
- **Parsing `ip addr` with fragile regex** — use `ip -j addr | jq` for scripting to get stable JSON output.
- **Not using `persistent=true`** on systemd timers that must not miss runs — without it, a reboot silently skips a missed window.
- **Running `kill -9` on a D-state process** — it won't work; the process is in uninterruptible kernel wait; investigate the I/O or storage layer.

---

## References

- [Arch Wiki: Network tools](https://wiki.archlinux.org/title/Network_tools)
- [Arch Wiki: Network configuration](https://wiki.archlinux.org/title/Network_configuration)
- [Red Hat RHEL 9: Monitoring and managing system status and performance](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/9/html/monitoring_and_managing_system_status_and_performance/index)
- [Red Hat RHEL 8: System Administrator's Guide — System Monitoring Tools](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/system_administrators_guide/ch-system_monitoring_tools)
- [DigitalOcean: How To Use journalctl to View and Manipulate systemd Logs](https://www.digitalocean.com/community/tutorials/how-to-use-journalctl-to-view-and-manipulate-systemd-logs)
- [DigitalOcean: How To Partition and Format Storage Devices in Linux](https://www.digitalocean.com/community/tutorials/how-to-partition-and-format-storage-devices-in-linux)
- [Cherry Servers: Top 20 Linux Network Commands](https://www.cherryservers.com/blog/linux-network-commands)
- [DevOpsCube: Linux Networking and Troubleshooting Commands](https://devopscube.com/list-linux-networking-troubleshooting-and-commands-beginners/)
- [Linux Journal: Linux Performance Monitoring: top, vmstat, and iostat](https://www.linuxjournal.com/content/linux-performance-monitoring-using-tools-top-vmstat-and-iostat)
- [Baeldung: ufw vs nftables vs iptables comparison](https://www.baeldung.com/linux/ufw-nftables-iptables-comparison)
- [CB.VU: Linux Firewall with iptables and nftables](https://cb.vu/networking/firewall-guide)
- [Last9: Understanding the Linux OOM Killer](https://last9.io/blog/understanding-the-linux-oom-killer/)
- [NetworkersHome: Linux Disk Management: LVM, RAID, Partitions & Filesystems](https://www.networkershome.com/fundamentals/linux/linux-disk-management-lvm-raid-partitions/)
- [osxhub: macOS Network Diagnostic and Configuration Commands](https://osxhub.com/macos-network-diagnostic-commands-guide/)
- [osxhub: Mastering macOS Process Management](https://osxhub.com/macos-process-management-ps-kill-launchctl-guide/)
- [LinuxBlog: Linux Disk Partitioning with fdisk, parted, and lsblk](https://linuxblog.io/linux-disk-partitioning-fdisk-parted-lsblk/)
- [LinuxBlog: Linux File Permissions Explained](https://linuxblog.io/linux-file-permissions-explained-chmod-chown-and-umask-in-practice/)
- [Red Hat Sysadmin: How to manage Linux permissions for users, groups, and others](https://www.redhat.com/sysadmin/manage-permissions)
- [scutil macOS Man Page / SCUtil Guide](https://cli.wiki/scutil---macOS-System-Configuration-Utility-Guide)
- [oneuptime: How to Troubleshoot DNS Resolution Failures on Linux](https://oneuptime.com/blog/post/2026-03-20-troubleshoot-dns-resolution-failures/view)
- [CB.VU: Systemd: Services, Timers, and Journal](https://cb.vu/linux/systemd-guide)
- [Jeff Geerling: Partition, format, and mount a large disk in Linux with parted](https://www.jeffgeerling.com/blog/2021/htgwa-partition-format-and-mount-large-disk-linux-parted/)

# check_akcp_sensorprobeXplus

Check plugin to query sensor data from AKCP sensor products via SNMP.

# Building

## Dependencies
There are some go dependencies, but go will fetch them for building dynamically from the internet.

Debian/Ubuntu
	apt-get install golang

## Compiling

Debian/Ubuntu
```
git clone $THIS_REPO
cd $REPO
go build
```

# Installation

The compiled binary is completely standalone, so copy it to place where you want to have it. In the case of using Icinga2 probably `/usr/lib/nagios/plugins/` (on Debian/Ubuntu).

# Usage
Executing the binary without any parameters will present the available parameters. The simplest form to use it productively would be to give only the _network address_ of the probe.
If not other paremeter is give, default settings will be used.

# Examples
```
check_apck -h 192.168.1.1

WARNING - states: unknown=1 warning=1 ok=3
[OK] Temperature Port 1: 27
[WARNING] Dual Humidity Port 2: 38%
[OK] Dual Temperature Port 2: 27
[OK] Airflow Port 3: 0%
```

# check_akcp_sensorprobeXplus

Check plugin to query sensor data from AKCP sensorProbeX+ via SNMP.

# Building

## Dependencies
There are some go dependencies, but go will fetch them for building dynamically from the internet.
Additionaly to the golang toolchain you might want to install git.

Debian/Ubuntu

	apt-get install golang

CentOS/RHEL

	yum install epel-release
	yum install golang


## Compiling

Debian/Ubuntu/CentOS/RHEL
```
git clone https://github.com/NETWAYS/check_akcp_sensorprobeXplus
cd check_akcp_sensorprobeXplus
go build
```

# Installation

The compiled binary is completely standalone, so copy it to place where you want to have it. In the case of using Icinga2 probably `/usr/lib/nagios/plugins/` (on Debian/Ubuntu).

	cp check_akcp_sensorprobeXplus /usr/lib/nagios/plugins/

# Usage
Executing the binary without any parameters will present the available parameters. The simplest form to use it productively would be to give only the _network address_ of the probe.
If not other paremeter is give, default settings will be used.

# Examples
```
check_akcp_sensorprobeXplus -h 192.168.1.1

WARNING - states: unknown=1 warning=1 ok=3
[OK] Temperature Port 1: 27
[WARNING] Dual Humidity Port 2: 38%
[OK] Dual Temperature Port 2: 27
[OK] Airflow Port 3: 0%
```

```
./check_akcp_sensorprobeXplus -h 192.168.1.1 -c public
WARNING - Device SPX+ Demo at location Room 217 (SPX+ F7 1.0.5233 May 12 2020 09:41:)

[WARNING] Dual Humidity Port 1: 25.0% is lower than warning threshold 32.0%
[WARNING] Dual Temperature Port 1: 31.1℃ is higher than warning threshold 30.0℃
[OK] Temp Sensor Test: 26.1℃
[OK] Motion Detector Port 4: Erkannt  | 'Dual Humidity Port 1'=25%;32:66;23:69 'Dual Temperature Port 1'=31.1C;20.7:30;10.6:40 'Temp Sensor Test'=26.1C;20.8:30;10.5:40 'Motion Detector Port 4'=0
```

## License

Copyright (c) 2022 [NETWAYS GmbH](mailto:info@netways.de) \
Copyright (c) 2022 [Lorenz Kästle](mailto:lorenz.kaestle@netways.de)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see [gnu.org/licenses](https://www.gnu.org/licenses/).


# GEOWall

Downloaded and updates IPtables based on Geo IP lists and country.

Works with both 'block country' and 'allow only' country modes.

IP block rules from:
- IPv4: http://www.ipdeny.com/ipblocks/data/aggregated/
- IPv6: http://www.ipdeny.com/ipv6/ipaddresses/aggregated/

Simply match the 'country code' to the country characters at the start of the {countrycode}-aggredated.zone file

For example, to only allow AU IP's to the server on eth0 and block everything else, use:

```
geowall -4 -6 -i eth0 start -c au -m allow
```

This will download and process:
http://www.ipdeny.com/ipblocks/data/aggregated/au-aggregated.zone to IPTables
http://www.ipdeny.com/ipv6/ipaddresses/aggregated/au-aggregated.zone to IP6Tables

To unload IPv4 and IPv6 rules on eth0 run:

```
geowall -4 -6 -i eth0 stop
```

To allow only AU and NZ IP's and block everything else on both IPv4 and IPv6, use:

```
geowall -4 -6 -i eth0 start -c au,nz -m allow
```

To block a specific country on both IPv4 and IPv6, use:

```
geowall -4 -6 -i eth0 start -c {countrycode} -m deny
```

To update rules for IPv4 to use AU and NZ use - this will clear existing rules and replace with AU and NZ:

```
geowall -4 update -c au,nz
```

To build this application:
Requires Golang 1.11.1

```
go get -u github.com/digitalsparky/geowall
```

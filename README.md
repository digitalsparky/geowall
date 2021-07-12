# GEOWall [!Build](https://github.com/digitalsparky/geowall/actions/workflows/main.yml/badge.svg?event=release)](https://github.com/digitalsparky/geowall/actions/workflows/main.yml)

Downloaded and updates IPtables based on Geo IP lists and country.

Works with both 'block country' and 'allow only' country modes.

IP block rules from:
- IPv4: http://www.ipdeny.com/ipblocks/data/aggregated/
- IPv6: http://www.ipdeny.com/ipv6/ipaddresses/aggregated/

Simply match the 'country code' to the country characters at the start of the {countrycode}-aggredated.zone file

```
NAME:
   geowall - GeoIP Based Firewall

USAGE:
   main [global options] command [command options] [arguments...]

AUTHOR:
   Matt Spurrier <matthew@spurrier.com.au>

COMMANDS:
     apply
     unload
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --v4, -4, --ipv4         Include IPv4 Rules
   --v6, -6, --ipv6         Include IPv6 Rules
   --iface value, -i value  Inbound Interface
   --save, -s               Save IPTables Rules
   --help, -h               show help
   --version, -v            print the version

COPYRIGHT:
   (c) 2019 Matt Spurrier
```

For example, to only allow AU IP's to the server on eth0 and block everything else, use:

```
geowall -4 -6 -i eth0 apply -c au -m allow
```

This will download and process:
http://www.ipdeny.com/ipblocks/data/aggregated/au-aggregated.zone to IPTables
http://www.ipdeny.com/ipv6/ipaddresses/aggregated/au-aggregated.zone to IP6Tables

To unload IPv4 and IPv6 rules on eth0 run:

```
geowall -4 -6 -i eth0 unload
```

To allow only AU and NZ IP's and block everything else on both IPv4 and IPv6, use:

```
geowall -4 -6 -i eth0 apply -c au,nz -m allow
```

To block a specific country on both IPv4 and IPv6, use:

```
geowall -4 -6 -i eth0 apply -c {countrycode} -m deny
```

To update rules for IPv4 to use AU and NZ use - this will clear existing rules and replace with AU and NZ:

```
geowall -4 -i eth0 apply -c au,nz
```

To allow only AU and NZ IP's and block everything else on both IPv4 and IPv6 and save the rules use:

```
geowall -4 -6 -i eth0 -s apply -c au,nz -m allow
```

To build this application:
Requires Golang 1.11.1

```
go get -u github.com/digitalsparky/geowall
```

# Like my stuff?

Would you like to buy me a coffee or send me a tip?
While it's not expected, I would really appreciate it.

[![Paypal](https://www.paypalobjects.com/webstatic/mktg/Logo/pp-logo-100px.png)](https://paypal.me/MattSpurrier) <a href="https://www.buymeacoffee.com/digitalsparky" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/white_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

# acrobits-sipis-exporter

Simple prometheus exporter to query your sipis instance for statistics and present it in prometheus format

Usage:

```
  $ go build
  $ acrobits-sipis-exporter --listen-address :8080 --instance https://sipis.example.com --every 15m
```

All command line flags are read from ENV as well: `SIPIS_LISTEN_ADDRESS`, `SIPIS_INSTANCE` and `SIPIS_EVERY`.

# acrobits-sipis-exporter 🚀

Simple Prometheus exporter to query your SIPIS instance for statistics and present it in Prometheus format 📊

## Usage:

```sh
  $ go build
  $ ./acrobits-sipis-exporter --listen-address :8080 --instance https://sipis.example.com --every 15m
```

You can also set multiple instances via environment variables:

```sh
  $ export SIPIS_INSTANCES="https://sipis1.example.com,https://sipis2.example.com"
  $ ./acrobits-sipis-exporter --listen-address :8080 --every 15m
```

All command line flags are read from ENV as well: `SIPIS_LISTEN_ADDRESS`, `SIPIS_INSTANCE`, and `SIPIS_EVERY`.

## Flags and Environment Variables

- `--listen-address` (or `SIPIS_LISTEN_ADDRESS`): The address to listen on for HTTP requests. Default is `:8080`.
- `--instance` (or `SIPIS_INSTANCE`): The URL of the SIPIS instance to query. Can be specified multiple times.
- `--every` (or `SIPIS_EVERY`): The interval at which to query the SIPIS instance. Default is `15m`.

## Example

```sh
  $ go build
  $ export SIPIS_INSTANCES="https://sipis1.example.com,https://sipis2.example.com"
  $ ./acrobits-sipis-exporter --listen-address :8080 --every 15m
```

This will start the exporter and query the specified SIPIS instances every 15 minutes. The metrics will be available at `http://localhost:8080/metrics` 🎉

Happy monitoring! 📈

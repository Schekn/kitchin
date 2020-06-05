# Delivery system simulation

Real-time system that emulates the fulfillment of delivery orders for a kitchen.

Run the app:
```bash
$ ./build/delivery -o orders.json
# or
$ ./build/delivery -o orders.json -c /path/to/config.yml
```

Type `p+Enter` to pause execution and `c+Enter` to continue.

Сonfiguration can be changed in the file `config.yml`. There can be configured list of shelves, orders and courier parameters. Note that valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.

For more convenient use, use the `Makefile`. To see all available commands run `make help`:
```bash
$ make help
build                          Build project
dev                            Install dev tools
doc                            Show documentation
help                           Display callable targets
lint                           Lint code
run                            Run application
test                           Run tests
```
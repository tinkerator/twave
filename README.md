# twave - a tool to view a dump.vcd file as text

## Overview

The `twave` tool is a lightweight program to view a `dump.vcd` file in
ASCII. I've used it for debugging in constrined environments. If
moving the `dump.vcd` file around is straightforward, you are much
better off using something more sophisticated like
[GTKWave](https://gtkwave.github.io/gtkwave/).

## Getting started

Build from source:
```
$ git clone https://github.com/tinkerator/twave.git
$ cd twave
$ go build twave.go
```

You can install it in your `~/go/bin/` directory with:
```
$ go install twave.go
```

You can invoke `./twave` against a VCD ([Value Change
Dump](https://en.wikipedia.org/wiki/Value_change_dump)) file as
follows:
```
$ ./twave --file=v/dump.vcd
[] : [$version Icarus Verilog $end]
           led_driver_test.reset__ack-+
                  led_driver_test.led-|-+
                  led_driver_test.clk-|-|-+
       led_driver_test.reset__disable-|-|-|-+
         led_driver_test.script[11:0]-|-|-|-|------------+
 led_driver_test.target.counter[11:0]-|-|-|-|------------|------------+
  led_driver_test.target.thresh[11:0]-|-|-|-|------------|------------|------------+
                                      | | | |            |            |            |
     2022-12-31 14:24:11.000000000000 x x x x xxxxxxxxxxxx xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000010000 x x 0 1 000000000000 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000020000 x x 0 1 000000000001 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000030000 x x 1 1 000000000010 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000040000 0 x 1 1 000000000011 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000050000 0 x 0 1 000000000100 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000060000 0 x 0 1 000000000101 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000070000 0 x 1 1 000000000110 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000080000 0 x 1 0 000000000111 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000090000 0 x 0 0 000000001000 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000100000 0 x 0 0 000000001001 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000110000 0 x 1 0 000000001010 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000120000 1 0 1 0 000000001011 000000000000 000000000001
     2022-12-31 14:24:11.000000130000 1 0 0 0 000000001100 000000000000 000000000001
     2022-12-31 14:24:11.000000140000 1 0 0 0 000000001101 000000000000 000000000001
     2022-12-31 14:24:11.000000150000 1 0 1 0 000000001110 000000000000 000000000001
     2022-12-31 14:24:11.000000160000 1 0 1 0 000000001111 000000000000 000000000001
     2022-12-31 14:24:11.000000170000 1 0 0 0 000000010000 000000000000 000000000001
     2022-12-31 14:24:11.000000180000 1 0 0 0 000000010001 000000000000 000000000001
     2022-12-31 14:24:11.000000190000 1 0 1 0 000000010010 000000000000 000000000001
     2022-12-31 14:24:11.000000200000 1 0 1 0 000000010011 000000000000 000000000001
     [... truncated ...]
```

The sample `v/dump.vcd` file can be regenerated using the Public
Domain verilog code we've included in the [`v/`](v) directory.

Note, `dump.vcd` files can contain a lot of state, which can cause the
output of `twave` to format poorly on a limited size terminal. The
`twave` program supports a `--syms` argument to limit the output to
specific symbol values only.

As with any command line tool that outputs text, you can combine
`twave` with tools like `grep`, `sed` and `awk` to quickly find
entries of interest. For example, the list of symbols in a dump can be
found with:
```
$ ./twave --file=v/dump.vcd | grep -i -E '^ *[a-z]'
```

## License info

The `twave` program is distributed with the same BSD 3-clause license
as that used by [golang](https://golang.org/LICENSE) itself.

## Reporting bugs and feature requests

The program `twave` has been developed purely out of self-interest and
a curiosity for debugging using command line programs only. Should you
find a bug or want to suggest a feature addition, please use the [bug
tracker](https://github.com/tinkerator/twave/issues).

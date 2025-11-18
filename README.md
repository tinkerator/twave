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
                  led_driver_test.clk-+
                  led_driver_test.led-|-+
           led_driver_test.reset__ack-|-|-+
       led_driver_test.reset__disable-|-|-|-+
         led_driver_test.script[11:0]-|-|-|-|------------+
 led_driver_test.target.counter[11:0]-|-|-|-|------------|------------+
  led_driver_test.target.thresh[11:0]-|-|-|-|------------|------------|------------+
                                      | | | |            |            |            |
     2022-12-31 14:24:11.000000000000 x x x x xxxxxxxxxxxx xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000010000 0 x x 1 000000000000 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000020000 0 x x 1 000000000001 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000030000 1 x x 1 000000000010 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000040000 1 x 0 1 000000000011 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000050000 0 x 0 1 000000000100 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000060000 0 x 0 1 000000000101 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000070000 1 x 0 1 000000000110 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000080000 1 x 0 0 000000000111 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000090000 0 x 0 0 000000001000 xxxxxxxxxxxx xxxxxxxxxxxx
     2022-12-31 14:24:11.000000100000 0 x 0 0 000000001001 xxxxxxxxxxxx xxxxxxxxxxxx
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

The `twave` program supports other output modes. The `--raw` mode
displays the integer time offset without rendering it in the form of a
clock timestamp. This can be useful in conjuction with the `--start`
and `--end` arguments. For example:

```
$ ./twave --raw --file=v/dump.vcd | tail -10
                                29200 1 0 0 1 000100100011 000000001011 000000010011
                                29300 0 0 0 1 000100100100 000000001011 000000010011
                                29400 0 0 0 1 000100100101 000000001011 000000010011
                                29500 1 0 0 1 000100100110 000000001011 000000010011
                                29600 1 0 0 1 000100100111 000000001100 000000010011
                                29700 0 0 0 1 000100101000 000000001100 000000010011
                                29800 0 0 0 1 000100101001 000000001100 000000010011
                                29900 1 0 0 1 000100101010 000000001100 000000010011
                                30000 1 0 0 1 000100101011 000000001101 000000010011
```

Which, instead of timestamps, shows us values that can be used for
`--start` and `--end` arguments. That is, to only display the last 5
entries, try this:

```
$ ./twave --file=v/dump.vcd --start=29600 --end=30000
[] : [$version Icarus Verilog $end]
                  led_driver_test.clk-+
                  led_driver_test.led-|-+
           led_driver_test.reset__ack-|-|-+
       led_driver_test.reset__disable-|-|-|-+
         led_driver_test.script[11:0]-|-|-|-|------------+
 led_driver_test.target.counter[11:0]-|-|-|-|------------|------------+
  led_driver_test.target.thresh[11:0]-|-|-|-|------------|------------|------------+
                                      | | | |            |            |            |
     2022-12-31 14:24:11.000002960000 1 0 0 1 000100100111 000000001100 000000010011
     2022-12-31 14:24:11.000002970000 0 0 0 1 000100101000 000000001100 000000010011
     2022-12-31 14:24:11.000002980000 0 0 0 1 000100101001 000000001100 000000010011
     2022-12-31 14:24:11.000002990000 1 0 0 1 000100101010 000000001100 000000010011
     2022-12-31 14:24:11.000003000000 1 0 0 1 000100101011 000000001101 000000010011
```

Partial support for [`wavy`](https://zappem.net/pub/project/wavy/)
text format is imlpemented in `twave`. To convert the `v/dump.vcd` file to a PNG image`:

```
$ ./twave --file=v/dump.vcd --wavy > led.wvy
$ ~/go/bin/wavy --input led.wvy --output led.png
```

Which generates the following:

![Generated led.png image from generated led.wvy](led.png "led.wvy converted to a PNG")

You can combine these flags with an explicit list of ordered symbols,
if you want to focus on a smaller number of signals over a specific
range:

```
$ ./twave --file=v/dump.vcd --wavy --start=10000 --end=20000 --syms led_driver_test.clk,led_driver_test.led > led-trunc.wvy
$ ~/go/bin/wavy --input led-trunc.wvy --output led-trunc.png
```

![Generated led-trunc.png image from generated led-trunc.wvy](led-trunc.png "led-trunc.wvy converted to a PNG")

## TODOs

- the `--wavy` output format does not support multibit signals yet.

## License info

The `twave` program is distributed with the same BSD 3-clause license
as that used by [golang](https://golang.org/LICENSE) itself.

## Reporting bugs and feature requests

The program `twave` has been developed purely out of self-interest and
a curiosity for debugging using command line programs only. Should you
find a bug or want to suggest a feature addition, please use the [bug
tracker](https://github.com/tinkerator/twave/issues).

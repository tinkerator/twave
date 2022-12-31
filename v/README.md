# led_driver

## Overview

This is some simple verilog code (`led_driver.v`) with a test wrapper
that drives it (`led_driver_test.v`) and pushes it through its
paces.

You can compile and run it with
[iverilog](http://iverilog.icarus.com/) as follows:
```
$ iverilog led_driver_test.v led_driver.v
$ ./a.out
VCD info: dumpfile dump.vcd opened for output.
led.bb(56):RESULT=PASS:0 @ testing completed
```
As indicated, this generates a `dump.vcd` file.

A file like this can be used as input to `twave`. We've included a
sample `dump.vcd` file in the git repository to act as an example
input for `twave`.

## About the verilog code

This code drives an output (`led`) at a slower and slower rate over
time. Eventually, when the threshold counter overflows, things speed
up again suddenly but then start to slow down again. Not very useful,
but the purpose is to generate a non-trivial `dump.vcd` file.

## License

The verilog code in this present directory is considered to be in the
Public Domain.

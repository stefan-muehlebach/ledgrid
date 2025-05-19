# LedGrid - Controlling NeoPixels with Go

The package `ledgrid` contains a number of functions and types, which allows
you to control a chain of NeoPixels, lets you draw shapes, images and custom
patterns and provides you with methods and types to animate most of your stuff.
Beside the library itself, this package also contains the following command
line tools and examples:

* `cmd/anim`: A collection of almost all visual effects, animations, shapes, etc.
  this package has to offer. The frontend is a fairly simple interactive
  command line tool.
* `cmd/colorEdit`: A ncurses application which lets you play with the NeoPixels
  interactively. It gives you direct access to the RGB-values of all LEDs and
  provides you with a small set of functions, like creating a gradient between
  two given LEDs.
* `cmd/gridController`: A daemon process which runs on the remote system the
  NeoPixels are connected to.
* `cmd/gridEmulator`: A GUI tool which lets you develop animations _without_
  the need of the NeoPixel hardware.


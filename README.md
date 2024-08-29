# Notify

Displays text read from stdin in a pop-up notification window.

Minimal and daemonless alternative to notification services like [Dunst](https://github.com/dunst-project/dunst).

## Features

* *Customizability* Set dimensions, placement, font, duration, borderwidth, bordercolor, bgcolor and fgcolor via command-line arguments.
* *Scripting* The notification text is read through stdin; Set the stdout text via a command-line argument; Control the exit code via left and right mousebutton clicks on the notification window.

## Build

```sh
$ make
# or...
$ make install
```

## Usage

See `$ notify --help` for command-line options

The notification text can be `[Title]Body` (i.e. title text in brackets and body text after) or `Body` (i.e. just body text).

```sh
$ curl example.com >/dev/null 2>&1 && \
  notify -B "#28A745" -d 1s <<< "[Curl]Download succeeded." || \
  notify -B "#DC3545" -d 1s <<< "[Curl]Download failed."
```
![Screenshot](screenshot.png)

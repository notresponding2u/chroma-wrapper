# Chroma heatmap
## Disclaimer
This project is developed and tested on Razer Black Widow chroma X, on Windows 10 x64.
There are bugs:
- The library used for detecting the keystrokes, does not detect the Fn key.
- Works only on Windows 10.
- No exit, except ctrl+c
- You need to install Connect from Razer synapse to be able to connect to the SDK.
- The library with the hook is throwing some error, because of US keyboard layout. But still works fine for the effect.

I will add some other functionalities later.

## Interface

F12 to save exit.

F11 to load all time build up heatmap.

F10 to save the map and start new one.

F9 discard current map.

## Build

To build as background process use this:

```shell
go build -ldflags -H=windowsgui
```
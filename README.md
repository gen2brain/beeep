## beeep
[![Build Status](https://github.com/gen2brain/beeep/actions/workflows/build.yml/badge.svg)](https://github.com/gen2brain/beeep/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/gen2brain/beeep.svg)](https://pkg.go.dev/github.com/gen2brain/beeep)

`beeep` provides a cross-platform library for sending desktop notifications, alerts and beeps.

### Installation

    go get -u github.com/gen2brain/beeep

### Build tags

* `nodbus` - disable `godbus/dbus` and use only `notify-send`

### Examples

```go
err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
if err != nil {
    panic(err)
}
```

```go
//go:embed testdata/info.png
var icon []byte

err := beeep.Notify("Title", "Message body", icon)
if err != nil {
    panic(err)
}
```

```go
beeep.AppName = "My App Name"

err := beeep.Alert("Title", "Message body", "testdata/warning.png")
if err != nil {
    panic(err)
}
```

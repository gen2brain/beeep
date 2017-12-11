// +build js

package beeep

import (
	"github.com/gopherjs/gopherjs/js"
)

// Notify sends desktop notification.
func Notify(title, message string) (err error) {
	defer func() {
		e := recover()

		if e == nil {
			return
		}

		if e, ok := e.(*js.Error); ok {
			err = e
		} else {
			panic(e)
		}
	}()

	n := js.Global.Get("window").Get("Notification")

	if n.Get("permission").String() == "granted" {
		n.New(title, map[string]interface{}{
			"body": message,
		})
	} else {
		n.Call("requestPermission", func(permission string) {
			if permission == "granted" {
				n.New(title, map[string]interface{}{
					"body": message,
				})
			}
		})
	}

	return
}

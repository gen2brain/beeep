package beeep

func ExampleBeep() {
	Beep(DefaultFreq, DefaultDuration)
}

func ExampleNotify() {
	Notify("Title", "MessageBody", "assets/icon128.png")
}

func ExampleNotify() {
	Alert("Title", "MessageBody", "assets/icon128.png")
}

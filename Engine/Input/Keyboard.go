package Input

var (
	keyState   = make(map[int]int8)
	mouseState = make(map[int]int8)
)

func OnKey(key, state int) {
	switch key {
	case KeyEsc:
		//running = false
	}
	switch state {
	case Key_Release:
		keyState[key] &= 2
	case Key_Press:
		if keyState[key] == 0 {
			keyState[key] = 3
		} else {
			keyState[key] |= 1
		}
	}
}

func UpdateInput() {
	for i, v := range keyState {
		keyState[i] = v & ^2
	}
	for i, v := range mouseState {
		mouseState[i] = v & ^2
	}
}

func KeyDown(key int) bool {
	return keyState[key]&1 != 0
}

func KeyUp(key int) bool {
	return keyState[key]&1 == 0
}

func KeyPress(key int) bool {
	return keyState[key]&2 != 0
}

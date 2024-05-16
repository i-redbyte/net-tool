package spinner

import (
	"fmt"
	"time"
)

func Spinner(done chan bool, endMessage string) {
	symbols := []rune{'\\', '|', '/', '-'}
	i := 0
	fmt.Print("Waiting for results ")
	for {
		select {
		case <-done:
			fmt.Printf("\r%s                       \r", endMessage)
			return
		default:
			fmt.Printf("\rWaiting for results %c", symbols[i])
			i = (i + 1) % len(symbols)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

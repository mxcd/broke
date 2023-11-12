package util

import (
	"fmt"
	"net/http"
	"time"
)

func WaitForServerUp(url string) {
	for {
		resp, err := http.Get(url)
		if err == nil && resp != nil && resp.StatusCode == 200 {
			fmt.Println("Server is up!")
			return
		}

		if err != nil {
			fmt.Printf("Failed to reach server: %s. Retrying...\n", err.Error())
		} else if resp != nil {
			fmt.Printf("Received status code %d. Retrying...\n", resp.StatusCode)
		}

		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(5 * time.Millisecond)
	}
}

func ListLimitOffset[T any](data []T, limit int, offset int) []T {
	if offset >= len(data) {
		return []T{}
	}
	if offset+limit > len(data) {
		return data[offset:]
	}
	return data[offset : offset+limit]
}

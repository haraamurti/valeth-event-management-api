package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

func main() {
    var wg sync.WaitGroup
    url := "http://localhost:8080/events/1/reserve"
    concurrent := 200
    wg.Add(concurrent)
    client := &http.Client{}
    for i:=0;i<concurrent;i++ {
        go func(i int) {
            defer wg.Done()
            body := fmt.Sprintf(`{"user_id":"user-%d"}`, i)
            req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
            req.Header.Set("Content-Type","application/json")
            resp, err := client.Do(req)
            if err != nil {
                fmt.Println("err:", err)
                return
            }
            defer resp.Body.Close()
            fmt.Println("status", resp.Status)
        }(i)
    }
    wg.Wait()
}

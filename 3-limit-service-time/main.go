//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"log"
	"sync/atomic"
	"time"
)

const freeDuration = 10 * time.Second

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

func (u *User) AddTime(consumed int64) int64 {
	return atomic.AddInt64(&u.TimeUsed, consumed)
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	result := make(chan time.Duration)
	timeCh := time.After(freeDuration - time.Duration(time.Duration(u.TimeUsed).Seconds()))
	log.Println(freeDuration - time.Duration(time.Duration(u.TimeUsed).Seconds()))

	go func() {
		start := time.Now()
		process()
		result <- time.Since(start)
	}()

	if u.IsPremium {
		<-result
		return true
	}

	select {
	case timeConsumed := <-result:
		u.TimeUsed += int64(timeConsumed.Seconds())
		return true
	case <-timeCh:
		u.TimeUsed = int64(freeDuration)
		return false
	}
}

func main() {
	RunMockServer()
}

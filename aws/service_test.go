package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"testing"
	"time"
)

func TestSessionLimiter_Concurrent(t *testing.T) {
	type args struct {
		sess          *session.Session
		maxConcurrent *int
		maxPerSecond  *int
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		Handlers: request.Handlers{
			Validate: request.HandlerList{},
			Complete: request.HandlerList{},
		},
	})
	if err != nil {
		t.Errorf("failed to create session: %s", err)
	}
	tests := []struct {
		name    string
		args    args
		runLoop int
		nextShouldBlock bool
	}{
		{
			"second request blocks when maxConcurrent is set to 1",
			args{
				sess.Copy(),
				aws.Int(1),
				aws.Int(100),
			},
			1,
			true,
		},
		{
			"fifth request blocks when maxConcurrent is set to 4",
			args{
				sess.Copy(),
				aws.Int(4),
				aws.Int(100),
			},
			4,
			true,
		},
		{
			"fifty requests should run without blocking when maxConcurrent is set to 0",
			args{
				sess.Copy(),
				aws.Int(0),
				aws.Int(100),
			},
			50,
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SessionLimiter(tt.args.sess, tt.args.maxConcurrent, tt.args.maxPerSecond)
			go func(){
				time.Sleep(10 * time.Second)
				t.Errorf("timed out, this is most likely caused by the handler unexpectedly blocking")
				t.Fail()
			}()

			for i := 0; i < tt.runLoop; i++ {
				tt.args.sess.Handlers.Validate.Run(&request.Request{})
			}

			done := make(chan int, 1)
			running := true
			// The next call to Validate should block, if it doesn't fail.
			go func() {
				tt.args.sess.Handlers.Validate.Run(&request.Request{})
				if running && tt.nextShouldBlock {
					t.Failed()
				}
				done <- 0
			}()

			// Give the goroutine a chance to run, if t.Failed() hasn't been called in this time we
			// consider the test passed, free up a slot, then wait for the goroutine to finish.
			time.Sleep(time.Millisecond * 100)
			running = false // prevents t.Failed() from being called in goroutine
			tt.args.sess.Handlers.Complete.Run(&request.Request{})
			<- done
		})
	}
}


func TestSessionLimiter_PerSecond(t *testing.T) {
	type args struct {
		sess          *session.Session
		maxConcurrent *int
		maxPerSecond  *int
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		Handlers: request.Handlers{
			Validate: request.HandlerList{},
			Complete: request.HandlerList{},
		},
	})
	if err != nil {
		t.Errorf("failed to create session: %s", err)
	}
	tests := []struct {
		name           string
		args           args
		maxReqsASecond int
	}{
		{
			"no more then two requests are run in a second when maxConcurrent is 1 and maxPerSecond is 1",
			args{
				sess.Copy(),
				aws.Int(1),
				aws.Int(1),
			},
			2,
		},
		{
			"no more then two requests are run in a second when maxConcurrent is 100 and maxPerSecond is 1",
			args{
				sess.Copy(),
				aws.Int(100),
				aws.Int(1),
			},
			2,
		},
		{
			"no more then 5 requests are run in a second when maxConcurrent is 1 and maxPerSecond is 4",
			args{
				sess.Copy(),
				aws.Int(1),
				aws.Int(4),
			},
			5,
		},
		{
			"no more then 12 requests are run in a second when maxConcurrent is 1 and maxPerSecond is 10",
			args{
				sess.Copy(),
				aws.Int(1),
				aws.Int(10),
			},
			12,
		},
		{
			"no more then 23 requests are run in a second when maxConcurrent is 1 and maxPerSecond is 20",
			args{
				sess.Copy(),
				aws.Int(1),
				aws.Int(20),
			},
			23,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SessionLimiter(tt.args.sess, tt.args.maxConcurrent, tt.args.maxPerSecond)

			go func(){
				time.Sleep(10 * time.Second)
				t.Errorf("timed out, this is most likely caused by the handler unexpectedly blocking")
				t.Fail()
			}()

			// Check the number of loops executed in the goroutine after a second, it should
			// not exceed maxReqsASecond.
			var i int
			done := false
			go func() {
				for !done {
					tt.args.sess.Handlers.Validate.Run(&request.Request{})
					tt.args.sess.Handlers.Complete.Run(&request.Request{})
					i++
				}
			}()
			time.Sleep(time.Second)
			done = true

			if i > tt.maxReqsASecond {
				t.Logf("number of executions: %d", i)
				t.Fail()
			}
		})
	}
}

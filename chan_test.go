package mml

import (
	"sync"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	s := newScheduler()
	defer s.close()

	c := newChan(s, 0)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		c.send(42)
		wg.Done()
	}()

	go func() {
		if v, ok := c.receive(); !ok || v != 42 {
			t.Error("invalid value received", ok, v)
		}

		wg.Done()
	}()

	wg.Wait()
}

func TestSelect(t *testing.T) {
	t.Run("single send", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c1 := newChan(s, 0)
		c2 := newChan(s, 0)
		c3 := newChan(s, 0)
		done := make(chan struct{})
		go func() {
			b1 := receiveItem(c1)
			b2 := sendItem(c2, 2)
			b3 := receiveItem(c3)
			s := selectItem(s, b1, b2, b3)
			if s != b2 {
				t.Error("wrong item selected")
			}

			close(done)
		}()

		v, ok := c2.receive()
		if !ok || v != 2 {
			t.Error("wrong value received", v)
		}

		<-done
	})

	t.Run("single receive", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c1 := newChan(s, 0)
		c2 := newChan(s, 0)
		c3 := newChan(s, 0)
		done := make(chan struct{})
		go func() {
			b1 := receiveItem(c1)
			b2 := sendItem(c2, 2)
			b3 := receiveItem(c3)
			s := selectItem(s, b1, b2, b3)

			if s != b1 {
				t.Error("wrong item selected")
			}

			if s.value != 1 {
				t.Error("wrong value received", s.value)
			}

			close(done)
		}()

		c1.send(1)
		<-done
	})

	t.Run("multiple send or receive", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c1 := newChan(s, 0)
		c2 := newChan(s, 0)
		c3 := newChan(s, 0)
		done := make(chan struct{})
		go func() {
			var sendDone, receiveDone bool
			for {
				b1 := receiveItem(c1)
				b2 := sendItem(c2, 2)
				b3 := receiveItem(c3)
				s := selectItem(s, b1, b2, b3)

				if s == b2 {
					sendDone = true
				}

				if s == b1 {
					receiveDone = true
					if !s.ok || s.value != 1 {
						t.Error("wrong value received", s.value)
					}
				}

				if sendDone && receiveDone {
					close(done)
					return
				}
			}
		}()

		go func() {
			v, ok := c2.receive()
			if !ok || v != 2 {
				t.Error("wrong value received")
			}
		}()

		go func() {
			c1.send(1)
		}()

		<-done
	})

	t.Run("default item", func(t *testing.T) {
		s := newScheduler()
		defer s.close()
		s.testHook = make(chan struct{}, 1)
		<-s.testHook

		c1 := newChan(s, 0)
		c2 := newChan(s, 0)

		go func() {
			c1.send(1)
		}()

		go func() {
			if v, ok := c2.receive(); !ok || v != 2 {
				t.Error("wrong value received")
			}
		}()

		<-s.testHook
		<-s.testHook

		s1 := receiveItem(c1)
		s2 := sendItem(c2, 2)
		d := defaultItem()

		si := selectItem(s, s1, s2, d)
		if si != s1 && si != s2 {
			t.Error("wrong item received")
		}

		<-s.testHook

		s1 = receiveItem(c1)
		s2 = sendItem(c2, 2)
		d = defaultItem()

		si = selectItem(s, s1, s2, d)
		if si != s1 && si != s2 {
			t.Error("wrong item received")
		}

		<-s.testHook

		s1 = receiveItem(c1)
		s2 = sendItem(c2, 2)
		d = defaultItem()

		si = selectItem(s, s1, s2, d)
		if si != d {
			t.Error("wrong item received")
		}

		<-s.testHook
	})

	t.Run("empty", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		s := newScheduler()
		defer s.close()

		go func() {
			selectItem(s)
			t.Fatal("should not reach here")
		}()

		time.Sleep(9 * time.Millisecond)
	})
}

func TestBuffered(t *testing.T) {
	t.Run("fill", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c := newChan(s, 3)
		c.send(1)
		c.send(2)
		c.send(3)

		if v, ok := c.receive(); !ok || v != 1 {
			t.Fatal("invalid value received", ok, v)
		}

		if v, ok := c.receive(); !ok || v != 2 {
			t.Fatal("invalid value received", ok, v)
		}

		if v, ok := c.receive(); !ok || v != 3 {
			t.Fatal("invalid value received", ok, v)
		}
	})

	t.Run("length", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c := newChan(s, 3)
		if c.len() != 0 {
			t.Error("invalid length")
		}

		c.send(1)
		c.send(2)
		c.send(3)

		if c.len() != 3 {
			t.Error("invalid length")
		}
	})

	t.Run("capacity", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c := newChan(s, 3)
		if c.cap() != 3 {
			t.Error("invalid capacity")
		}

		c.send(1)
		c.send(2)
		c.send(3)

		if c.cap() != 3 {
			t.Error("invalid capacity")
		}
	})
}

func TestClose(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c := newChan(s, 0)
		c.close()

		if _, ok := c.receive(); ok {
			t.Error("failed close the channel")
		}
	})

	t.Run("not empty", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c := newChan(s, 3)
		c.send(1)
		c.send(2)
		c.send(3)
		c.close()

		if v, ok := c.receive(); !ok || v != 1 {
			t.Error("invalid value received", ok, v)
		}

		if v, ok := c.receive(); !ok || v != 2 {
			t.Error("invalid value received", ok, v)
		}

		if v, ok := c.receive(); !ok || v != 3 {
			t.Error("invalid value received", ok, v)
		}

		if _, ok := c.receive(); ok {
			t.Error("failed close the channel")
		}
	})
}

func TestAxioms(t *testing.T) {
	t.Run("send to nil channel", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var c *channel
		go func() {
			c.send(42)
			t.Fatal("should not reach here")
		}()

		time.Sleep(9 * time.Millisecond)
	})

	t.Run("receive from a nil channel", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		var c *channel
		go func() {
			c.receive()
			t.Fatal("should not reach here")
		}()

		time.Sleep(9 * time.Millisecond)
	})

	t.Run("send to a closed channel", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c := newChan(s, 0)
		c.close()

		if err := c.send(42); err != errSendingOnClosed {
			t.Fatal("failed to panic")
		}
	})

	t.Run("receive from a closed channel", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c := newChan(s, 0)
		c.close()

		if _, ok := c.receive(); ok {
			t.Fatal("cannot receive values from a closed channel")
		}
	})

	t.Run("close of a closed channel", func(t *testing.T) {
		s := newScheduler()
		defer s.close()

		c := newChan(s, 0)
		c.close()

		if err := c.close(); err != errCloseClosed {
			t.Fatal("failed to panic")
		}
	})

	t.Run("close of a nil channel", func(t *testing.T) {
		var c *channel
		if err := c.close(); err != errCloseNil {
			t.Fatal("failed to panic")
		}
	})
}

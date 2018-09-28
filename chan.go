package mml

import (
	"errors"
	"math/rand"
)

type communicationType int

const (
	sendRequest communicationType = iota
	receiveRequest
	selectDefault
	lengthRequest
	closeChan
)

type communicationItem struct {
	typ     communicationType
	channel *channel
	value   interface{}
	ok      bool
	err     error
}

type communication struct {
	items       []*communicationItem
	done        chan *communicationItem
	defaultItem *communicationItem
}

type waitingItem struct {
	communication *communication
	item          *communicationItem
}

type scheduler struct {
	communication      chan *communication
	quit               chan struct{}
	sending, receiving map[*channel][]waitingItem
	testHook           chan struct{}
}

type channel struct {
	scheduler *scheduler
	capacity  int
	buffer    []interface{}
	closed    bool
}

var (
	errSendingOnClosed = errors.New("sending on closed channel")
	errCloseClosed     = errors.New("closing closed channel")
	errCloseNil        = errors.New("closing nil channel")
)

func newScheduler() *scheduler {
	s := &scheduler{
		communication: make(chan *communication),
		quit:          make(chan struct{}),
		sending:       make(map[*channel][]waitingItem),
		receiving:     make(map[*channel][]waitingItem),
	}

	go s.run()
	return s
}

func (s *scheduler) cleanup(w waitingItem) {
	for _, ci := range w.communication.items {
		waiting := s.receiving
		if ci.typ == sendRequest {
			waiting = s.sending
		}

		waiting[ci.channel] = waiting[ci.channel][1:]
		if len(waiting[ci.channel]) == 0 {
			delete(waiting, ci.channel)
		}
	}
}

func (s *scheduler) len(c *communication, ci *communicationItem) {
	ci.value = len(ci.channel.buffer)
	c.done <- ci
}

func (s *scheduler) closeChan(c *communication, ci *communicationItem) {
	if ci.channel.closed {
		ci.err = errCloseClosed
		c.done <- ci
		return
	}

	ci.channel.closed = true
	c.done <- ci

	for _, si := range s.sending[ci.channel] {
		si.item.err = errSendingOnClosed
		si.communication.done <- si.item
		s.cleanup(si)
	}

	for _, ri := range s.receiving[ci.channel] {
		ri.communication.done <- ri.item
		s.cleanup(ri)
	}
}

func (s *scheduler) send(c *communication, ci *communicationItem) bool {
	if ci.channel.closed {
		ci.err = errSendingOnClosed
		c.done <- ci
		return true
	}

	r, ok := s.receiving[ci.channel]

	if !ok && ci.channel.capacity == len(ci.channel.buffer) {
		return false
	}

	if !ok {
		ci.channel.buffer = append(
			ci.channel.buffer,
			ci.value,
		)

		c.done <- ci
		return true
	}

	r[0].item.value = ci.value
	r[0].item.ok = true
	c.done <- ci
	r[0].communication.done <- r[0].item

	s.cleanup(r[0])
	return true
}

func (s *scheduler) receive(c *communication, ci *communicationItem) bool {
	if len(ci.channel.buffer) > 0 {
		ci.value = ci.channel.buffer[0]
		ci.channel.buffer = ci.channel.buffer[1:]
		ci.ok = true
		c.done <- ci

		if sending, ok := s.sending[ci.channel]; ok {
			ci.channel.buffer = append(ci.channel.buffer, sending[0].item.value)
			sending[0].communication.done <- sending[0].item
			s.cleanup(sending[0])
		}

		return true
	}

	if ci.channel.closed {
		c.done <- ci
		return true
	}

	sending, ok := s.sending[ci.channel]
	if !ok {
		return false
	}

	ci.value = sending[0].item.value
	ci.ok = true
	sending[0].communication.done <- sending[0].item
	c.done <- ci

	s.cleanup(sending[0])
	return true
}

func (s *scheduler) schedule(c *communication) {
	for _, ci := range c.items {
		waiting := s.receiving
		if ci.typ == sendRequest {
			waiting = s.sending
		}

		waiting[ci.channel] = append(
			waiting[ci.channel],
			waitingItem{
				communication: c,
				item:          ci,
			},
		)
	}
}

func (s *scheduler) run() {
	for {
		if s.testHook != nil {
			s.testHook <- struct{}{}
		}

		var c *communication
		select {
		case c = <-s.communication:
		case <-s.quit:
			return
		}

		var handled bool
		for _, ci := range c.items {
			switch ci.typ {
			case sendRequest:
				handled = s.send(c, ci)
			case receiveRequest:
				handled = s.receive(c, ci)
			case lengthRequest:
				s.len(c, ci)
				handled = true
			case closeChan:
				s.closeChan(c, ci)
				handled = true
			}

			if handled {
				break
			}
		}

		if handled {
			continue
		}

		if c.defaultItem != nil {
			c.done <- c.defaultItem
			continue
		}

		s.schedule(c)
	}
}

func (s *scheduler) close() {
	close(s.quit)
}

func sendItem(c *channel, value interface{}) *communicationItem {
	return &communicationItem{
		typ:     sendRequest,
		channel: c,
		value:   value,
	}
}

func receiveItem(c *channel) *communicationItem {
	return &communicationItem{
		typ:     receiveRequest,
		channel: c,
	}
}

func defaultItem() *communicationItem {
	return &communicationItem{
		typ: selectDefault,
	}
}

func requestLength(c *channel) *communicationItem {
	return &communicationItem{
		typ:     lengthRequest,
		channel: c,
	}
}

func closeItem(c *channel) *communicationItem {
	return &communicationItem{
		typ:     closeChan,
		channel: c,
	}
}

func newChan(s *scheduler, capacity int) *channel {
	return &channel{
		scheduler: s,
		capacity:  capacity,
	}
}

func (c *channel) send(v interface{}) {
	if c == nil {
		select {}
	}

	selectItem(c.scheduler, sendItem(c, v))
}

func (c *channel) receive() (interface{}, bool) {
	if c == nil {
		select {}
	}

	response := selectItem(c.scheduler, receiveItem(c))
	return response.value, response.ok
}

func (c *channel) len() int {
	response := selectItem(c.scheduler, requestLength(c))
	return response.value.(int)
}

func (c *channel) cap() int { return c.capacity }

func (c *channel) close() {
	if c == nil {
		panic(errCloseNil)
	}

	selectItem(c.scheduler, closeItem(c))
}

func selectItem(s *scheduler, items ...*communicationItem) *communicationItem {
	citems := make([]*communicationItem, len(items))
	copy(citems, items)
	items = citems

	var defaultItem *communicationItem
	for i, item := range items {
		if item.typ == selectDefault {
			defaultItem = item
			items, items[len(items)-1] = append(items[:i], items[:i+1]...), nil
			break
		}
	}

	rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })

	comm := &communication{
		done:        make(chan *communicationItem),
		items:       items,
		defaultItem: defaultItem,
	}

	s.communication <- comm
	response := <-comm.done
	if response.err != nil {
		panic(response.err)
	}

	return response
}

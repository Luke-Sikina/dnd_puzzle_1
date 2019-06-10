package main

import (
	"log"
)

const ChannelSize = 10000

type Predicate func(subject interface{}) bool
type Stream chan interface{}
type Filterable struct {
	Streams []Stream
	Errored bool
}

func NewFilterable(unfiltered Stream) *Filterable {
	return &Filterable{
		[]Stream{unfiltered},
		false,
	}
}

func (objects *Filterable) Filter(predicate Predicate) *Filterable {
	if objects.Errored {
		log.Printf("Skipping because of previous error")
		return objects
	}
	streamNum := len(objects.Streams)
	current := objects.Streams[len(objects.Streams)-1]
	next := make(Stream, ChannelSize)
	objects.Streams = append(objects.Streams, next)
	// append isnt dont asynchronously, so a simple slice will do
	// However, it is important to hang on to current and next
	// rather than getting them in the goroutine itself

	go func() {
		log.Printf("Begining filtering from stream %d to stream %d", streamNum-1, streamNum)
		for object := range current {
			if predicate(object) {
				next <- object
			}
		}
		close(next) // this looks weird, but it's right. This is where the objects for next
		// are being generated, so this code is responsible for closing that channel
		log.Printf("Filtering complete and lock released")
	}()

	return objects
}

func (objects *Filterable) ToSlice() (filteredObjects []interface{}) {
	for object := range objects.Streams[len(objects.Streams)-1] {
		filteredObjects = append(filteredObjects, object)
	}
	return
}

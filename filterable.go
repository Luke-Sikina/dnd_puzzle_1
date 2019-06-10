package main

import (
	"context"
	"golang.org/x/sync/semaphore"
	"log"
)

const ChannelSize = 1000000

type Predicate func(subject interface{}) bool
type Filterable struct {
	Current chan interface{}
	Next    chan interface{}
	Ready   *semaphore.Weighted
	Context context.Context
	Errored bool
}

func NewFilterable(unfiltered chan interface{}) *Filterable {
	sem := semaphore.NewWeighted(1)
	return &Filterable{
		unfiltered,
		make(chan interface{}, ChannelSize),
		sem,
		context.TODO(),
		false,
	}
}

func (objects *Filterable) Filter(predicate Predicate) *Filterable {
	if objects.Errored {
		log.Printf("Skipping because of previous error")
		return objects
	}

	err := objects.Ready.Acquire(objects.Context, 1) // wait for last round of filtering to finish
	if err != nil {
		log.Printf("Error aquiring lock to filter: %v", err)
		return objects
	}

	log.Printf("Lock aquired")

	go func() {
		log.Printf("Begining filtering")
		for object := range objects.Current {
			if predicate(object) {
				objects.Next <- object
			}
		}
		close(objects.Next)
		objects.Current = objects.Next
		objects.Next = make(chan interface{}, ChannelSize)
		objects.Ready.Release(1)

		log.Printf("Filtering complete and lock released")
	}()

	return objects
}

func (objects *Filterable) ToSlice() (filteredObjects []interface{}) {
	err := objects.Ready.Acquire(objects.Context, 1) // wait for last round of filtering to finish
	if err != nil {
		log.Printf("Error aquiring lock to filter: %v", err)
		return
	}

	for object := range objects.Current {
		filteredObjects = append(filteredObjects, object)
	}
	return
}

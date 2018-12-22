// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package models

import (
	"container/list"
	"time"
)

type EventType int

const (
	EVENT_CONN = iota
	EVENT_DISCONN
	EVENT_JOIN
	EVENT_LEAVE
	EVENT_MESSAGE
	EVENT_LOGIN
	EVENT_LOGOUT
)

type Event struct {
	Type      EventType                   `json:"type"` // JOIN, LEAVE, MESSAGE
	Data      map[interface{}]interface{} `json:"data"`
	Timestamp int                         `time:"time"` // Unix timestamp (secs)
}

const archiveSize = 20

func NewEvent(ep EventType, data map[interface{}]interface{}) Event {
	return Event{Type: ep, Data: data, Timestamp: int(time.Now().Unix())}
}

// Event archives.
var archive = list.New()

// NewArchive saves new event to archive list.
func NewArchive(event Event) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

// GetEvents returns all events after lastReceived.
func GetEvents(lastReceived int) []Event {
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}

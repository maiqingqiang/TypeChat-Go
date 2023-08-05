// The following types define the structure of an object of type CalendarActions that represents a list of requested calendar actions
package main

type CalendarActions struct {
	Actions []Action
}

type Action struct {
	AddEventAction          *AddEventAction          `json:"add_event,omitempty"`
	RemoveEventAction       *RemoveEventAction       `json:"remove_event,omitempty"`
	AddParticipantsAction   *AddParticipantsAction   `json:"add_participants,omitempty"`
	ChangeTimeRangeAction   *ChangeTimeRangeAction   `json:"change_time_range,omitempty"`
	ChangeDescriptionAction *ChangeDescriptionAction `json:"change_description,omitempty"`
	FindEventsAction        *FindEventsAction        `json:"find_events,omitempty"`
	UnknownAction           *UnknownAction           `json:"unknown,omitempty"`
}

type AddEventAction struct {
	Event *Event
}

type RemoveEventAction struct {
	EventReference *EventReference
}

type AddParticipantsAction struct {
	// event to be augmented; if not specified assume last event discussed
	EventReference *EventReference `json:"event_reference,omitempty"`
	// new participants (one or more)
	Participants []string
}

type ChangeTimeRangeAction struct {
	// event to be changed
	EventReference *EventReference
	// new time range for the event
	TimeRange *EventTimeRange
}

type ChangeDescriptionAction struct {
	// event to be changed
	EventReference *EventReference `json:"event_reference,omitempty"`
	// new description for the event
	Description string
}

type FindEventsAction struct {
	// one or more event properties to use to search for matching events
	EventReference *EventReference
}

// UnknownAction if the user types text that can not easily be understood as a calendar action, this action is used
type UnknownAction struct {
	// text typed by the user that the system did not understand
	Text string
}

type EventTimeRange struct {
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Duration  string `json:"duration,omitempty"`
}

type Event struct {
	// date (example: March 22, 2024) or relative date (example: after EventReference)
	Day         string
	TimeRange   *EventTimeRange
	Description string
	Location    string `json:"location,omitempty"`
	// a list of people or named groups like 'team'
	Participants []string `json:"participants,omitempty"`
}

// EventReference properties used by the requester in referring to an event
// these properties are only specified if given directly by the requester
type EventReference struct {
	// date (example: March 22, 2024) or relative date (example: after EventReference)
	Day string `json:"day,omitempty"`
	// (examples: this month, this week, in the next two days)
	DayRange     string          `json:"day_range,omitempty"`
	TimeRange    *EventTimeRange `json:"time_range,omitempty"`
	Description  string          `json:"description,omitempty"`
	Location     string          `json:"location,omitempty"`
	Participants []string        `json:"participants,omitempty"`
}

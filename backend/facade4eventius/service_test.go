package facade4eventius

import "testing"

func TestEventScheduleState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		event Event
		want  ScheduleState
	}{
		{name: "title only", event: Event{Title: "Picnic"}, want: ScheduleStatePlanning},
		{name: "date only", event: Event{Title: "Picnic", Date: "2026-08-12"}, want: ScheduleStatePlanning},
		{name: "time only", event: Event{Title: "Picnic", Time: "18:30"}, want: ScheduleStatePlanning},
		{name: "location only", event: Event{Title: "Picnic", Location: "The park"}, want: ScheduleStatePlanning},
		{name: "scheduled", event: Event{Title: "Picnic", Date: "2026-08-12", Time: "18:30"}, want: ScheduleStateScheduled},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := test.event.ScheduleState(); got != test.want {
				t.Fatalf("ScheduleState() = %q, want %q", got, test.want)
			}
		})
	}
}

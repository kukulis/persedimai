package di

import (
	"darbelis.eu/persedimai/internal/aviation_edge"
	"fmt"
	"reflect"
	"testing"
)

func Custom() aviation_edge.ScheduleConsumer {
	return &aviation_edge.PrintScheduleConsumer{}
}

func TestWrap(t *testing.T) {
	Wrap(Custom)
	rez := Wrap(Custom)

	got := reflect.TypeOf(rez)
	want := reflect.TypeOf(&aviation_edge.PrintScheduleConsumer{})

	fmt.Printf("Types:  %v and %v\n", got, want)

	if got != want {
		t.Errorf("The type of rez %v does not match the expected type %v", got, want)
	}
}

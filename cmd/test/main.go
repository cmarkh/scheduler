package main

import (
	"fmt"
	"time"

	"github.com/cmarkh/scheduler"
)

func main() {
	testSort()

	schedule()

}

func testSort() {
	now := time.Now()

	times := []scheduler.RunTime{}
	times = append(times, scheduler.RunTime{Weekdays: []time.Weekday{now.Weekday()}, Hour: 1, Minute: 1, Second: 1})
	times = append(times, scheduler.RunTime{Weekdays: []time.Weekday{now.Weekday()}, Hour: 22, Minute: 6, Second: 1})
	times = append(times, scheduler.RunTime{Weekdays: []time.Weekday{now.Weekday()}, Hour: 22, Minute: 4, Second: 1})
	times = append(times, scheduler.RunTime{Weekdays: []time.Weekday{now.Weekday()}, Hour: 2, Minute: 1, Second: 5})

	timesParsed := scheduler.ParseTimes(times, false)

	fmt.Println(timesParsed)
}

func schedule() {

	schedule := scheduler.New()
	//schedule.SkipMissed = true

	schedule.Add(a)
	schedule.Add(b)
	//schedule.Add(c)
	//schedule.Add(d)

	//schedule.Run()
	schedule.RunHandleSleeps(0)

}

var a = scheduler.Func{
	Name: "a",
	Fn: func() error {
		fmt.Println("a")
		return nil
	},
	Times: []scheduler.RunTime{
		{Weekdays: []time.Weekday{time.Now().Weekday()}, Hour: time.Now().Hour(), Minute: time.Now().Minute() + 1, Second: 0},
		{Weekdays: []time.Weekday{time.Now().Weekday()}, Hour: time.Now().Hour(), Minute: time.Now().Minute() + 1, Second: 10},
	},
}

var b = scheduler.Func{
	Name: "b",
	Fn: func() error {
		fmt.Println("b")
		return nil
	},
	Times: []scheduler.RunTime{
		{Weekdays: []time.Weekday{time.Now().Weekday()}, Hour: time.Now().Hour(), Minute: time.Now().Minute() + 1, Second: 20},
		{Weekdays: []time.Weekday{time.Now().Weekday()}, Hour: time.Now().Hour(), Minute: time.Now().Minute() + 1, Second: 30},
	},
}

var c = scheduler.Func{
	Name: "c",
	Fn: func() error {
		fmt.Println("c")
		return nil
	},
	Times: []scheduler.RunTime{{Weekdays: []time.Weekday{4}, Hour: 1, Minute: 1, Second: 1}},
}

var d = scheduler.Func{
	Name: "d",
	Fn: func() error {
		fmt.Println("d")
		return nil
	},
	Times: []scheduler.RunTime{{Weekdays: []time.Weekday{4}, Hour: 1, Minute: 1, Second: 1}},
}

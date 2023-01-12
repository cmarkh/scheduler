package scheduler

import (
	"flag"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cmarkh/errs"
	"github.com/cmarkh/go-mail"
	"golang.org/x/exp/slices"
)

var (
	Weekdays = []time.Weekday{1, 2, 3, 4, 5}
	Everyday = []time.Weekday{0, 1, 2, 3, 4, 5, 6}
)

type Schedule struct {
	SkipMissed bool //skip funcs that should have already been run (time is in past)
	ErrEmails  []string
	functions  []Func
	wg         sync.WaitGroup
}

type Func struct {
	Name        string       //name of the program being run
	Fn          func() error //func to run
	Times       []RunTime    //times to run the func
	timesParsed []time.Time  //Times will get converted to this for use
}

type RunTime struct {
	Weekdays []time.Weekday
	Hour     int
	Minute   int
	Second   int
}

func New(emailOnErr ...string) Schedule {
	skipMissed := flag.Bool("skip", false, "skip functions scheduled for prior to start time")
	flag.Parse()

	return Schedule{
		ErrEmails:  emailOnErr,
		SkipMissed: *skipMissed,
	}
}

func (schedule *Schedule) Add(run Func) {
	run.timesParsed = ParseTimes(run.Times, schedule.SkipMissed)
	if len(run.timesParsed) == 0 {
		return
	}
	schedule.functions = append(schedule.functions, run)
}

// RunHandleSleeps schedules the functions and accounts for any computer sleeps on an interval basis (ie function won't be later than the interval)
// Note: If interval == 0, will use default interval (10 minutes)
func (schedule *Schedule) RunHandleSleeps(interval time.Duration) {
	schedule.startMsg()

	if interval == 0 {
		interval = time.Minute * 10
	}

	tick := time.NewTicker(time.Second) //start with just 1 second so first iteration fires immediately

	for range tick.C {
		tick.Reset(interval) //reset ticker to the correct interval after the first run

		now := time.Now()
		newfunctions := []Func{}
		for _, fn := range schedule.functions {
			newtimes := []time.Time{}
			for i, runtime := range fn.timesParsed {
				diff := runtime.Sub(now)
				if i != 0 && diff < 0 {
					continue //skip this runtime
				}
				if diff <= interval {
					f, t := fn, runtime
					schedule.wg.Add(1)
					time.AfterFunc(diff, func() { schedule.execute(f, t) })
				} else {
					newtimes = append(newtimes, runtime) //keep the future runtimes only
				}
			}
			fn.timesParsed = newtimes
			if len(fn.timesParsed) > 0 {
				newfunctions = append(newfunctions, fn)
			}
		}
		schedule.functions = newfunctions
		if len(schedule.functions) == 0 {
			tick.Stop()
			break
		}
	}

	schedule.wg.Wait()
	fmt.Println("\ndone")
}

// Run just schedules the functions to run. Any sleeps will cause delays to when the functions run. Use if the system will always be on
func (schedule *Schedule) Run() {
	schedule.startMsg()

	for _, fn := range schedule.functions {
		for _, runtime := range fn.timesParsed {
			f, t := fn, runtime
			schedule.wg.Add(1)
			time.AfterFunc(time.Until(t), func() { schedule.execute(f, t) })
		}
	}

	schedule.wg.Wait()
	fmt.Println("\ndone")
}

func (schedule *Schedule) startMsg() {
	fmt.Printf("Starting Scheduler - %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("Skip functions scheduled for prior to start time? - " + fmt.Sprint(schedule.SkipMissed))
	fmt.Println()
}

func (schedule *Schedule) execute(fn Func, t time.Time) {
	defer schedule.wg.Done()
	defer errs.CatchAndRecover()

	fmt.Printf("%s - Running %s - for %v\n", time.Now().Format("15:04:05"), fn.Name, t)

	err := fn.Fn()
	if err != nil {
		errs.Log(err)
		mail.Send(fn.Name+" Error", fmt.Sprintln(err), schedule.ErrEmails...)
	}
}

// ParseTimes returns just the times scheduled for today and sorts them
func ParseTimes(times []RunTime, skipMissed bool) (future []time.Time) {
	//remove times not scheduled for today's day of the week
	todays := []RunTime{}
	for _, t := range times {
		if slices.Contains(t.Weekdays, time.Now().Weekday()) {
			todays = append(todays, t)
		}
	}
	if len(todays) == 0 {
		return
	}

	//remove times that should have run earlier (leave one if skipMissed is false so runs at least once)
	now := time.Now()
	for i, t := range todays {
		if (i == 0 && !skipMissed) || !past(t, now) {
			future = append(future, time.Date(now.Year(), now.Month(), now.Day(), t.Hour, t.Minute, t.Second, 0, time.Local))
		}
	}

	sort.Slice(future, func(i, j int) bool {
		return future[i].Before(future[j])
	})

	return
}

func past(t RunTime, now time.Time) bool {
	if now.Hour() > t.Hour {
		return true
	}
	if now.Hour() == t.Hour && now.Minute() > t.Minute {
		return true
	}
	if now.Hour() == t.Hour && now.Minute() == t.Minute && now.Second() > t.Second {
		return true
	}

	return false
}

package main

import (
	"fmt"
	"time"
	"strings"
	"github.com/fatih/color"
	"os"
)

type ProgressPrinter struct {
	done chan error
	wait chan struct{}
	quite bool
}

func NewProgressPrinter(quite bool) *ProgressPrinter {
	return &ProgressPrinter{make(chan error), make(chan struct{}), quite}
}

func (pp *ProgressPrinter) print(a ...interface{}) {
	if pp.quite {
		return
	}
	fmt.Print(a...)
}

func (pp *ProgressPrinter) green(format string, a ...interface{}) {
	if pp.quite {
		return
	}
	color.Green(format, a...)
}

func (pp *ProgressPrinter) failOnError(err error) {
	if err == nil {
		return
	}
	color.Red(err.Error())
	os.Exit(1)
}

func (pp *ProgressPrinter) Start(message string) {
	go func(message string) {
		pp.print(message)
		cnt := 1
		for {
			select {
			case <-time.After(time.Millisecond * 500):
				pp.print("\r", message, strings.Repeat(".", cnt))
				cnt = cnt % 3
				cnt++
			case err := <-pp.done:
				pp.print("\r")
				if err != nil {
					pp.failOnError(fmt.Errorf("%s failed: %s", message, err.Error()))
				} else {
					pp.green("%s%s", message, strings.Repeat(".", 3))
				}
				pp.wait <- struct{}{}
				return
			}
		}
		pp.wait <- struct{}{}
		return
	}(message)
}

func (pp *ProgressPrinter) Complete() {
	pp.done <- nil
	<- pp.wait
}

func (pp *ProgressPrinter) CheckAndComplete(err error) {
	pp.done <- err
	<- pp.wait
}

func (pp *ProgressPrinter) CheckError(err error) {
	if err != nil {
		pp.done <- err
		<- pp.wait
	}
}
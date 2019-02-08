# go-dispatcher

[![Build Status](https://travis-ci.org/danielpoonwj/go-dispatcher.svg?branch=master)](https://travis-ci.org/danielpoonwj/go-dispatcher)

## Introduction

This library is a simple implementation of a worker pool in `go`, building on the concepts from [this](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/) article.

## Features

### Variable concurrency

Instead of spawning an unbounded number of goroutines, the number of concurrent jobs to be processed can be set. This is useful for managing system resources, as well as cases where the job processing involves interacting with another rate sensitive service.

### Variable buffer

If the rate at which jobs are processed is slower than the rate at which they are added, those waiting are added into a queue - the size of which can be set as well.

### Graceful shutdown

Calling `Stop()` waits for all added jobs to finish.

## Installing

```
go get github.com/danielpoonwj/go-dispatcher
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/danielpoonwj/go-dispatcher"
)

type Job struct {
	msg string
}

func (j *Job) Process() {
	fmt.Println(j.msg)
}

func NewJob(msg string) *Job {
	return &Job{msg: msg}
}

func main() {
	d := dispatcher.NewDispatcher(1, 3)
	d.Start()

	j1 := NewJob("Job 1")
	d.AddJob(j1)

	j2 := NewJob("Job 2")
	d.AddJob(j2)

	d.Stop()
}
```
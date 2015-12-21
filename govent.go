package main

//Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
//Use of this source code is governed by a BSD-style
//license that can be found in the LICENSE file.

import (
	"fmt"
	"log"
	"os"

	"github.com/MediaMath/govent/cligraphite"
	"github.com/MediaMath/govent/graphite"
	"github.com/codegangsta/cli"
)

var (
	whatFlag = cli.StringFlag{
		Name:   "what",
		Usage:  "The 'What' field in the event.",
		EnvVar: "GOVENT_WHAT",
	}

	tagsFlag = cli.StringSliceFlag{
		Name:   "tag",
		Usage:  "The 'Tag' field in the event.",
		EnvVar: "GOVENT_TAGS",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "govent"
	app.Usage = "send events to the graphite api"
	app.Flags = append(cligraphite.Flags, whatFlag, tagsFlag)

	app.Action = func(ctx *cli.Context) {
		event, err := eventFromCtx(ctx)
		if err != nil {
			log.Fatal(err)
		}

		client, err := cligraphite.NewClientFromContext(ctx)
		if err != nil {
			log.Fatal(err)
		}

		err = client.Publish(event)
		if err != nil {
			log.Fatal(err)
		}
	}

	app.Run(os.Args)
}

func eventFromCtx(ctx *cli.Context) (*graphite.Event, error) {
	what := ctx.String(whatFlag.Name)
	if what == "" {
		return nil, fmt.Errorf("%s is required", whatFlag.Name)
	}

	if len(ctx.Args()) != 1 {
		return nil, fmt.Errorf("Must provide data to post")
	}
	data := ctx.Args()[0]

	tags := ctx.StringSlice(tagsFlag.Name)

	return graphite.NewEvent(what, data, tags...), nil
}
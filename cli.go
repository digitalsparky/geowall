package main

import (
  "github.com/urfave/cli/v2"
)

type CLI struct {
  App *cli.App
}

func (*c CLI) Start() {
  c.App := &cli.App{
    Commands: []*cli.Command{
      {
          Name:    "add",
          Aliases: []string{"a"},
          Usage:   "add a task to the list",
          Action: func(cCtx *cli.Context) error {
              fmt.Println("added task: ", cCtx.Args().First())
              return nil
          },
      },
      {
          Name:    "complete",
          Aliases: []string{"c"},
          Usage:   "complete a task on the list",
          Action: func(cCtx *cli.Context) error {
              fmt.Println("completed task: ", cCtx.Args().First())
              return nil
          },
      },
    },
  }
}
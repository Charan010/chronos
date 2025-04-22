package main

import (
	"log"
	"os"

	"github.com/Charan010/chronos/internal"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert/yaml"
)

func SpinUpJobs(cfg *internal.Config, c *cron.Cron) {

	for _, job := range cfg.Jobs {
		if !job.Enabled {
			continue
		}

		j := job
		log.Printf(color.New(color.FgGreen).Sprint("ADDED: ")+"CRON job Name: %v | Schedule: %v\n", j.Name, j.Schedule)

		c.AddFunc(job.Schedule, func() {
			log.Printf(color.New(color.FgCyan).Sprint("EXEC: ")+"Executing CRON Job Name: %v | Schedule: %v\n", j.Name, j.Schedule)
			internal.DumpBackUp(j)
		})
	}
}

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg internal.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	c := cron.New()

	SpinUpJobs(&cfg, c)

	c.Start()

	select {}
}

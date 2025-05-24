package main

import (
	"log"
	"os"
	"flag"
	"context"
	"os/signal"
	"syscall"
	"sync"
	"github.com/Charan010/chronos/internal"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert/yaml"
)

func SpinUpJobs(cfg *internal.Config, c *cron.Cron, wp *internal.WorkerPool){

	for _,job := range cfg.Jobs{
		if !job.Enabled{
			continue
		}

		j := job
		log.Printf(color.New(color.FgGreen).Sprint("ADDED: ")+"CRON job Name: %v | Schedule: %v\n", j.Name, j.Schedule)

		
		c.AddFunc(j.Schedule ,func(){
			log.Printf(color.New(color.FgCyan).Sprint("EXEC: ")+"Dispatching job to worker pool: %v\n", j.Name)
			
			wp.JobQueue <- internal.Job{
				Name: j.Name,
					Execute: func(){
					internal.DumpBackup(j)
				},

			}
		})
	}
}


func main(){

	workers := flag.Int("workers",4, "Number of concurrent workers")
	queueSize := flag.Int("queue",20,"job queue size")

	flag.Parse()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg internal.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wp := internal.NewWorkerPool(*workers, *queueSize, &wg)

	c := cron.New()

	SpinUpJobs(&cfg, c, wp)

	c.Start()

	ctx, stop := signal.NotifyContext(context.Background(),os.Interrupt,syscall.SIGTERM)

	defer stop()

	<- ctx.Done()

	log.Println("Shutdown signal recieved")

	c.Stop()

	close(wp.JobQueue)

	wg.Wait()


	log.Println("All workers task completed succesfully,gracefully shutting down")


}

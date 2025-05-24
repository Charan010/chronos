## Chronos - Mysql DB Backup Scheduler

## Setup:
   . Put your jobs in config.yaml

   . Each jobs needs its own name,enabled flag,schedule(cron type),db_type, username,password,database and compressed option.

## Usage:
    . Run main.go from chronos directory.
    
    . thrown in flags like -workers and -queue to set concurrency and queueSize to execute multiple jobs simultaneously using goroutines.

## Features (or so i thought):

    . Parallel worker pool
    . Graceful shutdown
    . Supports gzip compression

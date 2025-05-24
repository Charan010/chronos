package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"sync"
	"github.com/fatih/color"
	"compress/gzip"
	"io"
)

type Config struct {
	Jobs []Job `yaml:"jobs"`
}

type Job struct {
	Name       string `yaml:"name"`
	Enabled    bool   `yaml:"enabled"`
	Status     string `yaml:"status"`
	Schedule   string `yaml:"schedule"`
	DB_type    string `yaml:"db_type"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
	Compressed bool   `yaml:"compressed"`
	Execute    func()   `yaml:"-"`
}




/* Workerpool picks up tasks from channel (think of it like a queue with limited 
 size) where go routine is spun up by go compiler to execute each task.
 
 */

type WorkerPool struct{
     JobQueue chan Job	

}

func NewWorkerPool(count int ,queueSize int, wg *sync.WaitGroup)*WorkerPool{

	wp := &WorkerPool{
		JobQueue: make(chan Job, queueSize),
	}

	for i := 0 ; i < count ; i++{
		wg.Add(1)	
		go wp.worker(i, wg)
	}
	return wp 

}

func (wp *WorkerPool) worker(id int, wg *sync.WaitGroup){
	defer wg.Done()
		
	for job := range wp.JobQueue {
		log.Printf("[Worker %d] Starting job: %s\n", id, job.Name)
		if job.Execute != nil {
			job.Execute()
		}
		log.Printf("[Worker %d] Finished job: %s\n", id, job.Name)
	}
}

func generateBackupName() string {

	currentTime := time.Now()
	backupName := currentTime.Format("2006-01-02 15:04:05")

	backupName = strings.ReplaceAll(backupName, " ", "_")
	backupName = strings.ReplaceAll(backupName, ":", "_")

	backupName += "_backup"

	return backupName
}


/* Creates new folder for each cron job (if it doesnt exist) and can store 
 backup files in compressed format using golang package.
 
 */
func DumpBackup(job Job) {

	backupName := generateBackupName()

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	backupDir := filepath.Join(cwd, "backups")

	if job.DB_type == "mysql" {
		backupDir = filepath.Join(backupDir, "mysql")
	}

	backupDir = filepath.Join(backupDir, job.Name)

	err = os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		log.Printf(color.New(color.BgRed).Sprint("ERROR: ")+"Error creating directory: %v", err)
		return
	}

	sqlFilePath := filepath.Join(backupDir, backupName+".sql")

	var outFile *os.File
	var outWriter io.Writer

	if job.Compressed {
		compressedFilePath := filepath.Join(backupDir, backupName+".sql.gz")

		outFile, err = os.Create(compressedFilePath)
		if err != nil {
			log.Printf(color.New(color.BgRed).Sprint("ERROR: ")+"Error creating compressed file: %v", err)
			return
		}
		defer outFile.Close()

		gzipWriter := gzip.NewWriter(outFile)
		defer gzipWriter.Close()

		outWriter = gzipWriter

	} else {
		outFile, err = os.Create(sqlFilePath)
		if err != nil {
			log.Printf(color.New(color.BgRed).Sprint("ERROR: ")+"Error creating file: %v", err)
			return
		}
		defer outFile.Close()

		outWriter = outFile
	}

	cmd := exec.Command("mysqldump", "-u", job.Username, job.Database)
	cmd.Env = append(os.Environ(), "MYSQL_PWD="+job.Password)

	cmd.Stdout = outWriter
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Printf(color.New(color.BgRed).Sprint("ERROR: ")+"Error dumping database: %v", err)
		return
	}

	if job.Compressed {
		log.Println("Database dump successfully saved to:", filepath.Join(backupDir, backupName+".sql.gz"))
	} else {
		log.Println("Database dump successfully saved to:", sqlFilePath)
	}
}
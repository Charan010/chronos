package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Jobs []Job `yaml:"jobs"`
}

type Job struct {
	Name     string `yaml:"name"`
	Schedule string `yaml:"schedule"`
	Status   string `yaml:"status"`
	Enabled  bool   `yaml:"enabled"`
	Db_type  string `yaml:"db_type"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func MarkTheMoment() string {

	currentTime := time.Now()
	backupName := currentTime.Format("2006-01-02 15:04:05")

	backupName = strings.ReplaceAll(backupName, " ", "_")
	backupName = strings.ReplaceAll(backupName, ":", "_")

	backupName += "_backup"

	return backupName
}

func DumpBackUp(job Job) {

	backupName := MarkTheMoment()

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	backupDir := filepath.Join(cwd, "backups")

	if job.Db_type == "mysql" {
		backupDir = filepath.Join(backupDir, "mysql")
	}

	outFilePath := filepath.Join(backupDir, backupName+".sql")

	outFile, err := os.Create(outFilePath)

	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	defer outFile.Close()

	cmd := exec.Command("mysqldump", "-u", job.Username, "-p"+job.Password, job.Database)

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	log.Println("Database dump succesfully saved to Backup Folder :)")

}

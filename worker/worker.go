package worker

import (
	"fmt"
	"jobScheduler/logger"
	"jobScheduler/models"
	"jobScheduler/scheduler"
	"os/exec"
	"strings"
	"time"

	"gorm.io/gorm"
)

var JobQueue chan models.Job

func StartWorkerPool(poolSize int, queueSize int, db *gorm.DB) {
	JobQueue = make(chan models.Job, queueSize)
	logger.L.Info("Job queue initialized", "size", queueSize)

	for i := 1; i <= poolSize; i++ {
		go worker(i, db)
	}
	logger.L.Info("Worker pool started", "workers", poolSize)

	go schedulerTicker(db)
	logger.L.Info("Scheduler started")
}

func worker(id int, db *gorm.DB) {
	for job := range JobQueue {
		logger.L.Info("Worker picked up a job", "worker_id", id, "job_id", job.ID)

		// Update main job status to "running"
		db.Model(&job).Updates(map[string]interface{}{"status": "running", "last_run_at": time.Now()})

		output, err := executeCommand(job.Command)

		executionStatus := "succeeded"
		if err != nil {
			executionStatus = "failed"
			logger.L.Error("Job execution failed", "job_id", job.ID, "error", err, "output", output)
		} else {
			logger.L.Info("Job execution succeeded", "job_id", job.ID, "output", output)
		}

		db.Model(&job).Update("status", executionStatus)

		// Create the detailed execution record
		executionRecord := models.JobExecution{
			JobID:      job.ID,
			Status:     executionStatus,
			Output:     output,
			FinishedAt: time.Now(),
		}
		if result := db.Create(&executionRecord); result.Error != nil {
			logger.L.Error("Failed to save job execution history", "job_id", job.ID, "error", result.Error)
		}
	}
}

func executeCommand(command string) (string, error) {
	if strings.HasPrefix(command, "http") {
		return fmt.Sprintf("Simulated HTTP GET to %s", command), nil
	}
	// #nosec G204
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func schedulerTicker(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {

		var pendingJobs []models.Job
		db.Where("status = ?", "pending").Find(&pendingJobs)

		if len(pendingJobs) == 0 {
			continue
		}
		logger.L.Info("Found pending jobs", "count", len(pendingJobs))

		for _, job := range pendingJobs {
			if scheduler.IsDue(job, t) {
				select {
				case JobQueue <- job:
					logger.L.Info("Job queued for execution", "job_id", job.ID)
				default:
					logger.L.Warn("Job queue is full. Cannot queue job.", "job_id", job.ID)
				}
			}
		}
	}
}

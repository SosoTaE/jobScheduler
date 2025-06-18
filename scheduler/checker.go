package scheduler

import (
	"jobScheduler/models"
	"time"
)

func contains[T comparable](slice []T, item T) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func IsDue(job models.Job, t time.Time) bool {
	schedule := job.Schedule // Get the schedule object from the job

	if len(schedule.Years) > 0 && !contains(schedule.Years, t.Year()) {
		return false
	}

	// --- The rest of the checks remain the same ---

	// Check Month
	if len(schedule.Months) > 0 && !contains(schedule.Months, int(t.Month())) {
		return false
	}

	// Check Day of Month
	if len(schedule.DaysOfMonth) > 0 && !contains(schedule.DaysOfMonth, t.Day()) {
		return false
	}

	// Check Day of Week
	if len(schedule.Weekdays) > 0 && !contains(schedule.Weekdays, t.Weekday()) {
		return false
	}

	// If we get here, the date part matches. Now check if any of the times match.
	for _, runTime := range schedule.Times {
		if runTime.Hour == t.Hour() && runTime.Minute == t.Minute() {
			return true // Found a matching time! The job is due.
		}
	}

	// No matching time was found for today.
	return false
}

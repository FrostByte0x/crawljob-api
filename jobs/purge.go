package jobs

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func purgeOldFiles(destinationFolder string, ageInHours int) {
	slog.Info("Purge job activating.")
	files, err := os.ReadDir(destinationFolder)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	for _, item := range files {
		fileInformation, err := item.Info()
		if err != nil {
			slog.Warn(err.Error())
			continue // go to the next file
		}
		lastModificationTime := fileInformation.ModTime()
		if time.Since(lastModificationTime) > time.Hour*time.Duration(ageInHours) {
			stringTime := lastModificationTime.Format("02-01-2006 15:04")
			slog.Info(fmt.Sprintf("Removing %s, last modified %s", fileInformation.Name(), stringTime))
			// Delete the folder and its content
			err := os.RemoveAll(filepath.Join(destinationFolder, fileInformation.Name()))
			if err != nil {
				slog.Warn(err.Error())
			}
		}
	}
	slog.Info("Purge job has completed successfully.")
}

func StartPurgeRoutine(destinationFolder string, ageInHours int) {
	go func() {
		for {
			slog.Info("Purge will run in one hour. Exit the program now to avoid file deletion.")
			time.Sleep(1 * time.Hour)
			purgeOldFiles(destinationFolder, ageInHours)
		}
	}()
}

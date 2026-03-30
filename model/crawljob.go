package model

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CrawlJob struct {
	URL                        string // text=
	Enabled                    bool   // true
	Comment                    string // Comment
	AutoStart                  bool   // true
	ExtractafterDownload       bool   // false
	ForcedStart                bool   // false
	DownloadFolder             string // /mnt/jDownloader/?
	OverwritePackagizerEnabled bool   // false
	AutoConfirm                bool   // true
}

func GenerateJobFile(url, destinationFolder, fileDestination string) error {
	jobFile := CrawlJob{
		URL:                        url,
		Enabled:                    true,
		Comment:                    "Created by crawljob-api",
		AutoConfirm:                true,
		AutoStart:                  true,
		ExtractafterDownload:       false,
		ForcedStart:                false,
		OverwritePackagizerEnabled: false,
		DownloadFolder:             destinationFolder,
	}
	// Create the file
	lines := []string{}
	lines = append(lines, fmt.Sprintf("enabled=%s", booltoString(jobFile.Enabled)))
	lines = append(lines, fmt.Sprintf("text=%s", jobFile.URL))
	lines = append(lines, fmt.Sprintf("comment=%s", jobFile.Comment))
	lines = append(lines, fmt.Sprintf("autoConfirm=%s", booltoString(jobFile.AutoConfirm)))
	lines = append(lines, fmt.Sprintf("autoStart=%s", booltoString(jobFile.AutoStart)))
	lines = append(lines, fmt.Sprintf("extractAfterDownload=%s", booltoString(jobFile.ExtractafterDownload)))
	lines = append(lines, fmt.Sprintf("forcedStart=%s", booltoString(jobFile.ForcedStart)))
	lines = append(lines, fmt.Sprintf("downloadFolder=%s", jobFile.DownloadFolder))
	lines = append(lines, fmt.Sprintf("overwritePackagizerEnabled=%s", booltoString(jobFile.OverwritePackagizerEnabled)))

	fileContent := strings.Join(lines, "\n")
	nameHash := sha256.Sum256([]byte(fileContent))
	fileName := fmt.Sprintf("%x", nameHash)[:16] + ".crawljob"
	// add the file name and the folder in which we drop the crawljobs to be picked up
	filePath := filepath.Join(fileDestination, fileName)

	err := os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

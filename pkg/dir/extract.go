package dir

import (
	"path"
	"strings"
	"os"
	"github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
)

// Split the source directory into working directory and source directory
func Split(source string) (src string, pwd string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatal("Unable to resolve working directory")
	}

	getterSource, err := getter.Detect(source, pwd, getter.Detectors)
	if err != nil {
		log.WithError(err).Errorf("Failed to detect source.")
		return src, ""
	}

	if strings.Index(getterSource, "file://") == 0 {
		log.Debug("File source detected. Changing source directory")
		rootPath := strings.Replace(getterSource, "file://", "", 1)

		fi, err := os.Stat(rootPath)
		if err != nil {
			log.WithError(err).Errorf("Failed to read %v", rootPath)
			return src, ""
		}

		if !fi.IsDir() {
			pwd = path.Dir(rootPath)
			src = path.Base(rootPath)
		} else {
			pwd = rootPath
			src = "."
		}

		log.Debugf("New source directory: %v", pwd)
		log.Debugf("New source: %v", src)
	}

	return src, pwd
}
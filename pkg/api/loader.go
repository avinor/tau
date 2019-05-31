package api

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/avinor/tau/pkg/config"
	getter "github.com/hashicorp/go-getter"
)

type Loader struct {
	hashing func(string) string
	pwd     string
	configs []*config.Config
}

type Module struct {
	// Executor
}

func NewLoader() *Loader {
	dir, err := os.Getwd()
	if err != nil {
		log.Panicf("Unable to get working directory: %s", err)
	}

	return &Loader{
		hashing: getDstPath,
		pwd:     dir,
	}
}

func (l *Loader) Loadfile(src string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	dst := fmt.Sprintf(".tau/%s/source", l.hashing(src))
	var err error

	// Try with .hcl and .tau extension if cannot find file
	for _, extsrc := range []string{src, fmt.Sprintf("%s.tau", src), fmt.Sprintf("%s.hcl", src)} {
		client := &getter.Client{
			Ctx:  ctx,
			Src:  extsrc,
			Dst:  dst,
			Pwd:  l.pwd,
			Mode: getter.ClientModeAny,
		}

		if err = client.Get(); err == nil {
			break
		}
	}

	return nil
}

func (l *Loader) GetModules() ([]*Module, error) {
	return nil, nil
}

func getDstPath(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

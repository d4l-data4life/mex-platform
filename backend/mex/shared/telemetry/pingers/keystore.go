package pingers

import (
	"errors"
	"fmt"
	"time"

	"github.com/d4l-data4life/mex/mex/shared/keys"
	"github.com/d4l-data4life/mex/mex/shared/utils"
)

type KeystorePinger struct {
	lastError error
	quit      chan<- struct{}
}

func (p *KeystorePinger) LastError() error {
	return p.lastError
}

func (p *KeystorePinger) Stop() {
	close(p.quit)
}

func NewKeystorePinger(client keys.TokenValidator, updateInterval time.Duration) *KeystorePinger {
	if client == nil {
		panic("Keystore client is nil")
	}

	p := KeystorePinger{
		lastError: fmt.Errorf("Keystore pinger did not run yet"),
	}

	p.quit = utils.ExponentialTicker(time.Second, updateInterval, func() {
		if !client.IsReady() {
			p.lastError = errors.New("remote key store not ready")
		} else {
			p.lastError = nil
		}
	})

	return &p
}

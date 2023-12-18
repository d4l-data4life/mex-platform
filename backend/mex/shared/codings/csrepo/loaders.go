package csrepo

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/d4l-data4life/mex/mex/shared/codings"
	"github.com/d4l-data4life/mex/mex/shared/utils/async"
)

// Function that takes the `config` property of an IndexDef message and returns a promise of the codingset.
type CodingsetLoader func(config *anypb.Any) async.Promise[codings.Codingset]

type loaderEntry struct {
	msg    proto.Message
	loader CodingsetLoader
}

var loaders = []loaderEntry{}

func InstallLoader(msg proto.Message, loader CodingsetLoader) {
	loaders = append(loaders, loaderEntry{msg: msg, loader: loader})
}

func GetLoader(msg *anypb.Any) CodingsetLoader {
	for _, loaderEntry := range loaders {
		if msg.MessageIs(loaderEntry.msg) {
			return loaderEntry.loader
		}
	}
	return nil
}

package remover

import (
	"github.com/figurecode/files-remover/scanner"
)

type Remover interface {
	Execute(files scanner.FoundFiles) error
}

type DebugRemover struct{}

func (d DebugRemover) Execute(files scanner.FoundFiles) error {
	return nil
}

type ActionRemover struct{}

func (a ActionRemover) Execute(scanner.FoundFiles) error {
	return nil
}

func NewRemover(isDemo bool) Remover {
	if isDemo {
		return DebugRemover{}
	}

	return ActionRemover{}
}

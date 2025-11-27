package remover

import (
	"bytes"
	"testing"

	"github.com/figurecode/files-remover/conf"
	"github.com/figurecode/files-remover/scanner"
	"github.com/stretchr/testify/assert"
)

func TestNewRemover(t *testing.T) {
	t.Run("Is Demo mode", func(t *testing.T) {
		remover := NewRemover(conf.Config{IsDemo: true})

		assert.Implements(t, (*Remover)(nil), remover)
		assert.IsType(t, DebugRemover{}, remover)
	})

	t.Run("Is Action mode", func(t *testing.T) {
		remover := NewRemover(conf.Config{IsDemo: false})

		assert.Implements(t, (*Remover)(nil), remover)
		assert.IsType(t, ActionRemover{}, remover)
	})
}

func TestDebugRemover_Execute(t *testing.T) {
	t.Run("Empty files list", func(t *testing.T) {
		remover := DebugRemover{
			outStream: &bytes.Buffer{},
		}

		err := remover.Execute(scanner.FoundFiles{})

		assert.NoError(t, err)
	})

	t.Run("Nil files list", func(t *testing.T) {
		remover := DebugRemover{
			outStream: &bytes.Buffer{},
		}

		err := remover.Execute(nil)

		assert.NoError(t, err)
	})

	t.Run("Non-empty files list", func(t *testing.T) {
		files := scanner.FoundFiles{
			"/tmp/file1.txt": 0,
			"/tmp/file2.txt": 0,
		}

		remover := DebugRemover{
			outStream: &bytes.Buffer{},
		}
		err := remover.Execute(files)

		assert.NoError(t, err)
	})
}

func TestActionRemover_Execute(t *testing.T) {
	t.Run("Empty files list", func(t *testing.T) {
		remover := ActionRemover{}

		err := remover.Execute(scanner.FoundFiles{})

		assert.NoError(t, err)
	})

	t.Run("Nil files list", func(t *testing.T) {
		remover := ActionRemover{}

		err := remover.Execute(nil)

		assert.NoError(t, err)
	})

	t.Run("Non-empty files list", func(t *testing.T) {
		files := scanner.FoundFiles{
			"/tmp/file1.txt": 0,
			"/tmp/file2.txt": 0,
		}

		remover := ActionRemover{}
		err := remover.Execute(files)

		assert.NoError(t, err)
	})
}

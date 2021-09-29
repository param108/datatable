package keybindings

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type KeyBinding struct {
	V   string
	Key gocui.Key
	Mod gocui.Modifier
	F   func(*gocui.Gui, *gocui.View) error
}

func (kb KeyBinding) Hash() string {
	return fmt.Sprintf("%d_%d", kb.Key, kb.Mod)
}

type KeyStore struct {
	g *gocui.Gui
	M map[string]KeyBinding
}

func NewKeyStore(g *gocui.Gui) *KeyStore {
	return &KeyStore{
		g: g,
		M: map[string]KeyBinding{},
	}
}

// Adds a key and checks for collisions
func (ks *KeyStore) AddKey(viewName string, key gocui.Key, mod gocui.Modifier,
	fn func(*gocui.Gui, *gocui.View) error) error {
	kb := KeyBinding{
		V:   viewName,
		Key: key,
		Mod: mod,
		F:   fn,
	}

	hash := kb.Hash()

	if _, ok := ks.M[hash]; ok {
		return errors.New("already found")
	}
	log.Infof("New Binding %s", viewName)
	if err := ks.g.SetKeybinding(viewName, key, mod, fn); err != nil {
		return errors.Wrap(err, "failed to bind")
	}

	ks.M[hash] = kb
	return nil
}

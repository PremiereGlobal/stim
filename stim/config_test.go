package stim

import (
	"testing"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/spf13/viper"
	"gotest.tools/assert"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestSimpleGetString(t *testing.T) {
	stim := &Stim{
		config: viper.New(),
	}
	stim.config.SetDefault("TEST-VALUE-ONE", "ONE")

	simple := stim.ConfigGetString("test-value-one")
	simple2 := stim.ConfigGetString("test.value-one")
	none := stim.ConfigGetString("test.value-two")

	assert.Equal(t, "ONE", simple, "Values not Equal")
	assert.Equal(t, "ONE", simple2, "Values not Equal")
	assert.Equal(t, "", none, "Values not Equal")
}

func TestSimpleGetRaw(t *testing.T) {
	stim := &Stim{
		config: viper.New(),
	}
	stim.config.SetDefault("TEST-VALUE-ONE", "ONE")

	simple := stim.ConfigGetRaw("test-value-one")
	simple2 := stim.ConfigGetRaw("test.value-one")
	none := stim.ConfigGetRaw("test.value-two")

	assert.Equal(t, "ONE", simple, "Values not Equal")
	assert.Equal(t, "ONE", simple2, "Values not Equal")
	assert.Equal(t, nil, none, "Values not Equal")
}

func TestSimpleGetBool(t *testing.T) {
	stim := &Stim{
		config: viper.New(),
		log:    stimlog.GetLogger(),
	}
	stim.config.SetDefault("TEST-VALUE-TRUE", true)
	stim.config.SetDefault("TEST-VALUE-TRUE-STRING", "true")
	stim.config.SetDefault("TEST-VALUE-TRUE-STRING2", "bad")
	stim.config.SetDefault("TEST-VALUE-TRUE-INT", 100)

	simple := stim.ConfigGetBool("test-value-true")
	simple2 := stim.ConfigGetBool("test.value-true")
	simple3 := stim.ConfigGetBool("test.value-true.string")
	simple4 := stim.ConfigGetBool("test.value-true.string2")
	simple5 := stim.ConfigGetBool("test.value-true.int")
	none := stim.ConfigGetBool("test.value-false")

	assert.Equal(t, true, simple, "Values not Equal")
	assert.Equal(t, true, simple2, "Values not Equal")
	assert.Equal(t, true, simple3, "Values not Equal")
	assert.Equal(t, false, simple4, "Values not Equal")
	assert.Equal(t, false, simple5, "Values not Equal")
	assert.Equal(t, false, none, "Values not Equal")
}

func TestHasValue(t *testing.T) {
	stim := &Stim{
		config: viper.New(),
	}
	stim.config.SetDefault("TEST-VALUE-TRUE", true)

	yes := stim.ConfigHasValue("test-value-true")
	yes2 := stim.ConfigHasValue("test.value-true")
	yes3 := stim.ConfigHasValue("test.value.true")
	no := stim.ConfigHasValue("test-value-false")

	assert.Equal(t, true, yes, "Values not Equal")
	assert.Equal(t, true, yes2, "Values not Equal")
	assert.Equal(t, true, yes3, "Values not Equal")
	assert.Equal(t, false, no, "Values not Equal")
}

package parallel

import "testing"

func TestDataChannel(t *testing.T) {
	p := New()
	t.Run("Initialization", func(t *testing.T) {
		err := p.NewDataChannel("TestChannel")
		testNil(t, err)
		testNotNil(t, p.dataChannels["TestChannel"])
	})
	t.Run("Closure", func(t *testing.T) {
		err := p.CloseDataChannel("TestChannel")
		testNil(t, err)
		_, exists := p.dataChannels["TestChannel"]
		testAssert(t, exists == false, "data channel still exists")
	})
}

package parallel

import "testing"

func TestDataChannel(t *testing.T) {
	p := New()
	t.Run("Initialization", func(t *testing.T) {
		err := p.NewDataChannel("TestChannel")
		test_Nil(t, err)
		test_NotNil(t, p.dataChannels["TestChannel"])
	})
	t.Run("Closure", func(t *testing.T) {
		err := p.CloseDataChannel("TestChannel")
		test_Nil(t, err)
		_, exists := p.dataChannels["TestChannel"]
		test_Assert(t, exists == false, "data channel still exists")
	})
}

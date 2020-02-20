package parallel

import "testing"

func TestDataChannel(t *testing.T) {
	p := New()
	t.Run("Initialization", func(t *testing.T) {
		err := p.NewDataChannel("TestChannel")
		Test_Nil(t, err)
		Test_NotNil(t, p.dataChannels["TestChannel"])
	})
	t.Run("Closure", func(t *testing.T) {
		err := p.CloseDataChannel("TestChannel")
		Test_Nil(t, err)
		_, exists := p.dataChannels["TestChannel"]
		Test_Assert(t, exists == false, "data channel still exists")
	})
}

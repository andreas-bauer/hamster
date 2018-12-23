package collect

import "testing"

func TestGet(t *testing.T) {
	client_60 := NewRetryHTTPClient(60, 1)
	bytes, err := client_60.Get("https://api.github.com/users/michaeldorner")
	if err != nil || bytes == nil {
		t.Errorf(err.Error())
	}

	client_1 := NewRetryHTTPClient(1, 2)
	bytes, err = client_1.Get("noreallserver")
	if err != ErrMaxRetries || bytes != nil {
		t.Errorf(err.Error())
	}
}

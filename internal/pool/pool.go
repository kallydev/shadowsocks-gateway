package pool

import "sync"

var instance = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func Get() []byte {
	return instance.Get().([]byte)
}

func Put(buffer []byte) {
	instance.Put(buffer)
}

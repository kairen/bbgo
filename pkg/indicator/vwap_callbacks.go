// Code generated by "callbackgen -type VWAP"; DO NOT EDIT.

package indicator

import ()

func (V *VWAP) OnUpdate(cb func(value float64)) {
	V.UpdateCallbacks = append(V.UpdateCallbacks, cb)
}

func (V *VWAP) EmitUpdate(value float64) {
	for _, cb := range V.UpdateCallbacks {
		cb(value)
	}
}

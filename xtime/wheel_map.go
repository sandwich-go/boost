package xtime

//go:generate gotemplate -outfmt gen_%v "../../base/container/templates/syncmap" "WheelMap(time.Duration,*Wheel)"

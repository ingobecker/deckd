package core

type sampleRingState struct {
	FillSize  int
	WriteSize int
	ReadWrap  bool
	WriteWrap bool
}

type SampleRing struct {
	readP   int
	writeP  int
	samples []float64
	size    int
}

// NewSampleRing creates a new SampleRing which is a ringbuffer
// implementation, that can be used in a locking free manner.
// It uses a pair of pointers where the write proccess
// is only allowed to modify the write pointer and read from the
// read pointer while the read proccess is only allowed to modify
// the read pointer and read the write pointer.
// This makes the implementation more convoluted but prevents
// possible race conditions which would occour if the read and
// write functions would modify the same data without proper
// synchronization. Locking free modification is required because
// this buffer is meant to be used in an audio application which
// doesn't allow blocking operations while manipulating the buffer.
// This only holds true if there are no concurrent write or read
// operations.
// Internally the array backing the ring is one element larger than
// the size given to this function. This is because the state of
// readP == writeP signals an empty buffer. By this definition
// the field writeP points to is always empty and the first field
// a write() call will write to. If we try to fill the whole array
// we always have one element less than the backing arrays size is.
// To compensate for that we simply create a backing array one field
// lager than specified.

func NewSampleRing(size int) *SampleRing {
	return &SampleRing{
		samples: make([]float64, size+1),
		size:    size,
	}
}

func (r *SampleRing) computeState(readP int, writeP int) *sampleRingState {
	state := &sampleRingState{}
	switch {
	case readP < writeP:
		state.ReadWrap = false
		state.WriteWrap = true
		state.FillSize = writeP - readP
	case readP > writeP:
		state.ReadWrap = true
		state.WriteWrap = false
		state.FillSize = r.size - readP + writeP + 1
	default:
		state.FillSize = 0
		state.ReadWrap = true
		state.WriteWrap = true
	}

	if readP == 0 {
		state.ReadWrap = false
		state.WriteWrap = false
	}
	state.WriteSize = r.size - state.FillSize
	return state
}

func (r *SampleRing) write(samples []float64) int {
	rP, wP := r.readP, r.writeP
	state := r.computeState(rP, wP)
	requested := len(samples)
	upper := r.size - wP + 1
	if requested > state.WriteSize {
		requested = state.WriteSize
	}
	if state.WriteWrap {
		if requested > upper {
			lower := requested - upper
			copy(r.samples[wP:], samples)
			copy(r.samples[:lower], samples[upper:upper+lower])
			r.writeP = lower
		} else {
			copy(r.samples[wP:wP+requested], samples[:requested])
			r.writeP = wP + requested
		}
	} else {
		copy(r.samples[wP:wP+requested], samples[:requested])
		r.writeP = wP + requested
	}
	return requested
}

func (r *SampleRing) read(samples []float64) int {
	rP, wP := r.readP, r.writeP
	state := r.computeState(rP, wP)
	requested := len(samples)
	upper := r.size - rP + 1
	if requested > state.FillSize {
		requested = state.FillSize
	}
	if state.ReadWrap {
		if requested > upper {
			lower := requested - upper
			copy(samples, r.samples[rP:])
			copy(samples[upper:upper+lower], r.samples[:lower])
			r.readP = lower
		} else {
			copy(samples, r.samples[rP:rP+requested])
			r.readP = rP + requested
		}
	} else {
		copy(samples, r.samples[rP:rP+requested])
		r.readP = rP + requested
	}
	return requested
}

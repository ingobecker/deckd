package core

import (
	"reflect"
	"testing"
)

func TestReadFromEmpty(t *testing.T) {
	r := NewSampleRing(6)
	out := []float64{0, 0}

	expectedSamples := 0
	actualSamples := r.read(out)

	if expectedSamples != actualSamples {
		t.Errorf("Samples returned: %v; expected %v", expectedSamples, actualSamples)
	}
}

func TestWriteMore(t *testing.T) {
	r := NewSampleRing(2)
	in := []float64{42, 23, 5}

	expectedSamples := 2
	actualSamples := r.write(in)

	if actualSamples != expectedSamples {
		t.Errorf("Number of samples written: %v; expected %v", actualSamples, expectedSamples)
	}
}

func TestWriteFull(t *testing.T) {
	r := NewSampleRing(2)
	in := []float64{42, 23}

	expectedSamples := 2
	actualSamples := r.write(in)

	if actualSamples != expectedSamples {
		t.Errorf("Number of samples written: %v; expected %v", actualSamples, expectedSamples)
	}
}

func TestWriteLess(t *testing.T) {
	r := NewSampleRing(2)
	in := []float64{42}

	expectedSamples := 1
	actualSamples := r.write(in)

	if actualSamples != expectedSamples {
		t.Errorf("Number of samples written: %v; expected %v", actualSamples, expectedSamples)
	}
}

func TestReadLess(t *testing.T) {
	r := NewSampleRing(2)
	in := []float64{42, 23}
	out := []float64{0}

	expectedSamples := 1
	r.write(in)
	actualSamples := r.read(out)

	if actualSamples != expectedSamples {
		t.Errorf("Number of samples red: %v; expected %v", actualSamples, expectedSamples)
	}

	if out[0] != in[0] {
		t.Errorf("Red %v; expected %v", out[0], in[0])
	}
}

func TestReadFull(t *testing.T) {
	r := NewSampleRing(2)
	in := []float64{42, 23}
	out := []float64{0, 0}

	expectedSamples := 2
	r.write(in)
	actualSamples := r.read(out)

	if actualSamples != expectedSamples {
		t.Errorf("Number of samples red: %v; expected %v", actualSamples, expectedSamples)
	}

	if !reflect.DeepEqual(in, out) {
		t.Errorf("Red %v; expected %v", out, in)
	}
}

func TestReadMore(t *testing.T) {
	r := NewSampleRing(2)
	in := []float64{42, 23}
	out := []float64{0, 0, 0}

	expectedSamples := 2
	r.write(in)
	actualSamples := r.read(out)

	if actualSamples != expectedSamples {
		t.Errorf("Number of samples red: %v; expected %v", actualSamples, expectedSamples)
	}

	if !reflect.DeepEqual(in, out[:2]) {
		t.Errorf("Red %v; expected %v", out, in)
	}
}

func TestWriteWithOffset(t *testing.T) {
	r := NewSampleRing(6)
	in := []float64{42, 23, 5, 1}
	out := []float64{0, 0, 0, 0}

	expectedSamples := 4
	actualSamplesW := r.write(in[0:2])
	actualSamplesW += r.write(in[2:4])
	actualSamplesR := r.read(out)

	if actualSamplesR != expectedSamples {
		t.Errorf("Number of samples red: %v; expected %v", actualSamplesR, expectedSamples)
	}

	if actualSamplesW != expectedSamples {
		t.Errorf("Number of samples written: %v; expected %v", actualSamplesW, expectedSamples)
	}

	if !reflect.DeepEqual(in, out) {
		t.Errorf("Red %v; expected %v", out, in)
	}
}

func TestReadWithOffset(t *testing.T) {
	r := NewSampleRing(6)
	in := []float64{42, 23, 5, 1}
	out := []float64{0, 0, 0, 0}

	expectedSamples := 4
	actualSamplesW := r.write(in)
	actualSamplesR := r.read(out[:2])
	actualSamplesR += r.read(out[2:4])

	if actualSamplesW != expectedSamples {
		t.Errorf("Number of samples written: %v; expected %v", actualSamplesW, expectedSamples)
	}

	if actualSamplesR != expectedSamples {
		t.Errorf("Number of samples red: %v; expected %v", actualSamplesR, expectedSamples)
	}

	if !reflect.DeepEqual(in, out) {
		t.Errorf("Red %v; expected %v", out, in)
	}
}

func TestWriteWrapTable(t *testing.T) {
	testCases := []struct {
		name string
		in   []float64
		out  []float64
	}{
		{"WriteWrapOnWrappable", []float64{42, 23, 5}, []float64{0, 0, 0}},
		{"WriteNoWrapOnWrappable", []float64{42}, []float64{0}},
		{"WriteWrapOnWrappableReadNoWrap", []float64{42, 23}, []float64{0}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewSampleRing(3)
			pre_in := []float64{3, 2, 1}
			pre_out := []float64{0, 0, 0}

			r.write(pre_in)
			r.read(pre_out)

			expectedSamples := len(tc.in)
			actualSamplesW := r.write(tc.in)

			actualSamplesR := r.read(tc.out)
			if actualSamplesW != expectedSamples {
				t.Errorf("Number of samples written: %v; expected %v", actualSamplesW, expectedSamples)
			}

			if !reflect.DeepEqual(tc.in[:actualSamplesR], tc.out) {
				t.Errorf("Red %v; expected %v", tc.out, tc.in[:actualSamplesR])
			}
		})
	}
}

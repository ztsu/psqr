// Package psqr provides dynamic calculation of quantiles without storing
// observations
// See http://www.cs.wustl.edu/~jain/papers/ftp/psqr.pdf for further reading
package psqr

import (
	"errors"
	"fmt"
	"math"
)

const points = 5

// Type, that satisfies psqr.Valuer can be use for calculations
type Valuer interface {
	Value() float64
}

// marker represent one of five points for estminate quantile
// not exported
type marker struct {
	height           float64
	actualPosition   int
	desiredPosition  float64
	desiredIncrement float64
}

// PSQR is a main type, holding markers and do calculations
type PSQR struct {
	markers      [points]marker
	observations int
	quantile     float64
}

func (psqr *PSQR) AddValue(data Interface) {
	if psqr.observations < points {
		psqr.init(data)
	} else {
		psqr.addValue(data)
	}
	psqr.observations += 1
}

func (psqr *PSQR) init(data Interface) {
	psqr.markers[psqr.observations] = marker{
		height:         data.Value(),
		actualPosition: psqr.observations + 1,
	}
	switch psqr.observations {
	case 0:
		psqr.markers[psqr.observations] = marker{
			desiredPosition:  1,
			desiredIncrement: 0,
		}
	case 1:
		psqr.markers[psqr.observations] = marker{
			desiredPosition:  1 + 2*psqr.quantile,
			desiredIncrement: psqr.quantile / 2,
		}
	case 2:
		psqr.markers[psqr.observations] = marker{
			desiredPosition:  1 + 4*psqr.quantile,
			desiredIncrement: psqr.quantile,
		}
	case 3:
		psqr.markers[psqr.observations] = marker{
			desiredPosition:  3 + 2*psqr.quantile,
			desiredIncrement: (1 + psqr.quantile) / 2,
		}
	case 4:
		psqr.markers[psqr.observations] = marker{
			desiredPosition:  5,
			desiredIncrement: 1,
		}
	}
}

func (psqr *PSQR) addValue(data Interface) {
	var cell int

	switch {
	case data.Value() < psqr.markers[0].height:
		psqr.markers[1].height = data.Value()
		cell = 1
	case data.Value() < psqr.markers[1].height:
		cell = 1
	case data.Value() < psqr.markers[2].height:
		cell = 2
	case data.Value() < psqr.markers[3].height:
		cell = 3
	case data.Value() < psqr.markers[4].height:
		cell = 4
	case data.Value() > psqr.markers[4].height:
		psqr.markers[4].height = data.Value()
		cell = 4
	}

	for i := cell; i < len(psqr.markers); i++ {
		psqr.markers[i].actualPosition += 1
	}
	for i := 0; i < len(psqr.markers); i++ {
		psqr.markers[i].desiredPosition += psqr.markers[i].desiredIncrement
	}

	for i := 1; i < 4; i++ {
		offset := psqr.markers[i].desiredPosition - float64(psqr.markers[i].actualPosition)
		if (offset >= 1 && (psqr.markers[i+1].actualPosition-psqr.markers[i].actualPosition) > 1) ||
			(offset <= -1 && (psqr.markers[i-1].actualPosition-psqr.markers[i].actualPosition < -1)) {
			offset = math.Abs(offset)
			height := psqr.parabolic(offset, i)
			if psqr.markers[i-1].height < height &&
				height < psqr.markers[i+1].height {
				psqr.markers[i].height = height
			} else {
				psqr.markers[i].height = psqr.linear(offset, i)
			}
			psqr.markers[i].actualPosition = psqr.markers[i].actualPosition + int(offset)
		}
	}
}

func (psqr *PSQR) parabolic(d float64, i int) float64 {
	// TODO:
	return 0.0
}

func (psqr *PSQR) linear(d float64, i int) float64 {
	// TODO:
	return 0.0
}

func (psqr *PSQR) Quantile() (float64, error) {
	if psqr.observations < 5 {
		return 0.0, errors.New(
			fmt.Sprintf(
				"cannot estimate quantile: ",
				"need at least: %d, got: %d\n",
				points,
				psqr.observations,
			),
		)
	}
	return psqr.quantile, nil
}

func New(percentile float64) *PSQR {
	return &PSQR{
		[5]marker{},
		0,
		percentile,
	}
}

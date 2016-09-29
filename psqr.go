// Package psqr provides dynamic calculation of quantiles without storing
// observations
// See http://www.cs.wustl.edu/~jain/papers/ftp/psqr.pdf for further reading
package psqr

import (
	"errors"
	"fmt"
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

func (psqr *PSQR) AddValue(data Valuer) {
	if psqr.observations < points {
		psqr.init(data)
	} else {
		psqr.addValue(data)
	}
	psqr.observations += 1
}

func (psqr *PSQR) init(data Valuer) {
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

func (psqr *PSQR) addValue(data Valuer) {
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

	fmt.Println(cell)
	for i := cell; i < points; i++ {
		psqr.markers[i].actualPosition += 1
	}
	for i := 0; i < points; i++ {
		psqr.markers[i].desiredPosition +=
			psqr.markers[i].desiredIncrement
	}

	for i := 1; i < points-1; i++ {
		actualCurr := psqr.markers[i].actualPosition
		actualNext := psqr.markers[i+1].actualPosition
		actualPrev := psqr.markers[i-1].actualPosition
		//heightNext := psqr.markers[i+1].height
		//heightPrev := psqr.markers[i-1].height
		desired := psqr.markers[i].desiredPosition
		offset := desired - float64(actualCurr)

		if (offset >= 1 && (actualNext-actualCurr) > 1) ||
			(offset <= -1 && (actualPrev-actualCurr < -1)) {

			offset := int(offset)
			//height := psqr.parabolic(offset, i)
			//if (height > heightPrev) && (height < heightNext) {
			//	psqr.markers[i].height += height
			//} else {
			psqr.markers[i].height += psqr.linear(offset, i)
			//}
			psqr.markers[i].actualPosition += int(offset)
		}
	}
}

func (psqr *PSQR) parabolic(d, i int) float64 {
	// TODO:
	return 0.0
}

func (psqr *PSQR) linear(d, i int) float64 {
	heightNext := psqr.markers[i+d].height
	actualCurr := psqr.markers[i].actualPosition
	actualNext := psqr.markers[i+d].actualPosition

	return heightNext * float64(d/(actualNext-actualCurr))
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
	return psqr.markers[3].height, nil
}

func New(quantile float64) *PSQR {
	return &PSQR{
		quantile: quantile,
	}
}

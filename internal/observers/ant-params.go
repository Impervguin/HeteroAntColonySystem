package observers

import "HeteroAntColonySystem/internal/core/colony"

type mapElement struct {
	Alpha float64
	Beta  float64
}

type AntParamsObserver struct {
	m map[uint]mapElement
}

func NewAntParamsObserver(gen uint) *AntParamsObserver {
	return &AntParamsObserver{
		m: make(map[uint]mapElement, gen),
	}
}

var _ colony.ColonyObserver = (*AntParamsObserver)(nil)

func (o *AntParamsObserver) Observe(dto *colony.ColonyObserverDTO) {
	avgAlpha := 0.0
	avgBeta := 0.0
	for _, ant := range dto.Ants {
		avgAlpha += ant.Alpha()
		avgBeta += ant.Beta()
	}
	avgAlpha /= float64(len(dto.Ants))
	avgBeta /= float64(len(dto.Ants))
	o.m[dto.Generation] = mapElement{
		Alpha: avgAlpha,
		Beta:  avgBeta,
	}
}

func (o *AntParamsObserver) Params(gen uint) (alpha, beta float64) {
	params := o.m[gen]
	return params.Alpha, params.Beta
}

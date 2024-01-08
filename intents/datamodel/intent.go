package datamodel

import "encoding/json"

type Intent struct {
	Name                   string    `json:"name,omitempty"`
	Idx                    int       `json:"idx,omitempty"`
	MaximumLoadAverage     int       `json:"maximum_load_average,omitempty"`
	MinimumThroughput      int       `json:"minimum_throughput,omitempty"`
	MaximumUEPerCell       int       `json:"maximum_ue_per_cell,omitempty"`
	MaximumAssociationRate int       `json:"maximum_association_rate,omitempty"`
	Condition              Condition `json:"condition,omitempty"`
	Objective              Objective `json:"objective,omitempty"`
}



type Condition struct {
	When   When   `json:"when"`
	Labels string `json:"labels,omitempty"`
}

type When struct {
	DayOfWeek string   `json:"DayOfWeek"`
	TimeSpan  TimeSpan `json:"time_span"`
}

type TimeSpan struct {
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
}

type Objective struct {
	MinimumCellOffset int `json:"minimum_cell_offset"`
	MaximumCellOffset int `json:"maximum_cell_offset"`
}

func (i Intent) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}

package places

import (
	admission_criteria "cuide/api/resource/admission-criteria"
	attendance_types "cuide/api/resource/attendance-types"
	reference_ways "cuide/api/resource/reference-ways"
	"cuide/api/resource/regionals"
	"cuide/api/resource/segments"
	service_types "cuide/api/resource/service-types"
)

type DTO struct {
	ID                uint8                                    `json:"id"`
	Name              string                                   `json:"name"`
	Address           string                                   `json:"address"`
	PhoneNumber       string                                   `json:"phone_number"`
	Website           string                                   `json:"website"`
	Observations      string                                   `json:"observations"`
	GoogleMapsLink    string                                   `json:"google_maps_link"`
	ServiceType       service_types.ServiceType                `json:"service_type"`
	Segment           segments.Segment                         `json:"segment"`
	Regionals         regionals.Regionals                      `json:"regionals"`
	ReferenceWays     reference_ways.ReferenceWays             `json:"reference_ways"`
	AdmissionCriteria admission_criteria.AdmissionCriteriaList `json:"admission_criteria"`
	AttendanceType    attendance_types.AttendanceTypes         `json:"attendance_types"`
}

type Form struct {
	Name                 string `json:"name"                   form:"required,max=255"`
	Address              string `json:"address"                form:"required,max=255"`
	PhoneNumber          string `json:"phone_number"           form:"max=255"`
	Website              string `json:"website"                form:"max=255"`
	Observations         string `json:"observations"           form:"max=255"`
	GoogleMapsLink       string `json:"google_maps_link"       form:"required,max=255"`
	ServiceTypeID        uint   `json:"service_type_id"        form:"required,min=1"`
	SegmentID            uint   `json:"segment_id"             form:"required,min=1"`
	RegionalIDs          []uint `json:"regional_ids"           form:"required,min=1"`
	ReferenceWayIDs      []uint `json:"reference_ways_ids"     form:"required,min=1"`
	AdmissionCriteriaIDs []uint `json:"admission_criteria_ids" form:"required,min=1"`
	AttendanceTypeIDs    []uint `json:"attendance_type_ids"    form:"required,min=1"`
}

type Place struct {
	ID                uint8
	Name              string
	Address           string
	PhoneNumber       string
	Website           string
	Observations      string
	GoogleMapsLink    string
	ServiceType       service_types.ServiceType
	Segment           segments.Segment
	Regionals         regionals.Regionals
	ReferenceWays     reference_ways.ReferenceWays
	AdmissionCriteria admission_criteria.AdmissionCriteriaList
	AttendanceType    attendance_types.AttendanceTypes
}

type Places []*Place

func (r *Place) ToDto() *DTO {
	return &DTO{
		ID:                r.ID,
		Name:              r.Name,
		Address:           r.Address,
		PhoneNumber:       r.PhoneNumber,
		Website:           r.Website,
		Observations:      r.Observations,
		GoogleMapsLink:    r.GoogleMapsLink,
		ServiceType:       r.ServiceType,
		Segment:           r.Segment,
		Regionals:         r.Regionals,
		ReferenceWays:     r.ReferenceWays,
		AdmissionCriteria: r.AdmissionCriteria,
		AttendanceType:    r.AttendanceType,
	}
}

func (rgs Places) ToDto() []*DTO {
	dtos := make([]*DTO, len(rgs))

	for i, v := range rgs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *Form) ToModel() Place {
	rs := make(regionals.Regionals, len(f.RegionalIDs))
	for i, r := range f.RegionalIDs {
		rs[i] = &regionals.Regional{
			ID: uint8(r),
		}
	}

	rws := make(reference_ways.ReferenceWays, len(f.ReferenceWayIDs))
	for i, r := range f.ReferenceWayIDs {
		rws[i] = &reference_ways.ReferenceWay{
			ID: uint8(r),
		}
	}

	acs := make(admission_criteria.AdmissionCriteriaList, len(f.AdmissionCriteriaIDs))
	for i, ac := range f.AdmissionCriteriaIDs {
		acs[i] = &admission_criteria.AdmissionCriteria{
			ID: uint8(ac),
		}
	}

	atts := make(attendance_types.AttendanceTypes, len(f.AttendanceTypeIDs))
	for i, att := range f.AttendanceTypeIDs {
		atts[i] = &attendance_types.AttendanceType{
			ID: uint8(att),
		}
	}

	return Place{
		Name:           f.Name,
		Address:        f.Address,
		PhoneNumber:    f.PhoneNumber,
		Website:        f.Website,
		Observations:   f.Observations,
		GoogleMapsLink: f.GoogleMapsLink,
		ServiceType: service_types.ServiceType{
			ID: uint8(f.ServiceTypeID),
		},
		Segment: segments.Segment{
			ID: uint8(f.SegmentID),
		},
		Regionals:         rs,
		ReferenceWays:     rws,
		AdmissionCriteria: acs,
		AttendanceType:    atts,
	}
}

package places

import (
	"cuide/api/resource/regionals"
	"cuide/api/resource/segments"
	service_types "cuide/api/resource/service-types"
)

type DTO struct {
	ID                  uint8                     `json:"id"`
	Name                string                    `json:"name"`
	Address             string                    `json:"address"`
	PhoneNumber         string                    `json:"phone_number"`
	Website             string                    `json:"website"`
	Observations        string                    `json:"observations"`
	GoogleMapsLink      string                    `json:"google_maps_link"`
	GoogleMapsEmbedLink string                    `json:"google_maps_embed_link"`
	AdmissionCriteria   string                    `json:"admission_criteria"`
	ReferenceWay        string                    `json:"reference_ways"`
	AttendanceType      string                    `json:"attendance_types"`
	ServiceType         service_types.ServiceType `json:"service_type"`
	Segment             segments.Segment          `json:"segment"`
	Regionals           regionals.Regionals       `json:"regionals"`
}

type Form struct {
	Name                string `json:"name"                   form:"required,max=2500"`
	Address             string `json:"address"                form:"required,max=2500"`
	PhoneNumber         string `json:"phone_number"           form:"max=2500"`
	Website             string `json:"website"                form:"max=2500"`
	Observations        string `json:"observations"`
	GoogleMapsLink      string `json:"google_maps_link"       form:"required"`
	GoogleMapsEmbedLink string `json:"google_maps_embed_link" form:"required"`
	AdmissionCriteria   string `json:"admission_criteria"     form:"required"`
	ReferenceWay        string `json:"reference_ways"         form:"required"`
	AttendanceType      string `json:"attendance_types"       form:"required"`
	ServiceTypeID       uint   `json:"service_type_id"        form:"required,min=1"`
	SegmentID           uint   `json:"segment_id"             form:"required,min=1"`
	RegionalIDs         []uint `json:"regional_ids"           form:"required,min=1"`
}

type Place struct {
	ID                  uint8
	Name                string
	Address             string
	PhoneNumber         string
	Website             string
	Observations        string
	GoogleMapsLink      string
	GoogleMapsEmbedLink string
	AdmissionCriteria   string
	ReferenceWay        string
	AttendanceType      string
	ServiceType         service_types.ServiceType
	Segment             segments.Segment
	Regionals           regionals.Regionals
}

type Places []*Place

type Filters struct {
	ServiceTypes []uint8
	Segments     []uint8
	Regionals    []uint8
	Name         string
}

type PaginationMetadata struct {
	Places   []*DTO `json:"places"`
	Metadata struct {
		TotalPlaces uint8 `json:"total_places"`
		Pages       uint8 `json:"pages"`
	} `json:"metadata"`
}

func (r *Place) ToDto() *DTO {
	return &DTO{
		ID:                  r.ID,
		Name:                r.Name,
		Address:             r.Address,
		PhoneNumber:         r.PhoneNumber,
		Website:             r.Website,
		Observations:        r.Observations,
		GoogleMapsLink:      r.GoogleMapsLink,
		GoogleMapsEmbedLink: r.GoogleMapsEmbedLink,
		AdmissionCriteria:   r.AdmissionCriteria,
		ReferenceWay:        r.ReferenceWay,
		AttendanceType:      r.AttendanceType,
		ServiceType:         r.ServiceType,
		Segment:             r.Segment,
		Regionals:           r.Regionals,
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

	return Place{
		Name:                f.Name,
		Address:             f.Address,
		PhoneNumber:         f.PhoneNumber,
		Website:             f.Website,
		Observations:        f.Observations,
		GoogleMapsLink:      f.GoogleMapsLink,
		GoogleMapsEmbedLink: f.GoogleMapsEmbedLink,
		AdmissionCriteria:   f.AdmissionCriteria,
		AttendanceType:      f.AttendanceType,
		ReferenceWay:        f.ReferenceWay,
		ServiceType: service_types.ServiceType{
			ID: uint8(f.ServiceTypeID),
		},
		Segment: segments.Segment{
			ID: uint8(f.SegmentID),
		},
		Regionals: rs,
	}
}

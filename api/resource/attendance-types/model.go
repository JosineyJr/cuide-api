package attendance_types

type DTO struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type Form struct {
	Name string `json:"name" form:"required,max=255"`
}

type AttendanceType struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type AttendanceTypes []*AttendanceType

func (r *AttendanceType) ToDto() *DTO {
	return &DTO{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (atts AttendanceTypes) ToDto() []*DTO {
	dtos := make([]*DTO, len(atts))

	for i, v := range atts {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *Form) ToModel() AttendanceType {
	return AttendanceType{
		Name: f.Name,
	}
}

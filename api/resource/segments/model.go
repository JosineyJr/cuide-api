package segments

type DTO struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type Form struct {
	Name string `json:"name" form:"required,max=255"`
}

type Segment struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type Segments []*Segment

func (r *Segment) ToDto() *DTO {
	return &DTO{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (rgs Segments) ToDto() []*DTO {
	dtos := make([]*DTO, len(rgs))

	for i, v := range rgs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *Form) ToModel() Segment {
	return Segment{
		Name: f.Name,
	}
}

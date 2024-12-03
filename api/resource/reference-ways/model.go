package reference_ways

type DTO struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type Form struct {
	Name string `json:"name" form:"required,max=255"`
}

type ReferenceWay struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type ReferenceWays []*ReferenceWay

func (r *ReferenceWay) ToDto() *DTO {
	return &DTO{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (rgs ReferenceWays) ToDto() []*DTO {
	dtos := make([]*DTO, len(rgs))

	for i, v := range rgs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *Form) ToModel() ReferenceWay {
	return ReferenceWay{
		Name: f.Name,
	}
}

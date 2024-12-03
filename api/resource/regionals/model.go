package regionals

type DTO struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type Form struct {
	Name string `json:"name" form:"required,max=255"`
}

type Regional struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type Regionals []*Regional

func (r *Regional) ToDto() *DTO {
	return &DTO{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (rgs Regionals) ToDto() []*DTO {
	dtos := make([]*DTO, len(rgs))

	for i, v := range rgs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *Form) ToModel() Regional {
	return Regional{
		Name: f.Name,
	}
}

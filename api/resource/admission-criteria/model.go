package admission_criteria

type DTO struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type Form struct {
	Name string `json:"name" form:"required,max=255"`
}

type AdmissionCriteria struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type AdmissionCriteriaList []*AdmissionCriteria

func (r *AdmissionCriteria) ToDto() *DTO {
	return &DTO{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (rgs AdmissionCriteriaList) ToDto() []*DTO {
	dtos := make([]*DTO, len(rgs))

	for i, v := range rgs {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *Form) ToModel() AdmissionCriteria {
	return AdmissionCriteria{
		Name: f.Name,
	}
}

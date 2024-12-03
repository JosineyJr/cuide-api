package service_types

type DTO struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type Form struct {
	Name string `json:"name" form:"required,max=255"`
}

type ServiceType struct {
	ID   uint8  `json:"id"`
	Name string `json:"name"`
}

type ServiceTypes []*ServiceType

func (r *ServiceType) ToDto() *DTO {
	return &DTO{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (sts ServiceTypes) ToDto() []*DTO {
	dtos := make([]*DTO, len(sts))

	for i, v := range sts {
		dtos[i] = v.ToDto()
	}

	return dtos
}

func (f *Form) ToModel() ServiceType {
	return ServiceType{
		Name: f.Name,
	}
}

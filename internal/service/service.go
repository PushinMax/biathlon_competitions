package service

type Service struct{

}

func New() *Service {
	return &Service{}
}

func (s * Service) Execute() (string, error) {
	return "", nil
}
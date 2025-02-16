package vault

type AuthLog struct {
	file string
	Current  *Credentials `json:"current"`
	Previous []*Credentials `json:"previous"`
}


func (a *AuthLog) ReadLogFile() ([]byte, error){

	return nil, nil
}

func (a *AuthLog) WriteLogFile() ([]byte, error){
	return nil, nil
}
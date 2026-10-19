package config

type configs struct{
	DbConfig *DbConfig
}

func CarregarConfig() (*configs, error){

	dbConfig, err := carregarConfigsDb()
	if err != nil{
		return nil, err
	}

	return &configs{
		DbConfig: dbConfig,
	}, nil
}
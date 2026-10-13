package config

type configs struct{
	DbConfig *DbConfig
}

func CarregarConfig() *configs{

	dbConfig := carregarConfigsDb()

	return &configs{
		DbConfig: dbConfig,
	}
}
Esse arquivo descreve as abstracoes criadas e suas motivacoes.

## ConfigManager

o sistema funciona em cima de partes menores e uma parte maior. Existirao arquivos "exemplo_config.go" que sao responsaveis por criar structs como: 

```Go
//db_config.go

package config

import (
	"time"
)

type DbConfig struct{
	Db_URL string
	MaxConn int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	ConnTimeout time.Duration
}
```

Seus métodos serao:

```go

func carregarConfigsDb() *DbConfig{
	return &DbConfig{
		// campos}
	}

```

E entao existira um arquivo "config_manager.go" com uma struct de configs. Seguindo esse modelo:

```go

//config_manager.go

package config

type Configs struct{
	DbConfig *DbConfig
}
```

Seus metodos serao:

```go

func CarregarConfig() *Configs{

	dbConfig := carregarConfigsDb()

	return &configs{
		DbConfig: dbConfig
	}
}

```
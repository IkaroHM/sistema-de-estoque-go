Esse arquivo descreve as abstracoes criadas e suas motivacoes.

## ConfigManager

o sistema funciona em cima de partes menores e uma parte maior. Existirao arquivos "exemplo_config.go" que sao responsaveis por criar structs de configuracoes.

Suas funcoes serao:

- carregarConfigModulo() (*structDeConfig, error)
Observase que "Modulo" é uma generalizacao, ou seja, ele pode ser "Db", "log", etc

E entao existira um arquivo "config_manager.go" com uma struct de configs.

Suas funcoes serao:

- carregarConfig() (*Configs, error)

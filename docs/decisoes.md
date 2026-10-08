# 📋 Gestor de Estoque — Documento de Arquitetura e Decisões

> **Versão:** 1.0 (MVP)

---

## 1. Visão Geral

Sistema de gestão de estoque para uso **pessoal e portfólio**. O objetivo é permitir a entrada e saída de produtos e controle de quantidade.

**Público-alvo:** 1 usuário (o próprio desenvolvedor).
**Contexto:** Projeto de aprendizado em Go, com foco em arquitetura limpa (HSR) e boas práticas.

---

## 2. Requisitos Funcionais

### 2.1 Entrada e Saída de Produtos

| Funcionalidade                   | Descrição                                                                    |
| :------------------------------- | :--------------------------------------------------------------------------- |
| **Criar produto**                | Cadastrar produto com ID, nome, descrição, preço (em centavos) e quantidade. |
| **Listar produtos**              | Retornar todos os produtos cadastrados.                                      |
| **Buscar por ID**                | Retornar um produto específico pelo ID.                                      |
| **Mudar nome**                   | Alterar o nome de um produto existente pelo ID.                              |
| **Deletar produto**              | Remover um produto do estoque.                                               |
| **Vender (subtrair quantidade)** | Reduzir a quantidade de um produto, validando se há estoque suficiente.      |
| **Repor (adicionar quantidade)** | Aumentar a quantidade de um produto.                                         |


### 2.2 Status do Produto

| Status | Condição |
| :--- | :--- |
| **Regular** | Quantidade acima do limite mínimo. |
| **Baixo** | Quantidade abaixo do limite mínimo (ex: < 5 unidades). |
| **Esgotado** | Quantidade igual a zero. |

### 2.3 Modelo de Domínio

```text

Estoque
├── id: int
├── nome: string
├── descricao: string
├── preco: int (centavos)
├── quantidade: int

```

---

## 3. Requisitos Não Funcionais

### 3.1 Disponibilidade e Performance

| Requisito | Meta |
| :--- | :--- |
| **Usuários simultâneos** | 1 |
| **Requisições por segundo** | Máximo 10 req/s |
| **Latência máxima** | 300ms |
| **Disponibilidade** | Leve e simples de manter |
| **Contexto de uso** | Pessoal e portfólio |

### 3.2 Segurança

| Item | Decisão |
| :--- | :--- |
| **Autenticação** | Login com sessão via cookie (1 usuário). |
| **Hash de senha** | bcrypt (V1). |
| **HTTPS** | Fornecido pela hospedagem (Fly.io/Render). |
| **Validação de input** | Manual no Service + Handler. |
| **Rate limiting** | Middleware simples (contador em memória). |

### 3.3 Observabilidade

| Item                  | Decisão                                       |
| :-------------------- | :-------------------------------------------- |
| **Logs**              | `log/slog` com JSON estruturado.              |


### 3.4 Deploy e Infraestrutura

| Item           | Decisão                       |
| :------------- | :---------------------------- |
| **Hospedagem** | Fly.io ou Render (a definir). |
| **Banco**      | Neon.                         |


### 3.5 Documentação

| Item | Decisão |
| :--- | :--- |
| **README** | Como rodar o projeto. |
| **ARQUITETURA.md** | Este documento. |

### 3.6 Convenções

| Item | Decisão |
| :--- | :--- |
| **Formatação** | `gofmt` (nativo). |
| **Commits** | Conventional Commits (`feat:`, `fix:`, `docs:`, etc). |
| **Linting** | `golangci-lint` (futuro). |

---

## 4. Arquitetura

### 4.1 Padrão HSR (Handler → Service → Repository)

```text
Cliente HTTP
    ↓
Router (chi)
    ↓
Handler (recebe DTO, valida, chama Service)
    ↓
Service (lógica de negócio, lança erros personalizados)
    ↓
Repository (acesso ao banco, SQL puro)
    ↓
SQLite
```

### 4.2 Fluxo de uma Requisição (Ex: Vender Produto)

1. **Router:** Recebe `POST /produtos/{id}/vender`.
2. **Handler:** Extrai o ID e a quantidade do body. Chama `MovimentacaoService.Vender(id, qnt)`.
3. **Service:** Busca o produto no Repository. Valida se há estoque suficiente. Se não, lança `ErrEstoqueInsuficiente`. Se sim, atualiza a quantidade.
4. **Repository:** Executa `UPDATE produtos SET quantidade = ? WHERE id = ?`.
5. **Handler:** Recebe o resultado. Se erro, traduz para HTTP (400, 404, etc). Se sucesso, devolve JSON.
6. **Cliente:** Recebe a resposta.

## 5. Stack

| Camada             | Tecnologia                                  | Justificativa                                 |
| :----------------- | :------------------------------------------ | :-------------------------------------------- |
| **Linguagem**      | Go 1.22+                                    | Tipada, compilada, baixo consumo de RAM.      |
| **Framework web**  | `chi`                                       | Leve, compatível com `net/http`, não opinado. |
| **Banco de dados** | Postgres                                    | Já usei, é escalavel e simples                |
| **Driver**         | PGX                                         | Ja usei, simples                              |
| **Migrations**     | SQL puro em arquivos organizados            | Simples, versionável, sem dependência.        |
| **Logs**           | `log/slog` (nativo)                         | Estruturado, sem dependência externa.         |
| **Testes**         | `testing` (nativo)                          | Nativo, sem dependência.                      |
| **Validação**      | Manual + `go-playground/validator` (futuro) | Controle total no V1.                         |
| **Hash de senha**  | `golang.org/x/crypto/bcrypt`                | Maduro, simples, seguro.                      |
| **Sessão**         | Cookie com store em memória ou SQLite       | Simples, nativo do navegador.                 |

---

## 6. Estrutura de Pastas

```text
gestor-estoque/
├── cmd/
│   └── api/
│       └── main.go              # Ponto de entrada
├── internal/
│   ├── db/
│	│	migrations/
│   │     └── 001_init.sql
│   ├── estoque/
│      ├── types.go         # Struct Produto, dtos, erros
│      ├── handler.go       # Handler HTTP
│      ├── service.go       # Lógica de negócio
│      └── repository.go    # Acesso ao banco
├── .env.example
├── .gitignore
├── config/
├── .env
├── go.mod
├── go.sum
├──docs/
├── └── README.md
└──  └── ARQUITETURA.md
```

---

## 7. Erros Personalizados

| Erro                      | Campos                                             | HTTP Status | Quando ocorre                           |
| :------------------------ | :------------------------------------------------- | :---------- | :-------------------------------------- |
| `ErrEstoqueNaoEncontrado` | `ID string`                                        | 404         | Buscar/vender/atualizar ID inexistente. |
| `ErrEstoqueInsuficiente`  | `ProdutoID string`, `Pedido int`, `Disponivel int` | 400         | Vender mais do que tem.                 |
| `ErrEstoqueJaExiste`      | `ProdutoID string`                                 | 409         | Criar produto com ID duplicado.         |
| `ErrQuantidadeInvalida`   | `Qnt int`                                          | 400         | Passar quantidade negativa ou zero.     |
| `ErrPrecoInvalido`        | `Preco int`                                        | 400         | Passar preço negativo.                  |

**Como funciona:**
- O Repository lança erros de infraestrutura (não encontrado).
- O Service lança erros de negócio (estoque insuficiente).
- O Handler usa `errors.As` para traduzir erro → HTTP status.
- Erros desconhecidos viram HTTP 500.

---

## 8. Decisões Técnicas

| Decisão                              | Justificativa                                                        |
| :----------------------------------- | :------------------------------------------------------------------- |
| **Go em vez de Python/TS**           | Tipagem forte, compilado, baixo consumo de RAM, concorrência nativa. |
| **`chi` em vez de `gin`**            | Leve, compatível com `net/http`, não opinado.                        |
| **`.env` em vez de YAML**            | Simples, conhecido                                                   |
| **`slog` em vez de `zerolog`/`zap`** | Nativo, estruturado, resolve 90% dos casos.                          |
| **bcrypt em vez de argon2**          | Maduro, simples, seguro para 1 usuário.                              |
| **Sessão com cookie em vez de JWT**  | Simples, nativo do navegador, sem refresh token.                     |
| **HSR em vez de MVC**                | Achei mais facil de se entender                                      |
| **Estrutura por módulo**             | Isola features, fácil de navegar, evita espaguete.                   |

---

**Fim do documento.**

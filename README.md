

# 🎰 Projeto Loteria - Resultados e Notificações

Este projeto tem como objetivo consultar os resultados de loterias da Caixa Econômica Federal, salvar os dados em um banco de dados SQLite, gerar PDFs com os resultados e enviar notificações via Discord.

---

## 🚀 Funcionalidades

- **Consulta de resultados**: Obtém os resultados de loterias em tempo real através da API da Caixa.
- **Banco de dados**: Salva os resultados em um banco de dados SQLite para consultas futuras.
- **Geração de PDFs**: Cria um PDF com os resultados da loteria.
- **Notificações no Discord**: Envia os resultados para um canal do Discord usando webhooks, tanto em formato de embed quanto em PDF.

---

## 📋 Pré-requisitos

Antes de começar, você precisará ter instalado:

- Go 1.16 ou superior
- Git (opcional, para clonar o repositório)

---

## 🛠️ Configuração

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/seu-repositorio.git
cd seu-repositorio
```

### 2. Crie um arquivo `.env`

Crie um arquivo `.env` na raiz do projeto e adicione as seguintes variáveis:

```plaintext
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/SEU_WEBHOOK_ID/SEU_WEBHOOK_TOKEN
```

Substitua `SEU_WEBHOOK_ID` e `SEU_WEBHOOK_TOKEN` pelo ID e token do seu webhook do Discord.

### 3. Instale as dependências

Execute o seguinte comando para instalar as dependências necessárias:

```bash
go mod tidy
```

### 4. Execute o projeto

```bash
go run main.go
```

---

## 🗂️ Estrutura do Projeto

```
.
├── main.go                  # Script principal
├── go.mod                   # Arquivo de dependências do Go
├── README.md                # Documentação do projeto
├── .env                     # Variáveis de ambiente
├── bancoloteria.sqlite3     # Banco de dados SQLite
└── resultado_loteria.pdf    # Exemplo de PDF gerado
```

---

## 🛑 Como Parar o Script

O script é executado em um loop infinito. Para interrompê-lo, pressione `Ctrl + C` no terminal.

---

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## 🤝 Contribuição

Contribuições são bem-vindas! Siga os passos abaixo:

1. Faça um fork do projeto.
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`).
3. Commit suas alterações (`git commit -m 'Adicionando nova feature'`).
4. Faça um push para a branch (`git push origin feature/nova-feature`).
5. Abra um Pull Request.

---

## 📧 Contato

Se tiver dúvidas ou sugestões, entre em contato:

- **Nome**: Gilberto Jr
- **E-mail**: gilberto@infinitytec.info


---

Feito com ❤️ por Gilberto Jr 👋

---

### Explicação das Seções:

1. **Título e Descrição**: Apresenta o projeto de forma clara e direta.
2. **Funcionalidades**: Lista as principais funcionalidades do projeto.
3. **Pré-requisitos**: Informa o que é necessário para rodar o projeto.
4. **Configuração**: Passo a passo para configurar e executar o projeto.
5. **Estrutura do Projeto**: Mostra a organização dos arquivos.
6. **Como Parar o Script**: Instruções para interromper a execução.
7. **Licença**: Informa sobre a licença do projeto.
8. **Contribuição**: Explica como contribuir para o projeto.
9. **Contato**: Fornece informações para contato.

---

### Como Usar

1. **Clone o repositório**:
   ```bash
   git clone https://github.com/seu-usuario/seu-repositorio.git
   cd seu-repositorio
   ```

2. **Configure o arquivo `.env`**:
   Adicione o URL do webhook do Discord no arquivo `.env`.

3. **Instale as dependências**:
   ```bash
   go mod tidy
   ```

4. **Execute o projeto**:
   ```bash
   go run main.go
   ```


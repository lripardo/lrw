# LibRipardoWeb

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=lripardo_lrw&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=lripardo_lrw)

## Sobre

O objetivo deste projeto é manter uma base de código única. Quase todos os projetos modernos que possuem uma API,
necessitam de alguns serviços em comum como o de autenticação, por exemplo. O projeto está baseado em alguns conceitos
como DDD, 12 Factors e outros.

## Execução

A execução pode ser realizada das seguintes formas:

- Docker
- Diretamente na máquina de desenvolvimento

Eu recomendo rodar diretamente na máquina de desenvolvimento, pois o debugger fica muito mais fácil de se utilizar do
que via Docker (debug via TCP). O servidor irá buscar as configurações diretamente das variáveis de ambiente com os
valores exemplificados em [.env.example](.env.example). Todas as configurações padrões permitem que o servidor inicie
com as funcionalidades internas (banco de dados, configurações, envio de e-mails, etc). Esses serviços podem ser
alterados conforme a necessidade. Por exemplo, o banco de dados e o cache por padrão são implementados em memória, ou
seja, nenhuma persistência em disco é utilizada. Para alterar este comportamento, basta alterar a variável de ambiente
GORM_DRIVER_TYPE para mysql, informando assim, que o servidor utilize o MySQL como a unidade de persistência. Claro que
após alterar este comportamento, outras variáveis também deverão ser declaradas por questões técnicas (Host, porta,
usuário, senha, quantidade de conexões simultâneas, etc). O mesmo comportamento das configurações servem para os outros
tipos de serviço:

- E-mail (Print no console, AWS SES)
- Cache (Redis, memória local)
- Configurações (Variáveis de ambiente, memória local)
- Validações de entrada de usuário (Hcaptcha, Recaptcha, Default)
- Servidor (Gin Gonic)

Para executar o projeto na máquina de desenvolvimento utilize as ferramentas padrões do go (build, run, test, etc), ou
configure uma execução pela própria IDE (usuários da Goland approves ;D).

Para executar o projeto com Docker, utilize o [Docker Compose](https://docs.docker.com/compose/install). Após instalado,
execute o seguinte comando:

```
docker-compose --env-file .env -f docker-compose.yaml up -d
```

Onde .env é o arquivo que contém toda a sua configuração customizada dos serviços.

O servidor irá iniciar por padrão em <http://localhost:8080>.

## Principais ferramentas e tecnologias

- Docker
- Docker Compose
- Go
- Go Modules
- Redis
- MySQL
- HCaptcha
- AWS
- Outras dependências em [go.mod](go.mod)

## Ferramentas de auxílio no desenvolvimento

- PhpMyAdmin: Esta ferramenta é útil durante a fase de desenvolvimento. Com poucas linhas, todo o ambiente sobe como um
  container utilizando as credenciais das variáveis de ambiente. Disponível por padrão em: <http://localhost:8081>
- PhpRedisAdmin: Esta ferramenta é similar ao PhpMyAdmin, porém desenvolvida para o Redis. Disponível por padrão
  em: <http://localhost:8082>

Para utilizar as ferramentas de auxílio:

```
docker-compose -f docker-compose-tools.yaml up -d
```

## Boas práticas e observações

- Versionar apenas os arquivos necessários (.gitignore)
- Dados de ambiente em variáveis de ambiente
- Assinatura de commits com GPG
- Commit Semântico
- Algumas ferramentas que também são úteis, porém não utilizadas aqui:
    - [Dockerize](https://github.com/jwilder/dockerize): Fazer o container esperar até que algum serviço especificado
      esteja disponível a nível de conexão TCP
- Testes dentro de diferentes pacotes com _test ao final

## Referências

- DDD: <https://ddd-crew.github.io>
- REST: <https://martinfowler.com/articles/richardsonMaturityModel.html>
- Padrões de projeto em Go: <https://refactoring.guru/pt-br/design-patterns/go>
- Twelve-factor app: <https://12factor.net>
- Headers de segurança:
    - <https://docs.spring.io/spring-security/site/docs/5.0.x/reference/html/headers.html>
    - <https://cheatsheetseries.owasp.org/cheatsheets/REST_Security_Cheat_Sheet.html#security-headers>
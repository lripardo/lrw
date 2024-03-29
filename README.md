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

Quem irá trabalhar em features do backend, executar o projeto diretamente na máquina de desenvolvimento é a melhor
opção, pois o debugger fica muito mais fácil de se utilizar. Caso seja um desenvolvedor frontend e deseja apenas
consumir a API a execução via Docker é a melhor opção. As configurações serão buscadas diretamente das variáveis de
ambiente com os valores exemplificados em [.env.example](.env.example). As configurações padrões permitem que os
serviços iniciem com as funcionalidades internas em memória (banco de dados, cache, envio de e-mails, etc). Esses
serviços podem ser alterados conforme a necessidade. Por exemplo, para alterar o serviço de banco de dados para MySQL,
basta alterar a variável de ambiente GORM_DRIVER_TYPE para mysql. Claro que após alterar este comportamento, outras
variáveis também deverão ser declaradas por questões técnicas (host, porta, usuário, senha, quantidade de conexões
simultâneas, etc). Outros tipos de serviço:

- E-mail (Print no console, AWS SES)
- Cache (Redis, memória local)
- Configurações (Variáveis de ambiente, memória local)
- Validações de entrada de usuário (Hcaptcha, Default)
- API (Gin Gonic)

Copie o arquivo [.env.example](.env.example) para .env e deixe somente as configurações adequadas ao seu ambiente.

```
cp .env.example .env
vi .env
```

Para executar o projeto na máquina de desenvolvimento utilize as ferramentas padrões do go (build, run, test, etc) ou
configure uma execução pela própria IDE. Algumas IDE's como a Goland © necessitam do plugin .env para carregar variáveis
de ambiente diretamente de arquivos.

Para executar o projeto com Docker, utilize o [Docker Compose](https://docs.docker.com/compose/install). Após instalado,
execute o seguinte comando:

```
docker-compose up -d
```

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

As ferramentas de auxílio trazem de forma rápida as implementações de banco de dados (MySQL) e cache (Redis) via docker.
Também são configuradas softwares web de consulta de dados. Que são:

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

## Produção

Execute o script [build_docker_image.sh](build_docker_image.sh) para gerar uma imagem de produção docker. Sinta-se livre
para renomear o nome da imagem. Alguns comandos de limpeza que também estão no script são opcionais. A idéia é que seja
gerada a versão de build do executável e esteja dentro de uma imagem docker "from scratch".

## Referências

- DDD: <https://ddd-crew.github.io>
- REST: <https://martinfowler.com/articles/richardsonMaturityModel.html>
- Padrões de projeto em Go: <https://refactoring.guru/pt-br/design-patterns/go>
- Twelve-factor app: <https://12factor.net>
- Headers de segurança:
    - <https://docs.spring.io/spring-security/site/docs/5.0.x/reference/html/headers.html>
    - <https://cheatsheetseries.owasp.org/cheatsheets/REST_Security_Cheat_Sheet.html#security-headers>
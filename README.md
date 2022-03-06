# LibRipardoWeb

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=lripardo_lrw&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=lripardo_lrw)


## Sobre

O objetivo deste projeto é manter uma base de código única. Quase todos os projetos modernos que possuem uma API,
necessitam de alguns serviços em comum como o de autenticação, por exemplo.

## Execução

Este projeto utiliza [Docker Compose](https://docs.docker.com/compose/install). Após instalado, execute o seguinte
comando:

```
docker-compose --env-file .env -f docker-compose.yaml up -d
```

O servidor irá iniciar por padrão em <http://localhost:8080>.

## Principais ferramentas e tecnologias

- Docker
- Docker Compose
- Go
- Go Modules
- Redis
- Outras dependências em [go.mod](go.mod)

## Ferramentas de auxílio no desenvolvimento

- PhpMyAdmin: Esta ferramenta é útil durante a fase de desenvolvimento. Com poucas linhas, todo o ambiente sobe como um
  container utilizando as credenciais das variáveis de ambiente. Disponível por padrão em: <http://localhost:8081>
- PhpRedisAdmin: Esta ferramenta é similar ao PhpMyAdmin, porém desenvolvida para o Redis. Disponível por padrão
  em: <http://localhost:8082>

Para utilizar as ferramentas de auxílio:

```
docker-compose -f docker-compose.tools.yaml up -d
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
- Headers de segurança: <https://docs.spring.io/spring-security/site/docs/5.0.x/reference/html/headers.html>
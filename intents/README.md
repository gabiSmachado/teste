# UMA SOLUÇÃO PARA REDES AUTO-ORGANIZADAS BASEADAS EM INTENÇÕES
Este trabalho tem o objetivo de criar redes 5g auto-organização da rede através de Intents, 
desenvolvido em cima da arquitetura da rede disponível pela O-RAN Alliance em https://docs.o-ran-sc.org/en/latest/architecture/architecture.html.

## Arquitetura

O Intent é definido pelo operador, através do intent interface que lê as solicitações do usuário, que são encaminhas para o Intent Receiver.

O Intent Receiver executa as operações de criação, listagem e visualização da descrição de intents, assim como deletá-los, através de uma interface REST com dados JSON. Quando um intente é criado o Intent Receiver verifica a sanidade e, se sã, é salvo em uma DataBase e publicado no Message Queue para consumo do Intent Broker.

O Intent Broker consome as mensagens e as traduz em em ações de controle que devem ser tomadas pelo near-RT RIC, publicando-as na interface A1 como políticas, que serão distribuídas para as xApps do near-RT RIC através do A1 Mediator.

![Alt text](/arqu.png)

## Requisitos 
Linguagem GO;

Mariadb: https://www.digitalocean.com/community/tutorials/how-to-install-mariadb-on-ubuntu-20-04-pt;

Apache Kafka: https://kafka.apache.org/quickstart;

Near-RT RIC: https://github.com/gabiSmachado/Near-RT-RIC_deploy;

## Execução

Após a instalação de todo os requisitos, inicialize o Mariadb crie a tabela para salvar os intents:

![Alt text](/imgs/createTable.png)


### Inicialize o Kafka: 
Responsável pela transferência das mensagens.
```
bin/zookeeper-server-start.sh config/zookeeper.properties
bin/kafka-server-start.sh config/server.properties
```
### Execute o Intent Receiver: 

O serviço ficara rodando no SMO, para realizar as operações com os Intents.

![Alt text](/imgs/receiver.png)

## Trabalhando com os Intents:

### Criar Intent: 
Na criação de um intent, deve ser informado alguns campos obrigatórios que serão solicitados pelo serviço.
Após informá-los a mensagem que o intent foi criado e seu id devem ser retornadas.

![Alt text](/imgs/createIntent.png)

### Listar Intents: 
É possível visualizar todos os intentes salvos até o momento com o comando:

![Alt text](/imgs/listIntent.png)

### Descrição de um Intent: 
Para visualizar as propriedade de um Intent, informe seu ID.

![Alt text](/imgs/show.png)

### Deletar um Intent:
É possível deletar um Intente informando seu ID.

![Alt text](/imgs/delete.png)

### Intent broker:
Em desenvolvimento.

#### É possível visualizar todas as operações realizadas através do log mantido pelo Intent Receiver:

![Alt text](/imgs/log.png)

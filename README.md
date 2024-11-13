# Trabajo de microservicios

```
Author: Emmanuel Guerreiro
Legajo: 47262
Año: 2024
```

## Descripción

Se implementa el microservicio de gestion de stocks basado en las definiciones de la catedra

## Requisitos

- go 1.22.0
- rabbitmq
- mongo
- docker

## Como correr

_Docker_

```
docker-compose up --build
```

_Local_

```
make run
```

Interfaz RABBITMQ

```
{
"correlation_id":"123123",
"routing_key":"Remote RoutingKey to Reply",
"exchange":"order-placed",
  "message": {
    "cartId": 123123,
    "userId": "testing",
    "articles": [
      {
        "articleId": 123,
        "quantity": 10
      }
    ]
  }
}
```

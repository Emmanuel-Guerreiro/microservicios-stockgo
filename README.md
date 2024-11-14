# Trabajo de microservicios

```
Author: Emmanuel Guerreiro
Legajo: 47262
Año: 2024
```

## Descripción

Se implementa el microservicio de gestion de stocks basado en las definiciones de la catedra

Este permite saber el stock de un producto y poder incrementar y decrementar el stock.

Implementa los patrones Event Sourcing y CQRS para llevar registro de los cambios en el stock y generar una vista de stock mas eficiente para la consulta.

## Requisitos

- go 1.22.0
- rabbitmq
- docker

## Como correrlo

_Docker_

```
docker-compose up --build
```

_Local_

Para correrlo localmente es requerido tener instalados:

- RabbitMQ
- MongoDB

```
make run
```

## Casos de uso

### CU: Consultar stock

- Ante un evento de consulta de stock recibido a través de un mensaje asincrono se devuelve, el stock actual del producto de manera asincrona.

#### Interfaz RABBITMQ

##### Escucha

- **Exchange:** stock_consulting
- **Tipo:** direct
- **Routing Key:** get_stock

###### Mensaje

```
{
  "articleId": string
}
```

_ejemplo:_

```
{
  "articleId": "123"
}
```

##### Escucha

- **Exchange:** stock_consulting
- **Tipo:** direct
- **Routing Key:** stock_response

###### Mensaje

```
{
  "articleId": string,
  "stock": int
}
```

_ejemplo:_

```
{
  "articleId": "123",
  "stock": 10
}
```

### CU: Configurar stock mínimo producto

- Se configura el stock mínimo para un producto a través de la interfaz REST. Definiendo a partir de que stock hay que emitir un evento de reposición a la cola de eventos asincrona.

#### Interfaz REST

##### POST

- **URL:** /article-config

###### Body

```
{
  "articleId": string,
  "alertMinQuantity": int
}
```

_ejemplo:_

```
{
  "articleId": "123",
  "alertMinQuantity": 10
}
```

### CU: reposición de stock

- Ante un evento de reposición de stock se almacena el evento de reposición y se re-calcula el stock del producto en la vista. Se realiza de manera asincrona.

#### Interfaz RABBITMQ

##### Escucha

- **Exchange:** stock_reposition
- **Tipo:** direct
- **Routing Key:** stock_reposition

###### Mensaje

```
 {
 	"message":{
		"articleId": "string",
		"amount": int
	}
}
```

_ejemplo:_

```
 {
 	"message":{
		"articleId": "123",
		"amount": 10
	}
}
```

### CU: disminución de stock

- Ante un evento de order_placed (Emitido actualmente por el microservicio de orders) se persiste un evento de decremento de stock y se calcula el stock del producto en la vista. Se realiza de manera asincrona.\
  Ante stock no disponible para completar la orden se emite un evento de not_enough_stock a través de la cola de eventos asincrona

#### Interfaz RABBITMQ

##### Escucha

- **Exchange:** order_placed
- **Tipo:** fanout

###### Mensaje

```
{
  "correlation_id": string,
  "message":{
    "cartId": string,
    "userId": string,
    "articles":[
      {
        "articleId": string,
        "quantity": int
      }
    ]
  }
}

```

_ejemplo:_

```
{
  "correlation_id":"123123",
  "message":{
    "cartId":123123,
    "userId":"testing",
    "articles":[
      {
        "articleId":123,
        "quantity":10
      }
    ]
  }
}

```

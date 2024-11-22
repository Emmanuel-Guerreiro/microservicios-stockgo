# Trabajo de microservicios

```
Author: Emmanuel Guerreiro
Legajo: 47262
Año: 2024
```

Ordenar:
Casos de uso
Interfaz

EXPLICAR MEJOR LOS CU
Aclarar el CU de re-calculo de la proyeccion de stock -> Cu interno, derivado

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

## Modelo de datos

_Evento_

```
{
  ID: string
  Type: string
  DecrementEvent: DecrementEvent
  RepositionEvent: RepositionEvent
  Created: date
  EventStatus: date
}

```

DecrementEvent

```
{
	ArticleId: string
	Quantity : int
}
```

RepositionEvent

```
{
  ArticleId: string
  Quantity: int
}
```

_Article config_

```
{
  articleId: string
  alertMinQuantity: int
  createdAt: date
  updatedAt: date
}
```

**Proyeccion stock**

```
{
  ArticleId: string
  Stock: string
  Created: date
  Updated: date
}
```

## Casos de uso

### CU: Cálculo de stock de un producto

- Se calcula a partir del historial de eventos para un producto
- No puede ser negativo

### CU: Re-compra de stock

- Si el stock cae por debajo del umbral configurado, se solicita una re-compra para un producto
- Se envia el evento, junto a un id para poder trackear la orden
- Si un producto no tiene eventos de stock, se define en 0.

### CU: Consultar stock

- Se devuelve el stock actual de un producto, a partir de su articleId
- Si no se encuentra el stock en la proyección, se ejecuta el `CU: Cálculo de stock de un producto`..

### CU: Configurar stock mínimo producto

- Un usuario registrado en el sistema puede configurar el stock mínimo para un producto.
- El stock mínimo no puede ser negativo
- Se almacenan la fecha de modificación

### CU: reposición de stock

- Desde otro servicio se puede generar un ingreso de stock para un producto
- La cantidad debe ser no negativa
- Se almacena el nuevo stock junto con la fecha de ingreso
- Se almacena el valor absoluto del movimiento de stock
- Ejecuta el `CU: Cálculo de stock de un producto`.

### CU: disminución de stock

- Otro servicio puede generar una disminución de stock para un producto
- La cantidad debe ser no negativa
- Se almacena el nuevo stock junto con la fecha de ingreso
- Se almacena el valor absoluto del movimiento de stock
- Ejecuta el `CU: Cálculo de stock de un producto`.

## Interfaz RABBITMQ

### Escucha - Exchange: stock_consulting - Routing Key: get_stock

- **Exchange:** stock_consulting
- **Tipo:** direct
- **Routing Key:** get_stock

##### Mensaje

```

{
  "articleId": string
}

```

_Ejemplo:_

```

{
  "articleId": "123"
}

```

### Emit - Exchange: stock_consulting - Routing Key: stock_response

- **Exchange:** stock_consulting
- **Tipo:** direct
- **Routing Key:** stock_response

#### Mensaje

```

{
  "articleId": string,
  "stock": int,
  "createdAt": string,
  "updatedAt": string
}

```

_Ejemplo:_

```

{
  "articleId": "123",
  "stock": 10,
  "createdAt": "2023-01-01T00:00:00Z",
  "updatedAt": "2023-01-01T00:00:00Z"
}

```

### Escucha - Exchange: stock_reposition - Routing Key: stock_reposition

- **Exchange:** stock_reposition
- **Tipo:** direct
- **Routing Key:** stock_reposition

#### Mensaje

```

{
  correlation_id: string,
  message:{
    articleId: "string",
    amount: int
  }
}

```

_Ejemplo:_

```

{
  "message":{
    "articleId": "123",
    "amount": 10
  }
}

```

### Escucha - Exchange: order_placed

- **Exchange:** order_placed
- **Tipo:** fanout

#### Mensaje

```

{
  "correlation_id": string,
  "message": {
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

_Ejemplo:_

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

## Interfaz REST

### POST - URL: /article-config

- **URL:** /article-config
- **Method:** POST
- **Header:**:
  - **Content-Type:** application/json
  - **Authorization**: Bearer <token>

#### Body

```

{
  "articleId": string,
  "alertMinQuantity": int
}

```

_Ejemplo:_

```

{
  "articleId": "123",
  "alertMinQuantity": 10
}

```

#### Respuesta:

##### Creado correctamente

- _Status Code_: 201 - Created
- _Body_:

```
{
  "articleId": string,
  "alertMinQuantity": int,
  "createdAt": string,
  "updatedAt": string
}
```

#### Errores

- _Status Code_: 400 - Bad Request

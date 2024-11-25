Mensajes utilizados para probar en la interfaz de RabbitMQ

GET STOCK - >
rk: get_stock
{
"articleId": "123"
}

ORDER PLACED ->
{
"correlation_id": "123123",
"message": {
"orderId": "123123",
"cartId": "123123",
"articles": [
{
"articleId": "1211",
"quantity": 10
}
]
}
}

stock reposition ->
rk: stock_reposition

{
"message":{
"articleId": "1211",
"amount": 10
}
}

# Фреймворк для коммуникации между сервером и клиентом

Позволяет вызывать удалённые процедуры на стороне сервера и принимать PUSH уведомления
Механизм запрос/ответ работает поверх протокола HTTP 1.1
PUSH уведомления получаются по Websocket

# TODO генератор клиента и интерфейса серверных методов

struct PushMessage1 {
    name:string
}

struct PushMessage2 {
    login:string
    email:string
}

func Sum(n1:int, n2:int):int
push MyChannel(PushMessage1, PushMessage2)
## Телеграм-бот для пересылки email-сообщений в чаты

## Настройка
.env:
- TG_TOKEN - токен бота
- ADMIN_ID - Id пользователя-администратора (отправляет команды боту)
- DEBUG - Debug режим (логи ТГ-бота)
- STORAGE_PREFIX - путь для хранения email-ов в html 
- HTTP_PREFIX - префикс сервера в адресе по которому будут доступны сообщения ( HTTP_PREFIX + STORAGE_PREFIX + UUID ) - адрес html сообщения
- HTTP_ADDR - адрес на котором запускается сервер (host:port)

## Деплой со сборкой
>```docker build --tag '<image-name>' .```

>```docker run --env-file ./.env -v ./storage:/app/storage <image-name>```

## Деплой с образом из Dockerhub
>```docker run --env-file ./.env -v ./storage:/app/storage vladdoth/mail2tg:latest```


## Управление
Настройка IMAP:
>```/set {"imapServer": "imap.yandex.ru:993", "imapUser": "vladDotH", "imapToken": "1234567890`"}```
 
Установка правила:
>```/set {"name": "default", "box": "INBOX", "chatId": 841733382}```

Опционально delay, originalMsgId 

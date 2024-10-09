## Телеграм-бот для пересылки email-сообщений в чаты

## Настройка
.env:
- TG_TOKEN - токен бота
- ADMIN_ID - Id пользователя-администратора (отправляет команды боту)
- DEBUG - Debug режим (логи ТГ-бота)

## Деплой
>```docker build --tag '<image-name>' .```

>```docker run --env-file ./.env -v ./storage:/app/storage <image-name>```

## Управление
Настройка IMAP:
>```/set {"imapServer": "imap.yandex.ru:993", "imapUser": "vladDotH", "imapToken": "1234567890`"}```
 
Установка правила:
>```/set {"name": "default", "box": "INBOX", "chatId": 841733382}```

Опционально delay, originalMsgId 
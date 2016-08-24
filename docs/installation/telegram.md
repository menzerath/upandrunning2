# Telegram-Bot / -Notifications
Receiving Telegram-notifications using the built-in Telegram-bot is quite easy but needs a bit of configuration:

* Create a new Telegram-bot using the [Botfather](https://telegram.me/BotFather).
* Insert the received API-Token into your configuration-file as `telegramBotApiKey`.
* Start a conversation with the newly created bot to receive your personal user-id.
* Set this user-id as notification-target for the selected websites.
* Done!

## Available commands
You can use the following commands when chatting with your Telegram-bot:

| Command   | Description                                                     |
|-----------|-----------------------------------------------------------------|
| `/start`  | Start a conversation and allow the bot to send messages to you. |
| `/id`     | Receive your user-id.                                           |
| `/server` | Get some information about the bot's server.                    |
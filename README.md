# TelegoBot

TelegoBot is a Telegram bot designed to deliver real-time news articles and summaries to users based on their interests. It provides a convenient way for users to stay updated with the latest news on various topics without leaving the Telegram app.


## Features

- Real-time News Delivery
- Summarization
- Source Management

Getting Started

To run TelegoBot locally or deploy it to a server, follow these steps:

```bash
git clone https://github.com/yourusername/TelegoBot.git
```
Install the required dependencies:

```bash
# Navigate to the project directory
cd TelegoBot
# Install dependencies (example using Go modules)
go mod tidy
```
Configure the bot and database settings by creating a .env file in the project root and filling in the necessary environment variables. You can use the provided .env.example as a template.
Build and run the application:

```bash
go build -o TelegoBot
./TelegoBot
```
## Roadmap

- Add docker

- Add tests for primary func's

- Add Chat Mod

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)



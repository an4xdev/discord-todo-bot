# Simple todo Discord bot

Create simple todo list in channel, mark tasks as done, and clear tasks from channel.

## Why?

I wanted to create a simple bot and practice Github Actions and deploying to VPS.

## How to run

1. Create a Discord bot and get the token.
2. Rename `.env.example` in the root directory to `.env` and add your bot token and bot id.
3. Install the required packages:
    ```
    npm install
    ```
4. Register commands in Discord:
   ```
   npm run register-commands
   ```
   This will register the commands defined in `commands.js` with your bot.
5. Run the bot:
   ```
   npm start
   ```
   This will start the bot and it will begin listening for commands in your Discord server.
6. Invite the bot to your server using the OAuth2 URL generated in the Discord Developer Portal.
7. Use the commands in your Discord server to manage your todo list.
## Commands
- `/todo add <task>`: Add a new task to the todo list.
- `/todo list`: List all tasks in the todo list.
- `/todo reset`: Reset the todo list, removing all tasks.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


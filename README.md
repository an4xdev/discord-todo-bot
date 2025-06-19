# Simple todo Discord bot

Create simple todo list in channel, mark tasks as done, and clear tasks from channel.

## Why?

I wanted to create a simple bot and practice Github Actions and deploying to VPS.

## How to run

1. **Create a Discord bot and get the token:**
   - Go to [Discord Developer Portal](https://discord.com/developers/applications)
   - Create a new application and bot
   - Copy the bot token and client ID

2. **Set up environment variables:**
   - Rename `.env.example` to `.env` in the root directory
   - Add your bot token and client ID:
     ```
     DISCORD_TOKEN=your_bot_token_here
     CLIENT_ID=your_client_id_here
     ```

3. **Install Go dependencies:**
   ```bash
   go mod download
   ```

4. **Build the application:**
   ```bash
   go build -o discord-bot main.go
   ```

5. **Run the bot:**
   ```bash
   ./discord-bot
   ```
   The bot will automatically register commands on startup and begin listening for interactions in your Discord server.

6. **Invite the bot to your server:**
   - Use the OAuth2 URL generator in the Discord Developer Portal
   - Select "bot" and "applications.commands" scopes
   - Select necessary permissions (Send Messages, Use Slash Commands, etc.)

7. **Use the commands in your Discord server:**
   - `/todo add <task>` - Add a new todo item
   - `/todo list` - Show all todos for the current channel
   - `/todo reset` - Remove all todos from the current channel


## Production Deployment

The bot includes automatic deployment via GitHub Actions with systemd service management. See the CI/CD workflow for automated testing and deployment.

## Database

The bot uses SQLite database (`todos.db`) which will be created automatically on first run. The database stores:
- Todo items with completion status
- Channel-specific todo lists
- Message tracking for cleanup operations
## Commands
- `/todo add <task>`: Add a new task to the todo list.
- `/todo list`: List all tasks in the todo list.
- `/todo reset`: Reset the todo list, removing all tasks.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


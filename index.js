require("dotenv").config();
const {
    Client,
    GatewayIntentBits,
    Collection,
    EmbedBuilder,
} = require("discord.js");
const db = require("./database/db");

const client = new Client({
    intents: [GatewayIntentBits.Guilds],
});

client.commands = new Collection();

const todoCommand = require("./commands/todo");
client.commands.set(todoCommand.data.name, todoCommand);

client.once("ready", () => {
    console.log(`Bot ${client.user.tag} jest online!`);
    console.log(`Bot jest na ${client.guilds.cache.size} serwerach`);
});

client.on("debug", (info) => {
    if (info.includes("heartbeat")) return;
    console.log("Debug:", info);
});

client.on("error", (error) => {
    console.error("Błąd klienta:", error);
});

client.on("interactionCreate", async (interaction) => {
    if (interaction.isChatInputCommand()) {
        const command = client.commands.get(interaction.commandName);
        if (!command) return;

        try {
            await command.execute(interaction);
        } catch (error) {
            console.error(error);
            await interaction.reply({
                content: "Wystąpił błąd podczas wykonywania komendy!",
                ephemeral: true,
            });
        }
    }

    if (interaction.isButton()) {
        if (interaction.customId.startsWith("complete_todo_")) {
            const todoId = interaction.customId.split("_")[2];

            db.run(
                "UPDATE todos SET completed = 1 WHERE id = ?",
                [todoId],
                function (err) {
                    if (err) {
                        console.error(err);
                        return;
                    }

                    db.get(
                        "SELECT * FROM todos WHERE id = ?",
                        [todoId],
                        async (err, row) => {
                            if (err || !row) {
                                console.error(err);
                                return;
                            }

                            const completedEmbed = new EmbedBuilder()
                                .setTitle("✅ TODO - Wykonane")
                                .setDescription(`~~${row.content}~~`)
                                .setColor(0x90ee90)
                                .setTimestamp();

                            await interaction.update({
                                embeds: [completedEmbed],
                                components: [],
                            });
                        }
                    );
                }
            );
        }
    }
});

client.login(process.env.DISCORD_TOKEN);

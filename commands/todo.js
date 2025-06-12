const {
    SlashCommandBuilder,
    EmbedBuilder,
    ActionRowBuilder,
    ButtonBuilder,
    ButtonStyle,
} = require("discord.js");
const db = require("../database/db");

module.exports = {
    data: new SlashCommandBuilder()
        .setName("todo")
        .setDescription("Zarządzanie zadaniami TODO")
        .addSubcommand((subcommand) =>
            subcommand
                .setName("add")
                .setDescription("Dodaj nowe zadanie")
                .addStringOption((option) =>
                    option
                        .setName("task")
                        .setDescription("Treść zadania")
                        .setRequired(true)
                )
        )
        .addSubcommand((subcommand) =>
            subcommand
                .setName("list")
                .setDescription("Pokaż wszystkie zadania na tym kanale")
        )
        .addSubcommand((subcommand) =>
            subcommand
                .setName("reset")
                .setDescription("Usuń wszystkie zadania z tego kanału")
        ),

    async execute(interaction) {
        const subcommand = interaction.options.getSubcommand();

        if (subcommand === "add") {
            const task = interaction.options.getString("task");
            const channelId = interaction.channel.id;

            const todoEmbed = new EmbedBuilder()
                .setTitle("📝 Nowe TODO")
                .setDescription(task)
                .setColor(0x0099ff)
                .setTimestamp();

            const row = new ActionRowBuilder().addComponents(
                new ButtonBuilder()
                    .setCustomId(`complete_todo_temp`)
                    .setLabel("✅ Oznacz jako wykonane")
                    .setStyle(ButtonStyle.Success)
            );

            const message = await interaction.reply({
                embeds: [todoEmbed],
                components: [row],
                fetchReply: true,
            });

            db.run(
                "INSERT INTO todos (channel_id, message_id, content) VALUES (?, ?, ?)",
                [channelId, message.id, task],
                function (err) {
                    if (err) {
                        console.error(err);
                        return;
                    }

                    const updatedRow = new ActionRowBuilder().addComponents(
                        new ButtonBuilder()
                            .setCustomId(`complete_todo_${this.lastID}`)
                            .setLabel("✅ Oznacz jako wykonane")
                            .setStyle(ButtonStyle.Success)
                    );

                    message.edit({ components: [updatedRow] });
                }
            );
        } else if (subcommand === "list") {
            const channelId = interaction.channel.id;

            db.all(
                "SELECT * FROM todos WHERE channel_id = ? ORDER BY created_at DESC",
                [channelId],
                async (err, rows) => {
                    if (err) {
                        console.error(err);
                        return interaction.reply({
                            content: "Błąd podczas pobierania zadań!",
                            ephemeral: true,
                        });
                    }

                    if (rows.length === 0) {
                        return interaction.reply({
                            content: "Brak zadań na tym kanale!",
                            ephemeral: true,
                        });
                    }

                    const todoList = rows
                        .map((row) => {
                            const status = row.completed ? "✅" : "❌";
                            const text = row.completed
                                ? `~~${row.content}~~`
                                : row.content;
                            return `${status} ${text}`;
                        })
                        .join("\n");

                    const listEmbed = new EmbedBuilder()
                        .setTitle("📋 Lista TODO")
                        .setDescription(todoList)
                        .setColor(0x00ae86)
                        .setTimestamp();

                    try {
                        const listMessage = await interaction.reply({
                            embeds: [listEmbed],
                            fetchReply: true,
                        });

                        db.run(
                            "INSERT INTO list_messages (channel_id, message_id) VALUES (?, ?)",
                            [channelId, listMessage.id],
                            function (err) {
                                if (err) {
                                    console.error(err);
                                }
                            }
                        );
                    } catch (error) {
                        console.error("Błąd podczas wysyłania listy:", error);
                    }
                }
            );
        } else if (subcommand === "reset") {
            const channelId = interaction.channel.id;

            db.all(
                "SELECT message_id FROM todos WHERE channel_id = ?",
                [channelId],
                async (err, rows) => {
                    if (err) {
                        console.error(err);
                        return interaction.reply({
                            content: "Błąd podczas pobierania zadań!",
                            ephemeral: true,
                        });
                    }

                    let deletedMessages = 0;
                    for (const row of rows) {
                        try {
                            const message =
                                await interaction.channel.messages.fetch(
                                    row.message_id
                                );
                            await message.delete();
                            deletedMessages++;
                        } catch (error) {
                            console.log(
                                `Nie można usunąć wiadomości ${row.message_id}:`,
                                error.message
                            );
                        }
                    }

                    db.all(
                        `SELECT message_id FROM list_messages WHERE channel_id = ?`,
                        [channelId],
                        async (err, listRows) => {
                            if (err) {
                                console.error(err);
                                return interaction.reply({
                                    content:
                                        "Błąd podczas pobierania listy zadań!",
                                    ephemeral: true,
                                });
                            }

                            for (const listRow of listRows) {
                                try {
                                    const message =
                                        await interaction.channel.messages.fetch(
                                            listRow.message_id
                                        );
                                    await message.delete();
                                    deletedMessages++;
                                } catch (error) {
                                    console.log(
                                        `Nie można usunąć wiadomości listy ${listRow.message_id}:`,
                                        error.message
                                    );
                                }
                            }
                        }
                    );

                    db.run(
                        "DELETE FROM todos WHERE channel_id = ?",
                        [channelId],
                        function (err) {
                            if (err) {
                                console.error(err);
                                return interaction.reply({
                                    content: "Błąd podczas resetowania zadań!",
                                    ephemeral: true,
                                });
                            }

                            const resetEmbed = new EmbedBuilder()
                                .setTitle("🗑️ Reset TODO")
                                .setDescription(
                                    `Usunięto ${this.changes} zadań z bazy danych i ${deletedMessages} wiadomości z kanału.`
                                )
                                .setColor(0xff6b6b)
                                .setTimestamp();

                            interaction.reply({ embeds: [resetEmbed] });
                        }
                    );
                }
            );
        }
    },
};

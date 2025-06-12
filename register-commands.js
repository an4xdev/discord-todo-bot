require("dotenv").config();
const { REST, Routes } = require("discord.js");

const commands = [require("./commands/todo").data.toJSON()];

const rest = new REST({ version: "10" }).setToken(process.env.DISCORD_TOKEN);

(async () => {
    try {
        console.log("Rozpoczęcie rejestracji komend...");

        await rest.put(Routes.applicationCommands(process.env.CLIENT_ID), {
            body: commands,
        });

        console.log("Komendy zostały zarejestrowane!");
    } catch (error) {
        console.error(error);
    }
})();

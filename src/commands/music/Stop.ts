import { SlashCommandBuilder, type CommandInteraction } from 'discord.js';

import type { Bot } from '@/structures/Bot';
import { Command } from '@/structures/Command';

export class Stop extends Command {
  constructor(bot: Bot) {
    super(bot, {
      isActive: true,
      build: new SlashCommandBuilder()
        .setName('stop')
        .setDescription('Stops the music playback')
    });
  }

  async execute(interaction: CommandInteraction) {
    this.bot.musicPlaybackHandler.stop(true);
    await this.bot.messageChannelHandler.signCommandExecution(interaction);
  }
}
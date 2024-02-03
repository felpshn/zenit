package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/larahfelipe/saturn/internal/command"
	"github.com/larahfelipe/saturn/pkg/discord"
)

type Bot struct {
	Token   string
	Command *command.Command
	*discord.DiscordService
}

func New(token string, command *command.Command) (*Bot, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("token cannot be empty")
	}

	// ft := &Feature{
	// 	External: &External{
	// 		Youtube: &youtube.Client{},
	// 	},
	// }

	// mq := &music.MusicQueue{
	// 	PlaybackState: make(chan music.PlaybackState, 5),
	// 	Songs:         []music.Song{},
	// }

	ds, err := discord.NewService(token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Token:          token,
		Command:        command,
		DiscordService: ds,
	}

	bot.CommandMessageCreateHandler(command.Handle, command.Prefix)

	return bot, nil
}

func (bot *Bot) BuildErrorMessageEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "❌ Oops, a wild error appeared! 😱",
			IconURL: bot.Session.State.User.AvatarURL("256"),
		},
		Description: message,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Please try again later",
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Color:     0xFB3640,
	}
}

func (bot *Bot) BuildMessageEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    bot.Session.State.User.Username,
			IconURL: bot.Session.State.User.AvatarURL("256"),
		},
		Description: message,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "From space",
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Color:     0x6E76E5,
	}
}

func (bot *Bot) MakeVoiceConnection(m *command.Message) (*discordgo.VoiceConnection, error) {
	for _, guild := range bot.Session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				if vs.UserID != bot.Session.State.User.ID {
					bot.Session.ChannelMessageSendEmbed(m.ChannelID, bot.BuildMessageEmbed(fmt.Sprintf("Yay! Joining the party on <#%s>", vs.ChannelID)))
				}

				vc, err := bot.Session.ChannelVoiceJoin(guild.ID, vs.ChannelID, false, true)
				if err != nil {
					bot.Session.ChannelMessageSendEmbed(m.ChannelID, bot.BuildErrorMessageEmbed("It seems that I'm not in the mood for partying right now. Maybe later?"))
					return nil, err
				}

				return vc, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find a voice channel for the user who requested the song")
}

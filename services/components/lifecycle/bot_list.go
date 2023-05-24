package lifecycle

import (
	"github.com/forta-network/forta-node/config"
)

// FindUpdatedBots finds the updated bots in the first list by comparing with the second list.
func FindUpdatedBots(list1, list2 []config.AgentConfig) (resultList []config.AgentConfig) {
	for _, bot1 := range list1 {
		findBotAndDo(bot1, list2, func(foundBot config.AgentConfig) {
			if !bot1.Equal(foundBot) {
				resultList = append(resultList, foundBot)
			}
		})
	}
	return
}

// Drop drops a bot from the list.
func Drop(bot config.AgentConfig, botList []config.AgentConfig) (resultList []config.AgentConfig) {
	for _, currBot := range botList {
		if bot.ContainerName() == currBot.ContainerName() {
			continue
		}
		resultList = append(resultList, currBot)
	}
	return
}

func findBotAndDo(bot config.AgentConfig, botList []config.AgentConfig, do func(eachBot config.AgentConfig)) {
	for _, currBot := range botList {
		if bot.ContainerName() == currBot.ContainerName() {
			do(currBot)
			return
		}
	}
}

// FindExtraBots looks at the second list and finds the items that doesn't exist in the first list.
func FindExtraBots(original, extraIn []config.AgentConfig) (resultList []config.AgentConfig) {
	return FindMissingBots(extraIn, original)
}

// FindMissingBots looks at the first list and finds the items that doesn't exist in the second list.
func FindMissingBots(original, missingIn []config.AgentConfig) (resultList []config.AgentConfig) {
	for _, bot := range original {
		missBotAndDo(bot, missingIn, func(missingBot config.AgentConfig) {
			resultList = append(resultList, missingBot)
		})
	}
	return
}

func missBotAndDo(bot config.AgentConfig, botList []config.AgentConfig, do func(eachBot config.AgentConfig)) {
	for _, currBot := range botList {
		if bot.ContainerName() == currBot.ContainerName() {
			return
		}
	}
	do(bot)
}

// GetBotIDs makes a new slice of bot IDs.
func GetBotIDs(botList []config.AgentConfig) (ids []string) {
	for _, bot := range botList {
		ids = append(ids, bot.ID)
	}
	return
}

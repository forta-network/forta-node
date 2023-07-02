import axios from 'axios'
import TelegramBot from 'node-telegram-bot-api'
import { Controller, Param, Body, Get, Post, Put, Delete, QueryParam, BodyParam, JsonController } from 'routing-controllers'
import { BigNumber, utils } from 'ethers'
import { APIEmbed, EmbedBuilder } from 'discord.js'

// const botKey = https://api.telegram.org/bot6231325118:AAG_tGdf7Db7ndDuIDuX0MHGNgtAcACMz_g/sendMessage?chat_id=-858059392&text=something
const botKey = '6231325118:AAG_tGdf7Db7ndDuIDuX0MHGNgtAcACMz_g'
const chatId = '-858059392'
const explorerUrl = 'https://securityalliance.dev/'
const bot = new TelegramBot(botKey, { polling: false })
const discordWebhookId = '1125156706855436308/RHD91-NS-GFPwitOy2zo_vD-okJzs6u2w92gr0b5zu-4h_WJpWRORGdMwSMmugm1a6Uj'

interface FortaAlert {
    addresses: string[]
    alertId: string
    createdAt: string
    description: string
    metadata: {
        anchorPrice: string
        cTokenAddress: string
        protocolVersion: string
        reporterPrice: string
        underlyingTokenAddress: string
        validatorProxyAddress: string
    }
    name: string
    severity: string
    source: {
        block: {
            chainId: number
            hash: string
            number: string
            timestamp: string
        }
        bot: {
            id: string
        }
    }
    transactionHash: string
}

const formatPrice = (price: string) => {
    const priceBn = BigNumber.from(price)
    const remainder = priceBn.mod(1e5)
    const formatted = utils.formatUnits(priceBn.sub(remainder), 6)
    return formatted
}

function toEscapeMsg(str: string): string {
    return str.replace(/_/gi, '\\_').replace(/-/gi, '\\-').replace('~', '\\~').replace(/`/gi, '\\`').replace(/\./g, '\\.')
}

const formatTelegramAlert = (alert: FortaAlert) => {
    return toEscapeMsg(
        `⚠️ [TX](${explorerUrl}tx/${alert.source.block.hash}) Reported price of ${
            alert.metadata.cTokenAddress
        } was rejected\n Anchor Price: ${formatPrice(alert.metadata.anchorPrice)}\n Reporter Price: ${formatPrice(alert.metadata.reporterPrice)}`
    )
}

const formatDiscordAlert = (alert: FortaAlert) => {
    const embed = new EmbedBuilder()
        .setColor(0x0099ff)
        .setTitle('Oracle Price Monitor (Simulation)')
        .setURL(`${explorerUrl}tx/${alert.source.block.hash}`)
        .setAuthor({
            name: 'Forta (Simulation)',
            iconURL: 'https://github-production-user-asset-6210df.s3.amazonaws.com/4401444/250399693-e75d751d-5442-4bbb-b166-845e17c7393a.png',
            url: 'https://securityalliance.dev/',
        })
        .setDescription(formatTelegramAlert(alert))
        .setTimestamp()

    return embed.data
}

const sendTelegramAlert = async (chat: string, message: string) => {
    const res = await bot.sendMessage(chat, message, { parse_mode: 'MarkdownV2' })
    return res
}

const sendDiscordAlert = async (webhook: string, embeds: APIEmbed[]) => {
    const url = `https://discord.com/api/webhooks/${webhook}`
    const result = await axios.post(
        url,
        {
            embeds,
        },
        {
            headers: {
                'Content-Type': 'application/json',
            },
        }
    )
    return result.status
}

@JsonController()
export class NotificationController {
    @Get('/health')
    async health() {
        return 'I am alive'
    }

    @Get('/telegram')
    async sendTelegram(@QueryParam('message') message: string) {
        const status = await sendTelegramAlert(chatId, message)
        return 'This action sent a message: ' + message + ` with status: ` + status
    }

    @Get('/discord')
    async sendDiscord(@QueryParam('message') message: string) {
        const exampleEmbed = new EmbedBuilder()
            .setColor(0x0099ff)
            .setTitle('Oracle Price Monitor (Simulation)')
            .setAuthor({
                name: 'Forta (Simulation)',
                iconURL: 'https://github-production-user-asset-6210df.s3.amazonaws.com/4401444/250399693-e75d751d-5442-4bbb-b166-845e17c7393a.png',
                url: 'https://securityalliance.dev/',
            })
            .setDescription(message)
            .setTimestamp()

        await sendDiscordAlert(discordWebhookId, [exampleEmbed.data])
        return 'This action sent alert'
    }

    @Post('/forta')
    async post(@BodyParam('alerts') alerts: FortaAlert[]) {
        for await (const alert of alerts) {
            console.log({ alert })
            const message = formatTelegramAlert(alert)
            await sendTelegramAlert(chatId, message)
            const discordMessage = formatDiscordAlert(alert)
            await sendDiscordAlert(discordWebhookId, [discordMessage])
        }
        return 'This action sent alert'
    }
}

// at the top of your file

// inside a command, event listener, etc.

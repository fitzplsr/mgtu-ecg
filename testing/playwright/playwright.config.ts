import { defineConfig } from '@playwright/test'
import dotenv from 'dotenv'

dotenv.config()

export default defineConfig({
    use: {
        baseURL: process.env.API_BASE_URL ?? 'http://localhost:4000',
    },
    reporter: [
        ['list'],
        ['allure-playwright'],
    ],
})


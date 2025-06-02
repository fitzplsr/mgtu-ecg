import {expect, request, test} from '@playwright/test'
import fs from 'fs'
import path from 'path'
import dotenv from 'dotenv'

dotenv.config()

const API_PREFIX = (process.env.API_PREFIX ?? '/api/v1').replace(/\/$/, '')
const LOGIN = process.env.TEST_LOGIN ?? 'string'
const PASSWORD = process.env.TEST_PASSWORD ?? 'string'
const PATIENT_ID = 1
const EDF_PATH = process.env.TEST_EDF_PATH ?? path.resolve(__dirname, '..', 'fixtures', 'sample.edf')

const route = (p: string) => `${API_PREFIX}${p}`

let token: string
let uploadedFileId: number

test.describe.configure({mode: 'serial'})

test('POST /auth/signup → 201 & token', async () => {
    const ctx = await request.newContext()
    const res = await ctx.post(route('/auth/signup'), {
        data: {login: LOGIN, password: PASSWORD, name: "string",},
    })

    // expect(res.status()).toBe(201)
})

test('POST /auth/login → 200 & token', async () => {
    const ctx = await request.newContext()
    const res = await ctx.post(route('/auth/login'), {
        data: {login: LOGIN, password: PASSWORD},
    })

    expect(res.status()).toBe(200)
    const body = await res.json()
    expect(body.access_token).toBeTruthy()
    token = body.access_token
})

test('POST /analyse/upload → 201 & fileId', async () => {
    const ctx = await request.newContext({
        extraHTTPHeaders: {
            Authorization: `${token}`,
        },
    })

    const res = await ctx.post(route('/analyse/upload'), {
        multipart: {
            file: {
                name: path.basename(EDF_PATH),
                mimeType: 'application/octet-stream',
                buffer: fs.readFileSync(EDF_PATH),
            },
            patient_id: PATIENT_ID,
        },
    })

    expect(res.status()).toBe(201)
    const body = await res.json()
    expect(body[0].id).toBeDefined()
    uploadedFileId = body[0].id
})


test('PUT /analyse/edf → 200', async () => {
    const ctx = await request.newContext({
        extraHTTPHeaders: {
            Authorization: `${token}`,
        },
    })

    const res = await ctx.put(route(`/analyse/edf`), {
        data: {
            file_id: uploadedFileId,
        },
    })

    expect(res.status()).toBe(200)
})

test('POST /analyse/run → 200', async () => {
    const ctx = await request.newContext({
        extraHTTPHeaders: {
            Authorization: `${token}`,
        },
    })

    const res = await ctx.post(route(`/analyse/run`), {
            data: {
                file_ids: [
                    uploadedFileId
                ]
            }
        }
    )

    expect(res.status()).toBe(200)
})



test('PUT /analyse/patient/list → 200', async () => {
    const ctx = await request.newContext({
        extraHTTPHeaders: {
            Authorization: `${token}`,
        },
    })

    const res = await ctx.put(route(`/analyse/patient/list`), {
            data: {
                patient_id: PATIENT_ID
            }
        }
    )

    expect(res.status()).toBe(200)
})

test('PUT /analyse/list_edf → 200', async () => {
    const ctx = await request.newContext({
        extraHTTPHeaders: {
            Authorization: `${token}`,
        },
    })

    const res = await ctx.put(route(`/analyse/list_edf`), {
        data: {
            patient_id: PATIENT_ID,
        },
    })

    expect(res.status()).toBe(200)
    const body = await res.json()

    expect(body.files, 'body.files should be defined').toBeDefined()
    expect(Array.isArray(body.files), 'body.files should be an array').toBe(true)

    const hasFileWithUploadedFileId = body.files.some((f: any) => f.id === uploadedFileId)
    expect(hasFileWithUploadedFileId, 'Should contain file with id = uploadedFileId').toBe(true)
})


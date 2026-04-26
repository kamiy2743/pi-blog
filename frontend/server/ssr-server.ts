import { createServer, type IncomingMessage, type ServerResponse } from 'http'
import type { Page, PageProps } from '@inertiajs/core'
import { renderPage } from './render'

type JsonObject = Record<string, unknown>

type RouteHandler = (request: IncomingMessage) => Promise<JsonObject>

const portRaw = process.env.PORT
if (!portRaw) {
  throw new Error('.env に PORT が未設定です')
}
const port = Number.parseInt(portRaw, 10)
if (Number.isNaN(port)) {
  throw new Error(`PORT が不正です: ${portRaw}`)
}

const streamToString = (stream: IncomingMessage): Promise<string> =>
  new Promise((resolve, reject) => {
    let data = ''
    stream.on('data', (chunk: string | Buffer) => {
      data += chunk.toString()
    })
    stream.on('end', () => resolve(data))
    stream.on('error', (err: Error) => reject(err))
  })

const routes: Record<string, RouteHandler> = {
  '/health': async () => ({ status: 'OK', timestamp: Date.now() }),
  '/render': async (request: IncomingMessage) => {
    const requestBody = await streamToString(request)
    const payload = JSON.parse(requestBody) as Page<PageProps>
    const result = await renderPage(payload)

    return {
      head: result.head,
      body: result.body
    }
  }
}

createServer(async (request: IncomingMessage, response: ServerResponse) => {
  const urlPath = request.url ?? ''
  const handler = routes[urlPath]

  try {
    if (!handler) {
      response.writeHead(404, { 'Content-Type': 'application/json' })
      response.end(JSON.stringify({ status: 'NOT_FOUND', timestamp: Date.now() }))
      return
    }

    const result = await handler(request)
    response.writeHead(200, { 'Content-Type': 'application/json', Server: 'Inertia SSR' })
    response.end(JSON.stringify(result))
  } catch (err) {
    response.writeHead(500, { 'Content-Type': 'application/json' })
    response.end(JSON.stringify({ status: 'エラー', message: String(err) }))
  }
}).listen(port, () => {
  console.log(`Inertia SSR listening on :${port}`)
})

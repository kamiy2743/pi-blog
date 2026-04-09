import { createServer, type IncomingMessage, type ServerResponse } from 'http'
import { createInertiaApp, type ResolvedComponent } from '@inertiajs/svelte'
import type { Page, PageProps } from '@inertiajs/core'
import { render as renderSvelte } from 'svelte/server'

type InertiaSSRResponse = {
  head: string[]
  body: string
}

type JsonObject = Record<string, unknown>

type RouteHandler = (request: IncomingMessage) => Promise<JsonObject>

const pages = import.meta.glob<ResolvedComponent>('./pages/**/*.svelte', { eager: true })

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

const render = async (page: Page<PageProps>): Promise<InertiaSSRResponse> => {
  const result = await createInertiaApp({
    page,
    resolve: (name: string) => {
      const path = `./pages/${name}.svelte`
      const pageModule = pages[path]
      if (!pageModule) {
        throw new Error(`ページが見つかりません: ${name}`)
      }
      return pageModule
    },
    setup: ({ App, props }) => {
      return renderSvelte(App, { props })
    }
  })

  if (!result) {
    throw new Error('SSR レンダリング結果が空です')
  }

  return {
    head: result.head,
    body: result.body
  }
}

const routes: Record<string, RouteHandler> = {
  '/health': async () => ({ status: 'OK', timestamp: Date.now() }),
  '/render': async (request: IncomingMessage) => {
    const requestBody = await streamToString(request)
    const payload = JSON.parse(requestBody) as Page<PageProps>
    const result = await render(payload)
  
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

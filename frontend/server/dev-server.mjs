import { createServer as createHttpServer } from 'node:http'
import { createServer as createViteServer } from 'vite'

const streamToString = (stream) =>
  new Promise((resolve, reject) => {
    let data = ''
    stream.on('data', (chunk) => {
      data += chunk.toString()
    })
    stream.on('end', () => resolve(data))
    stream.on('error', (err) => reject(err))
  })

const httpServer = createHttpServer()

const vite = await createViteServer({
  configFile: './vite.config.vite-dev.js',
  server: {
    hmr: {
      server: httpServer,
    },
    middlewareMode: {
      server: httpServer,
    },
  },
})

httpServer.on('request', async (request, response) => {
  const urlPath = request.url ?? ''

  try {
    if (urlPath === '/health') {
      response.writeHead(200, { 'Content-Type': 'application/json' })
      response.end(JSON.stringify({ status: 'OK', timestamp: Date.now() }))
      return
    }

    if (urlPath === '/render' && request.method === 'POST') {
      const requestBody = await streamToString(request)
      const payload = JSON.parse(requestBody)
      const { renderPage } = await vite.ssrLoadModule('/server/render.ts')
      const result = await renderPage(payload)

      response.writeHead(200, { 'Content-Type': 'application/json', Server: 'Vite Inertia SSR' })
      response.end(JSON.stringify(result))
      return
    }

    vite.middlewares(request, response, () => {
      response.writeHead(404, { 'Content-Type': 'application/json' })
      response.end(JSON.stringify({ status: 'NOT_FOUND', timestamp: Date.now() }))
    })
  } catch (err) {
    vite.ssrFixStacktrace(err)
    response.writeHead(500, { 'Content-Type': 'application/json' })
    response.end(JSON.stringify({ status: 'エラー', message: String(err) }))
  }
})

const port = Number.parseInt(process.env.PORT, 10)

httpServer.listen(port, () => {
  console.log(`Vite dev server with SSR listening on :${port}`)
})

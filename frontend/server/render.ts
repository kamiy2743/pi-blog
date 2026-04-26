import { createInertiaApp, type ResolvedComponent } from '@inertiajs/svelte'
import type { Page, PageProps } from '@inertiajs/core'
import { render as renderSvelte } from 'svelte/server'

export type InertiaSSRResponse = {
  head: string[]
  body: string
}

const pages = import.meta.glob<ResolvedComponent>('../src/pages/**/*.svelte', { eager: true })

export async function renderPage(page: Page<PageProps>): Promise<InertiaSSRResponse> {
  const result = await createInertiaApp({
    page,
    resolve: (name: string) => {
      const path = `../src/pages/${name}.svelte`
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

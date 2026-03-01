import './app.css'
import { createInertiaApp } from '@inertiajs/svelte'
import type { ComponentType } from 'svelte'

type PageModule = {
  default: ComponentType
}

const pages = import.meta.glob<PageModule>('./pages/**/*.svelte')

createInertiaApp({
  resolve: async (name: string) => {
    const path = `./pages/${name}.svelte`
    const pageImporter = pages[path]
    if (!pageImporter) {
      throw new Error(`ページが見つかりません: ${name}`)
    }
    return pageImporter()
  },
  setup({ el, App, props }) {
    if (!(el instanceof HTMLElement)) {
      throw new Error('マウント先の要素が見つかりません')
    }
    new App({ target: el, props, hydrate: true })
  }
})

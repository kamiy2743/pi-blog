import './app.css'
import { createInertiaApp, type ResolvedComponent } from '@inertiajs/svelte'
import { hydrate } from 'svelte'

const pages = import.meta.glob<ResolvedComponent>('./pages/**/*.svelte')

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
    return hydrate(App, { target: el, props })
  }
})

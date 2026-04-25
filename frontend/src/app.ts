import './app.css'
import { createInertiaApp, type ResolvedComponent } from '@inertiajs/svelte'

const pages = import.meta.glob<ResolvedComponent>('./pages/**/*.svelte')

createInertiaApp({
  resolve: async (name: string) => {
    const path = `./pages/${name}.svelte`
    const pageImporter = pages[path]
    if (!pageImporter) {
      throw new Error(`ページが見つかりません: ${name}`)
    }
    return pageImporter()
  }
})

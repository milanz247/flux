// Flux application bootstrap. Reads the initial page object the Go server
// embedded in #app[data-page], resolves the matching Vue component from
// resources/js/Pages/**, and mounts a root that swaps pages on navigation.

import { createApp, defineComponent, h, shallowRef, type Component } from 'vue'
import { createPinia } from 'pinia'
import { initRouter, setCurrentPage } from './flux'
import type { Page } from './flux'
import { useAuthStore } from './stores/auth'
import '../css/app.css'

// Every page component, code-split by Vite. "Users/Index" →
// ./Pages/Users/Index.vue — exactly what req.View("Users/Index", ...) targets.
const pages = import.meta.glob<{ default: Component }>('./Pages/**/*.vue')

async function resolvePage(name: string): Promise<Component> {
  const importer = pages[`./Pages/${name}.vue`]
  if (!importer) {
    throw new Error(
      `Flux: page "${name}" not found. Create resources/js/Pages/${name}.vue ` +
        `(or run: flux make:page ${name})`,
    )
  }
  return (await importer()).default
}

const el = document.getElementById('app')
if (!el) throw new Error('Flux: missing #app element in the root template')

const initialPage: Page = JSON.parse(el.dataset.page ?? '{}')

void (async () => {
  const initialComponent = await resolvePage(initialPage.component)

  const component = shallowRef<Component>(initialComponent)
  const props = shallowRef(initialPage.props)

  const pinia = createPinia()

  const Root = defineComponent({
    name: 'FluxRoot',
    setup() {
      return () => h(component.value, { ...props.value, key: undefined })
    },
  })

  const app = createApp(Root)
  app.use(pinia)

  const auth = useAuthStore(pinia)

  const applyPage = (page: Page) => {
    setCurrentPage(page)
    auth.sync(page)
  }

  applyPage(initialPage)

  initRouter(initialPage, async (page) => {
    const next = await resolvePage(page.component)
    applyPage(page)
    component.value = next
    props.value = page.props
    document.title = `${page.component.split('/').join(' · ')} — ${auth.appName}`
  })

  app.mount(el)
})()

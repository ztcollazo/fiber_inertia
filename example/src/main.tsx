import { ComponentType, lazy, Suspense } from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import { createInertiaApp } from '@inertiajs/inertia-react'
import 'vite/modulepreload-polyfill'

const pages = import.meta.globEager('./pages/**/*.tsx')

createInertiaApp({
  resolve: name => pages[`./pages/${name}.tsx`],
  setup({ el, App, props }) {
    ReactDOM.render(<App {...props} />, el)
  },
})

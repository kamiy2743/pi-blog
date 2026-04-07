import '../app.css'

const reloadButton = document.querySelector<HTMLButtonElement>('[data-reload-button]')

reloadButton?.addEventListener('click', () => {
  window.location.reload()
})

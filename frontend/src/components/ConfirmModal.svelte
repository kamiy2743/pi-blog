<script lang="ts">
  import { Form } from '@inertiajs/svelte'

  export let title: string
  export let message: string
  export let action: string
  export let confirmLabel = 'OK'
  export let cancelLabel = 'キャンセル'
  export let confirmClass = 'admin-button'
  export let onClose: () => void
  export let onSuccess: () => void = onClose
</script>

<div
  class="fixed inset-0 z-50 bg-black/50"
  aria-hidden="true"
></div>
<div
  class="fixed inset-0 z-[60] flex items-center justify-center px-6"
  role="dialog"
  aria-modal="true"
>
  <div class="admin-panel w-full max-w-md rounded-lg border bg-[var(--admin-surface)] px-6 py-6 shadow-2xl">
    <h2 class="text-xl font-semibold">{title}</h2>
    <p class="admin-copy mt-3 text-sm leading-6">
      {message}
    </p>

    <Form
      class="mt-6 flex flex-wrap justify-center gap-3"
      {action}
      method="post"
      options={{ preserveScroll: true, preserveState: false }}
      {onSuccess}
    >
      <button
        class={`${confirmClass} min-w-28 rounded-lg px-5 py-3 text-sm font-semibold`}
        type="submit"
      >
        {confirmLabel}
      </button>
      <button
        class="admin-secondary-button min-w-28 rounded-lg border px-5 py-3 text-sm font-semibold"
        type="button"
        on:click={onClose}
      >
        {cancelLabel}
      </button>
    </Form>
  </div>
</div>

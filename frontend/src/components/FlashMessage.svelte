<script lang="ts">
  import { onDestroy } from 'svelte'
  import { fly } from 'svelte/transition'
  import type { Flash } from '../types/flash'

  type Props = {
    flash?: Flash
  }

  let { flash = {} }: Props = $props()
  let visible = $state(false)
  let hideTimer: ReturnType<typeof setTimeout> | null = null

  $effect(() => {
    const hasMessage = Boolean(flash.success || flash.error)
    if (hasMessage) {
      showMessage()
    } else {
      visible = false
      clearHideTimer()
    }
  })

  function showMessage() {
    visible = true
    clearHideTimer()
    hideTimer = setTimeout(hideMessage, 3000)
  }

  function hideMessage() {
    visible = false
    hideTimer = null
  }

  function clearHideTimer() {
    if (!hideTimer) {
      return
    }
    clearTimeout(hideTimer)
    hideTimer = null
  }

  onDestroy(clearHideTimer)
</script>

{#if visible}
  <div
    class="pointer-events-none fixed inset-x-0 bottom-4 z-50 flex justify-center px-4 sm:bottom-6"
    transition:fly={{ y: 16, duration: 180 }}
  >
    {#if flash.success}
      <div
        class="flash-message flash-message-success pointer-events-auto w-full max-w-xl rounded-lg border px-4 py-3 text-sm font-semibold shadow-lg"
        role="status"
        aria-live="polite"
      >
        {flash.success}
      </div>
    {/if}

    {#if flash.error}
      <div
        class="flash-message flash-message-error pointer-events-auto w-full max-w-xl rounded-lg border px-4 py-3 text-sm font-semibold shadow-lg"
        role="alert"
        aria-live="assertive"
      >
        {flash.error}
      </div>
    {/if}
  </div>
{/if}

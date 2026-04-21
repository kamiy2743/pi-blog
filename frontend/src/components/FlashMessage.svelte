<script lang="ts">
  import type { Flash } from '../types/flash'

  export let flash: Flash = {}
</script>

{#if flash.success}
  <div class="mb-6 rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm font-semibold text-emerald-700">
    {flash.success}
  </div>
{/if}

{#if flash.error}
  <div class="mb-6 rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm font-semibold text-red-700">
    {flash.error}
  </div>
{/if}

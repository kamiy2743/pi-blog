<script lang="ts">
  import type { Snippet } from 'svelte'
  import FlashMessage from './FlashMessage.svelte'
  import type { Flash } from '../types/flash'

  type Props = {
    flash?: Flash
    children?: Snippet
  }

  let { flash = {}, children }: Props = $props()
</script>

{#if children}
  {@render children()}
{/if}
<FlashMessage {flash} />

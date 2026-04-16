<script lang="ts">
  import { Link } from '@inertiajs/svelte'
  export let className = ''
</script>

<Link class={`blog-eyebrow mb-4 inline-flex text-sm font-semibold tracking-[0.24em] ${className}`.trim()} href="/">
  panda-dev.net / blog
</Link>

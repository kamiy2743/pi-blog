<svelte:head>
  <title>エラーが発生しました</title>
</svelte:head>

<script lang="ts">
  import { Link } from '@inertiajs/svelte'

  export let statusCode: number
  export let statusText: string
  export let message: string
  export let description: string

  function reloadPage(): void {
    window.location.reload()
  }
</script>

<div class="mx-auto max-w-5xl px-6 py-10 sm:px-8 lg:px-10 lg:py-14">
  <div class="blog-shell overflow-hidden rounded-[2rem] border">
    <div class="blog-hero px-6 py-10 sm:px-8 lg:px-10 lg:py-12">
      <div class="max-w-3xl space-y-4">
        <p class="blog-eyebrow text-sm font-semibold tracking-[0.24em]">{statusCode} / {statusText}</p>
        <h1 class="blog-title text-3xl font-semibold tracking-tight sm:text-3xl">
          {message}
        </h1>
        <p class="blog-copy max-w-2xl text-base leading-8 sm:text-lg">
          {description}
        </p>
        <div class="flex flex-wrap gap-3">
          {#if statusCode >= 500 && statusCode < 600}
            <button
              class="blog-button inline-flex items-center gap-2 rounded-full px-5 py-3 text-sm font-semibold transition"
              type="button"
              on:click={reloadPage}
            >
              再読み込み
              <span aria-hidden="true">↻</span>
            </button>
          {/if}
          <Link
            class="blog-button inline-flex items-center gap-2 rounded-full px-5 py-3 text-sm font-semibold transition"
            href="/"
          >
            トップページへ戻る
            <span aria-hidden="true">→</span>
          </Link>
        </div>
      </div>
    </div>
  </div>
</div>

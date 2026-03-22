<svelte:head>
  <title>Blog Top</title>
</svelte:head>

<script lang="ts">
  import { Link } from '@inertiajs/svelte'

  const latestArticles = [
    {
      title: 'Raspberry Pi 5 の常時稼働メモ',
      summary: '自宅サーバーとして安定運用するために見直した電源、冷却、監視のポイントを整理しました。',
      meta: '2025.02.10 / Infra'
    },
    {
      title: 'Go と SQLite で小さな CMS を組む',
      summary: '記事投稿の流れを単純に保つために、なぜこの構成にしたのかを実装視点でまとめています。',
      meta: '2025.02.04 / Backend'
    },
    {
      title: 'Cloudflare Tunnel で公開構成を整理する',
      summary: 'ポート公開を減らしつつ、外部公開と管理導線をどう分けるかを検証した記録です。',
      meta: '2025.01.28 / Network'
    }
  ]

  const categoryLinks = [
    { name: 'Infra', href: '/article?category=infra' },
    { name: 'Backend', href: '/article?category=backend' },
    { name: 'Network', href: '/article?category=network' },
    { name: 'Raspberry Pi', href: '/article?category=raspberry-pi' }
  ]
</script>

<div class="mx-auto max-w-7xl px-6 py-12 sm:px-8 lg:px-10 lg:py-16">
  <div class="overflow-hidden rounded-[2rem] border border-stone-200 bg-[#fffdf9] shadow-[0_18px_48px_rgba(28,25,23,0.08)]">
    <div class="border-b border-stone-200 bg-[#fcf6ee] px-6 py-10 sm:px-8 lg:px-10 lg:py-12">
      <div class="max-w-3xl space-y-4">
        <p class="text-sm font-semibold tracking-[0.24em] text-amber-800 uppercase">panda-dev.net / blog</p>
        <h1 class="text-4xl font-semibold tracking-tight text-stone-900 sm:text-5xl">
          静かな学びを記録する Blog
        </h1>
        <p class="max-w-2xl text-base leading-8 text-stone-600 sm:text-lg">
          Raspberry Pi 5 で動く自作ブログ。Go + SQLite を基盤に、
          小さく運用しながら学んだことを順番に残していきます。
        </p>
      </div>

      <div class="mt-6 flex flex-wrap gap-3">
        <span class="rounded-full border border-stone-200 bg-white px-3 py-1.5 text-sm font-medium text-stone-600">Go</span>
        <span class="rounded-full border border-stone-200 bg-white px-3 py-1.5 text-sm font-medium text-stone-600">SQLite</span>
        <span class="rounded-full border border-stone-200 bg-white px-3 py-1.5 text-sm font-medium text-stone-600">Inertia</span>
        <span class="rounded-full border border-stone-200 bg-white px-3 py-1.5 text-sm font-medium text-stone-600">Svelte</span>
        <span class="rounded-full border border-stone-200 bg-white px-3 py-1.5 text-sm font-medium text-stone-600">Raspberry Pi 5</span>
      </div>
    </div>

    <div class="grid gap-10 px-6 py-8 sm:px-8 lg:grid-cols-[minmax(0,1.7fr)_minmax(300px,0.9fr)] lg:gap-10 lg:px-10 lg:py-10">
      <section class="space-y-5">
        <div class="space-y-2">
          <p class="text-xs font-bold tracking-[0.24em] text-amber-800 uppercase">Latest</p>
          <h2 class="text-2xl font-semibold text-stone-900">最新記事</h2>
        </div>

        <div class="space-y-4">
          {#each latestArticles as article}
            <Link
              class="group block rounded-[1.5rem] border border-stone-200 bg-white px-6 py-6 transition duration-200 hover:-translate-y-0.5 hover:border-amber-200 hover:shadow-[0_16px_32px_rgba(120,53,15,0.08)]"
              href="/article"
            >
              <p class="text-sm font-medium text-amber-800">{article.meta}</p>
              <h3 class="mt-2 text-xl font-semibold text-stone-900 group-hover:text-amber-900">
                {article.title}
              </h3>
              <p class="mt-3 leading-8 text-stone-600">{article.summary}</p>
            </Link>
          {/each}
        </div>
      </section>

      <aside class="space-y-5">
        <div class="rounded-[1.5rem] border border-amber-100 bg-[#faf4ec] p-6">
          <p class="text-xs font-bold tracking-[0.24em] text-amber-800 uppercase">Navigation</p>
          <h2 class="mt-2 text-2xl font-semibold text-stone-900">記事一覧</h2>
          <p class="mt-3 leading-7 text-stone-600">公開済みの記事をまとめて確認できます。</p>
          <Link
            class="mt-5 inline-flex items-center gap-2 rounded-full bg-stone-900 px-5 py-3 text-sm font-semibold text-white transition hover:bg-amber-900"
            href="/article"
          >
            記事一覧へ
            <span aria-hidden="true">→</span>
          </Link>
        </div>

        <div class="rounded-[1.5rem] border border-amber-100 bg-[#faf4ec] p-6">
          <p class="text-xs font-bold tracking-[0.24em] text-amber-800 uppercase">Categories</p>
          <h2 class="mt-2 text-2xl font-semibold text-stone-900">カテゴリ</h2>
          <div class="mt-5 space-y-3">
            {#each categoryLinks as category}
              <Link
                class="block rounded-2xl border border-stone-200 bg-white px-4 py-3 font-semibold text-stone-700 transition hover:border-amber-200 hover:bg-amber-50 hover:text-amber-900"
                href={category.href}
              >
                {category.name}
              </Link>
            {/each}
          </div>
        </div>
      </aside>
    </div>
  </div>
</div>

<svelte:head>
  <title>Blog Top</title>
</svelte:head>

<script lang="ts">
  import { Link } from '@inertiajs/svelte'

  const latestArticles = [
    {
      title: 'Raspberry Pi 5 の常時稼働メモ',
      date: '2025.02.10',
      category: 'Infra'
    },
    {
      title: 'Go と SQLite で小さな CMS を組む',
      date: '2025.02.04',
      category: 'Backend'
    },
    {
      title: 'Cloudflare Tunnel で公開構成を整理する',
      date: '2025.01.28',
      category: 'Network'
    }
  ]

  const categoryLinks = [
    { name: 'Infra', href: '/article?category=infra' },
    { name: 'Backend', href: '/article?category=backend' },
    { name: 'Network', href: '/article?category=network' },
    { name: 'Raspberry Pi', href: '/article?category=raspberry-pi' }
  ]
</script>

<div class="mx-auto max-w-7xl px-6 py-10 sm:px-8 lg:px-10 lg:py-12">
  <div class="blog-shell overflow-hidden rounded-[2rem] border">
    <div class="blog-hero border-b px-6 py-10 sm:px-8 lg:px-10 lg:py-12">
      <div class="max-w-3xl space-y-4">
        <p class="blog-eyebrow text-sm font-semibold tracking-[0.24em]">panda-dev.net / blog</p>
        <h1 class="blog-title text-3xl font-semibold tracking-tight sm:text-3xl">
          パンダの個人開発ブログ
        </h1>
        <p class="blog-copy max-w-2xl text-base leading-8 sm:text-lg">
          Raspberry Pi で自宅サーバーを開発・運用しながら学んだことを発信していきます
        </p>
      </div>
    </div>

    <div class="grid gap-8 px-6 py-8 sm:px-8 lg:grid-cols-[minmax(100px,0.5fr)_minmax(0,1.7fr)] lg:gap-8 lg:px-10 lg:py-10">
      <aside class="space-y-5">
        <div class="blog-side-card rounded-[1.5rem] border p-6">
          <h2 class="blog-section-title text-2xl font-semibold">カテゴリ</h2>
          <div class="mt-5 space-y-3">
            {#each categoryLinks as category}
              <Link
                class="blog-category-link block rounded-2xl border px-4 py-3 font-semibold transition"
                href={category.href}
              >
                {category.name}
              </Link>
            {/each}
          </div>
          <Link
            class="blog-button mt-5 inline-flex items-center gap-2 rounded-full px-5 py-3 text-sm font-semibold transition"
            href="/article"
          >
            記事一覧へ
            <span aria-hidden="true">→</span>
          </Link>
        </div>
      </aside>

      <section class="space-y-5">
        <div class="space-y-2">
          <h2 class="blog-section-title text-2xl font-semibold">最新記事</h2>
        </div>

        <div class="space-y-4">
          {#each latestArticles as article}
            <Link
              class="blog-article-card group block rounded-[1.35rem] border px-5 py-5 transition"
              href="/article"
            >
              <p class="blog-accent text-sm font-semibold">{article.date}</p>
              <h3 class="blog-card-title mt-2 text-xl font-semibold">
                {article.title}
              </h3>
              <div class="mt-3">
                <span class="blog-category-pill px-3 py-1 text-sm font-semibold">{article.category}</span>
              </div>
            </Link>
          {/each}
        </div>
      </section>
    </div>
  </div>
</div>

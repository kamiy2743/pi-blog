<script lang="ts">
  import { Link, router } from '@inertiajs/svelte'
  import { ChevronRight } from 'lucide-svelte'
  import { formatDate } from '../../utils/date'

  type Article = {
    id: number
    title: string
    date: string
    isPublished: boolean
    categoryNames: string[]
  }

  type Category = {
    id: number
    name: string
  }

  type InitialProps = {
    categories: Category[]
  }

  type PartialSearchProps = {
    title: string
    categoryIds: number[]
    page: number
    totalCount: number
    totalPages: number
    articles: Article[]
  }

  export let initial: InitialProps = {
    categories: [],
  }
  export let partialSearch: PartialSearchProps = {
    title: '',
    categoryIds: [],
    page: 1,
    totalCount: 0,
    totalPages: 1,
    articles: [],
  }
  export let errors: Record<string, string> = {}

  let titleInput = partialSearch.title
  let selectedCategoryIds: number[] = [...partialSearch.categoryIds]

  function search() {
    router.visit(buildUrl(1), {
      method: 'get',
      preserveState: true,
      preserveScroll: true,
      replace: true,
      only: ['partialSearch', 'errors'],
    })
  }

  function buildUrl(targetPage: number): string {
    const params = new URLSearchParams()
    if (titleInput !== '') {
      params.set('title', titleInput)
    }
    for (const categoryId of selectedCategoryIds) {
      params.append('categoryId', String(categoryId))
    }
    params.set('page', String(targetPage))

    const query = params.toString()
    return query === '' ? '/admin' : `/admin?${query}`
  }

  function resetFilters() {
    titleInput = ''
    selectedCategoryIds = []
  }

  function paginationItems(currentPage: number, lastPage: number): Array<number | 'ellipsis'> {
    if (lastPage <= 7) {
      return Array.from({ length: lastPage }, (_, index) => index + 1)
    }

    if (currentPage <= 4) {
      return [1, 2, 3, 4, 5, 'ellipsis', lastPage]
    }

    if (currentPage >= lastPage - 3) {
      return [1, 'ellipsis', lastPage - 4, lastPage - 3, lastPage - 2, lastPage - 1, lastPage]
    }

    return [1, 'ellipsis', currentPage - 2, currentPage - 1, currentPage, currentPage + 1, currentPage + 2, 'ellipsis', lastPage]
  }
</script>

<svelte:head>
  <title>管理 - 記事一覧</title>
</svelte:head>

<div class="admin-page px-6 py-8 sm:px-8 lg:px-10 lg:py-10">
  <div class="mx-auto max-w-7xl">
    <header class="admin-hero rounded-lg border px-6 py-8 sm:px-8 lg:px-10">
      <div class="max-w-3xl space-y-3">
        <p class="admin-eyebrow text-sm font-semibold">Admin</p>
        <h1 class="text-3xl font-semibold tracking-tight sm:text-4xl">管理画面</h1>
        <p class="admin-copy text-base leading-8">
          記事の作成、編集、公開状態の確認を行います。
        </p>
      </div>
    </header>

    <main class="mt-6 grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
      <section class="space-y-4">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 class="text-2xl font-semibold">記事一覧</h2>
            <p class="admin-copy mt-1 text-sm">公開・非公開を含む記事を更新日時順で表示します。</p>
          </div>
        </div>

        <details class="admin-panel rounded-lg border px-5 py-5 group">
          <summary class="flex cursor-pointer list-none items-center gap-3">
            <ChevronRight class="h-4 w-4 transition-transform group-open:rotate-90" aria-hidden="true" />
            <h2 class="text-lg font-semibold">検索</h2>
          </summary>

          <form class="mt-5" action="/admin" method="get" on:submit|preventDefault={search}>
            <label class="block text-sm font-semibold" for="admin-title">タイトル</label>
            <input
              id="admin-title"
              name="title"
              class="mt-2 w-full rounded-lg border px-4 py-3"
              type="text"
              bind:value={titleInput}
              placeholder="Go, Docker など"
            />
            {#if errors.title}
              <p class="mt-2 text-sm font-semibold text-red-600">{errors.title}</p>
            {/if}

            <fieldset class="mt-5">
              <legend class="text-sm font-semibold">カテゴリ</legend>
              <div class="mt-3 flex flex-wrap gap-3">
                {#each initial.categories as category}
                  <label
                    class={`admin-filter-chip inline-flex items-center gap-2 rounded-lg border px-4 py-2 text-sm font-semibold transition ${
                      selectedCategoryIds.includes(category.id) ? 'admin-filter-chip-selected' : ''
                    }`}
                  >
                    <input
                      type="checkbox"
                      name="categoryId"
                      value={category.id}
                      bind:group={selectedCategoryIds}
                    />
                    <span>{category.name}</span>
                  </label>
                {/each}
              </div>
            </fieldset>

            <div class="mt-5 flex flex-wrap gap-3">
              <button
                class="admin-button inline-flex w-28 items-center justify-center rounded-lg px-5 py-3 text-sm font-semibold transition"
                type="submit"
              >
                検索
              </button>
              <button
                class="admin-secondary-button inline-flex w-28 items-center justify-center rounded-lg border px-5 py-3 text-sm font-semibold transition"
                type="button"
                on:click={resetFilters}
              >
                リセット
              </button>
            </div>
          </form>
        </details>

        <p class="admin-accent text-sm font-semibold">全 {partialSearch.totalCount} 件</p>

        {#if partialSearch.articles.length === 0}
          <div class="admin-panel rounded-lg border px-6 py-8">
            <p class="text-base font-semibold">該当する記事はありません。</p>
            <p class="admin-copy mt-2">検索条件を変えて、もう一度試してください。</p>
          </div>
        {:else}
          <div class="admin-panel overflow-x-auto rounded-lg border">
            <table class="min-w-full border-collapse text-left text-sm">
              <thead>
                <tr class="admin-table-head">
                  <th class="whitespace-nowrap px-4 py-3 font-semibold">更新日</th>
                  <th class="min-w-64 px-4 py-3 font-semibold">タイトル</th>
                  <th class="whitespace-nowrap px-4 py-3 font-semibold">状態</th>
                  <th class="min-w-48 px-4 py-3 font-semibold">カテゴリ</th>
                  <th class="whitespace-nowrap px-4 py-3 font-semibold">操作</th>
                </tr>
              </thead>
              <tbody>
                {#each partialSearch.articles as article}
                  <tr class="admin-table-row">
                    <td class="whitespace-nowrap px-4 py-3">{formatDate(article.date)}</td>
                    <td class="px-4 py-3 font-semibold">{article.title}</td>
                    <td class="whitespace-nowrap px-4 py-3">
                      <span
                        class={`admin-pill w-16 justify-center px-3 py-1 text-xs font-semibold ${
                          article.isPublished ? 'admin-pill-accent' : ''
                        }`}
                      >
                        {article.isPublished ? '公開' : '非公開'}
                      </span>
                    </td>
                    <td class="px-4 py-3">
                      <div class="flex flex-wrap gap-1.5">
                        {#each article.categoryNames as categoryName}
                          <span class="admin-pill px-2.5 py-1 text-xs font-semibold">{categoryName}</span>
                        {/each}
                      </div>
                    </td>
                    <td class="whitespace-nowrap px-4 py-3">
                      <Link
                        class="admin-secondary-button inline-flex items-center justify-center rounded-lg border px-4 py-2 text-sm font-semibold transition"
                        href={`/admin/article/edit/${article.id}`}
                      >
                        編集
                      </Link>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}

        {#if partialSearch.totalPages > 1}
          <div class="flex items-center justify-center gap-3 pt-2">
            {#if partialSearch.page > 1}
              <Link
                class="admin-secondary-button inline-flex h-10 w-10 items-center justify-center rounded-lg border text-sm font-semibold"
                href={buildUrl(partialSearch.page - 1)}
                aria-label="前のページ"
              >
                ←
              </Link>
            {/if}

            <div class="flex flex-wrap items-center justify-center gap-2">
              {#each paginationItems(partialSearch.page, partialSearch.totalPages) as item}
                {#if item === 'ellipsis'}
                  <span class="inline-flex h-10 min-w-10 items-center justify-center text-sm font-semibold">...</span>
                {:else if item === partialSearch.page}
                  <span class="admin-button inline-flex h-10 min-w-10 items-center justify-center rounded-lg px-3 text-sm font-semibold">
                    {item}
                  </span>
                {:else}
                  <Link
                    class="admin-secondary-button inline-flex h-10 min-w-10 items-center justify-center rounded-lg border px-3 text-sm font-semibold"
                    href={buildUrl(item)}
                    aria-label={`ページ ${item}`}
                  >
                    {item}
                  </Link>
                {/if}
              {/each}
            </div>

            {#if partialSearch.page < partialSearch.totalPages}
              <Link
                class="admin-secondary-button inline-flex h-10 w-10 items-center justify-center rounded-lg border text-sm font-semibold"
                href={buildUrl(partialSearch.page + 1)}
                aria-label="次のページ"
              >
                →
              </Link>
            {/if}
          </div>
        {/if}
      </section>

      <aside class="space-y-4">
        <section class="admin-panel rounded-lg border px-5 py-5">
          <h2 class="text-lg font-semibold">作成</h2>
          <p class="admin-copy mt-2 text-sm leading-6">新しい記事を追加します。</p>
          <Link
            class="admin-button mt-5 inline-flex w-full items-center justify-center rounded-lg px-5 py-3 text-sm font-semibold transition"
            href="/admin/article/new"
          >
            新規作成
          </Link>
        </section>

        <section class="admin-panel rounded-lg border px-5 py-5">
          <h2 class="text-lg font-semibold">公開側ページ</h2>
          <p class="admin-copy mt-2 text-sm leading-6">公開側ページの記事一覧を確認します。</p>
          <Link
            class="admin-secondary-button mt-5 inline-flex w-full items-center justify-center rounded-lg border px-5 py-3 text-sm font-semibold transition"
            href="/article"
          >
            公開記事一覧
          </Link>
        </section>
      </aside>
    </main>
  </div>
</div>

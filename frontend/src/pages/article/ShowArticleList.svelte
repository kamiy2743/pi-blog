<script lang="ts">
  import { Link, router } from '@inertiajs/svelte'
  import { ChevronRight } from 'lucide-svelte'
  import PublicSiteLink from '../../components/PublicSiteLink.svelte'
  import { formatDate } from '../../utils/date'

  type Article = {
    id: number
    title: string
    date: string
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
    return query === '' ? '/article' : `/article?${query}`
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
  <title>記事一覧</title>
</svelte:head>

<div class="mx-auto max-w-7xl px-6 py-10 sm:px-8 lg:px-10 lg:py-12">
  <div class="blog-shell overflow-hidden rounded-[2rem] border">
    <div class="blog-hero border-b px-6 py-10 sm:px-8 lg:px-10 lg:py-12">
      <div class="max-w-3xl space-y-4">
        <PublicSiteLink />
        <h1 class="blog-title text-3xl font-semibold tracking-tight sm:text-3xl">記事一覧</h1>
        <p class="blog-copy max-w-2xl text-base leading-8 sm:text-lg">
          公開済みの記事を一覧できます。タイトル検索とカテゴリ絞り込みに対応しています。
        </p>
      </div>
    </div>

    <div class="space-y-6 px-6 py-8 sm:px-8 lg:px-10 lg:py-10">
      <details class="blog-side-card mx-auto max-w-4xl rounded-[1.5rem] border p-6 group">
        <summary class="flex cursor-pointer list-none items-center gap-3">
          <ChevronRight class="h-4 w-4 transition-transform group-open:rotate-90" aria-hidden="true" />
          <h2 class="text-lg font-semibold text-slate-500">検索</h2>
        </summary>

        <form class="mt-5" action="/article" method="get" on:submit|preventDefault={search}>
          <label class="block text-sm font-semibold" for="title">タイトル</label>
          <input
            id="title"
            name="title"
            class="mt-2 w-full rounded-2xl border px-4 py-3"
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
                  class={`blog-filter-chip inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-semibold transition ${
                    selectedCategoryIds.includes(category.id)
                      ? 'blog-filter-chip-selected'
                      : ''
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

          <div class="mt-5 flex gap-3">
            <button
              class="blog-button inline-flex w-28 items-center justify-center gap-2 rounded-full px-5 py-3 text-sm font-semibold transition"
              type="submit"
            >
              検索
            </button>
            <button
              class="inline-flex w-28 items-center justify-center rounded-full border px-5 py-3 text-sm font-semibold transition"
              type="button"
              on:click={resetFilters}
            >
              リセット
            </button>
          </div>
        </form>
      </details>

      <section class="mx-auto max-w-4xl space-y-5">
        <div>
          <p class="blog-accent pl-3 text-sm font-semibold">全 {partialSearch.totalCount} 件</p>
        </div>

        {#if partialSearch.articles.length === 0}
          <div class="blog-side-card rounded-[1.5rem] border p-6">
            <p class="text-base font-semibold">該当する記事はありません。</p>
            <p class="muted mt-2">検索条件を変えて、もう一度試してください。</p>
          </div>
        {:else}
          <div class="space-y-4">
            {#each partialSearch.articles as article}
              <Link
                class="blog-article-card group block rounded-[1.35rem] border px-5 py-5 transition"
                href={`/article/${article.id}`}
              >
                <p class="blog-accent text-sm font-semibold">{formatDate(article.date)}</p>
                <h3 class="blog-card-title mt-2 text-xl font-semibold">
                  {article.title}
                </h3>
                <div class="mt-3 flex flex-wrap gap-2">
                  {#each article.categoryNames as categoryName}
                    <span class="blog-category-pill px-3 py-1 text-sm font-semibold">{categoryName}</span>
                  {/each}
                </div>
              </Link>
            {/each}
          </div>
        {/if}

        {#if partialSearch.totalPages > 1}
          <div class="flex items-center justify-center gap-3 pt-2">
            {#if partialSearch.page > 1}
              <Link
                class="inline-flex h-10 w-10 items-center justify-center rounded-full border text-sm font-semibold"
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
                  <span class="inline-flex h-10 min-w-10 items-center justify-center rounded-full border bg-stone-900 text-sm font-semibold text-stone-50">
                    {item}
                  </span>
                {:else}
                  <Link
                    class="inline-flex h-10 min-w-10 items-center justify-center rounded-full border text-sm font-semibold"
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
                class="inline-flex h-10 w-10 items-center justify-center rounded-full border text-sm font-semibold"
                href={buildUrl(partialSearch.page + 1)}
                aria-label="次のページ"
              >
                →
              </Link>
            {/if}
          </div>
        {/if}
      </section>
    </div>
  </div>
</div>

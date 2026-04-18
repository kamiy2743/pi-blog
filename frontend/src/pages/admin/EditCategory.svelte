<script lang="ts">
  type Category = {
    id: number
    name: string
  }

  export let categories: Category[] = []
  export let errors: Record<string, string> = {}

  let categoryNames: Record<number, string> = {}
  let deleteTarget: Category | null = null

  function categoryName(category: Category): string {
    return categoryNames[category.id] ?? category.name
  }

  function updateCategoryName(category: Category, event: Event) {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    categoryNames = {
      ...categoryNames,
      [category.id]: target.value,
    }
  }

  function isCategoryChanged(category: Category): boolean {
    return categoryName(category) !== category.name
  }

  function openDeleteModal(category: Category) {
    deleteTarget = category
  }

  function closeDeleteModal() {
    deleteTarget = null
  }
</script>

<svelte:head>
  <title>管理 - カテゴリ編集</title>
</svelte:head>

<div class="admin-page px-6 py-8 sm:px-8 lg:px-10 lg:py-10">
  <div class="mx-auto max-w-5xl">
    <div class="admin-panel rounded-lg border px-6 py-8 sm:px-8">
      <p class="admin-eyebrow text-sm font-semibold">Category</p>
      <h1 class="mt-3 text-3xl font-semibold tracking-tight">カテゴリ編集</h1>
      <p class="admin-copy mt-3 text-base leading-8">カテゴリの追加と名前の更新を行います。</p>

      <form class="mt-6 space-y-5" method="post" action="/admin/category">
        <label class="block">
          <span class="text-sm font-semibold">新しいカテゴリ名</span>
          <input
            class="mt-2 w-full rounded-lg border px-4 py-3"
            name="name"
            type="text"
            autocomplete="off"
          />
          {#if errors.name}
            <p class="mt-2 text-sm font-semibold text-red-600">{errors.name}</p>
          {/if}
        </label>

        <div class="flex flex-wrap gap-3">
          <button class="admin-button inline-flex items-center justify-center rounded-lg px-5 py-3 text-sm font-semibold" type="submit">
            追加
          </button>
        </div>
      </form>
    </div>

    <div class="admin-panel mt-6 rounded-lg border px-6 py-8 sm:px-8">
      <h2 class="text-xl font-semibold">登録済みカテゴリ</h2>

      {#if categories.length === 0}
        <p class="admin-copy mt-4">カテゴリはまだありません。</p>
      {:else}
        <div class="mt-5 space-y-4">
          {#each categories as category}
            <form
              class="grid gap-3 rounded-lg border border-[var(--admin-border)] p-4 sm:grid-cols-[minmax(0,1fr)_auto]"
              method="post"
              action={`/admin/category/${category.id}`}
            >
              <label class="block">
                <span class="text-sm font-semibold">カテゴリ名</span>
                <input
                  class="mt-2 w-full rounded-lg border px-4 py-3"
                  name="name"
                  type="text"
                  value={categoryName(category)}
                  autocomplete="off"
                  on:input={(event) => updateCategoryName(category, event)}
                />
                {#if errors.name}
                  <p class="mt-2 text-sm font-semibold text-red-600">{errors.name}</p>
                {/if}
              </label>
              <div class="flex flex-wrap gap-3 self-end">
                <button
                  class={`rounded-lg border px-5 py-3 text-sm font-semibold disabled:cursor-not-allowed disabled:opacity-40 ${
                    isCategoryChanged(category) ? 'admin-button' : 'admin-secondary-button'
                  }`}
                  type="submit"
                  disabled={!isCategoryChanged(category)}
                >
                  更新
                </button>
                <button
                  class="admin-secondary-button rounded-lg border px-5 py-3 text-sm font-semibold"
                  type="button"
                  on:click={() => openDeleteModal(category)}
                >
                  削除
                </button>
              </div>
            </form>
          {/each}
        </div>
      {/if}

      <a class="admin-secondary-button mt-6 inline-flex items-center justify-center rounded-lg border px-5 py-3 text-sm font-semibold" href="/admin">
        戻る
      </a>
    </div>
  </div>
</div>

{#if deleteTarget}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 px-6"
    role="dialog"
    aria-modal="true"
    aria-labelledby="delete-category-title"
  >
    <div class="admin-panel w-full max-w-md rounded-lg border px-6 py-6">
      <h2 id="delete-category-title" class="text-xl font-semibold">カテゴリを削除しますか</h2>
      <p class="admin-copy mt-3 text-sm leading-6">
        {deleteTarget.name} を削除します。
      </p>

      <form class="mt-6 flex flex-wrap gap-3" method="post" action={`/admin/category/${deleteTarget.id}/delete`}>
        <button
          class="admin-button rounded-lg px-5 py-3 text-sm font-semibold"
          type="submit"
        >
          削除
        </button>
        <button
          class="admin-secondary-button rounded-lg border px-5 py-3 text-sm font-semibold"
          type="button"
          on:click={closeDeleteModal}
        >
          キャンセル
        </button>
      </form>
    </div>
  </div>
{/if}

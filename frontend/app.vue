<template>
  <div class="app-shell">
    <header class="topbar">
      <div>
        <p class="eyebrow">
          Small Retail POS
        </p>
        <h1>POS and Inventory</h1>
      </div>
      <nav class="nav">
        <NuxtLink to="/dashboard">
          Dashboard
        </NuxtLink>
        <NuxtLink to="/products">
          Products
        </NuxtLink>
        <NuxtLink to="/pos">
          POS
        </NuxtLink>
        <NuxtLink to="/reports">
          Reports
        </NuxtLink>
      </nav>
    </header>

    <main class="page">
      <NuxtPage />
    </main>
  </div>
</template>

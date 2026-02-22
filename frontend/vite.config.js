import { sveltekit } from '@sveltejs/kit/vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';

const config = {
  plugins: [
    sveltekit(),
    SvelteKitPWA({
      registerType: 'autoUpdate',
      manifest: {
        name: 'Garden Planner',
        short_name: 'Garden',
        description: 'Local-first garden planning workspace.',
        theme_color: '#2f5d50',
        background_color: '#f6f1e9',
        display: 'standalone',
        start_url: '/',
        icons: [
          {
            src: '/icons/icon.svg',
            sizes: 'any',
            type: 'image/svg+xml'
          },
          {
            src: '/icons/icon-maskable.svg',
            sizes: 'any',
            type: 'image/svg+xml',
            purpose: 'maskable'
          }
        ]
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,svg,png,ico,txt,woff2}'],
        navigateFallback: '/index.html'
      },
      devOptions: {
        enabled: true
      }
    })
  ]
};

export default config;

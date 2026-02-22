/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        earth: {
          bg: '#F6F1E7',
          surface: '#EFE6D6',
          text: '#2B2A27',
          border: '#D8CBB5',
          forest: '#2F5D50',
          terracotta: '#B45A3C',
          mustard: '#C2A24C'
        }
      }
    }
  },
  plugins: []
};

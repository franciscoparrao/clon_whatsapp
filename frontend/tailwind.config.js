/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'whatsapp-green': '#25D366',
        'whatsapp-dark': '#075E54',
        'whatsapp-light': '#DCF8C6',
        'whatsapp-gray': '#ECE5DD',
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
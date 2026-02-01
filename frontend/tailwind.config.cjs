/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/**/*.{js,jsx,ts,tsx}',
    './public/index.html'
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f5f9ff',
          100: '#e6f0ff',
          500: '#2563eb'
        }
      },
      borderRadius: {
        'xl-2': '1rem'
      }
    }
  },
  plugins: []
}

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "./src/views/*.ejs" ],
  theme: {
    extend: {
      keyframes: {
        'popIn': {
          '0%': {
            'transform': 'scale(0.4)',
            'opacity': '0'
          },
          '100%': {
            'transform': 'scale(1)',
            'opacity': '1'
          }
        }
      },
      animation: {
        'popIn-600': 'popIn 600ms ease-in-out'
      }
    },
  },
  plugins: [],
}

const plugin = require('tailwindcss/plugin')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['**/*.{html,templ}'],
  daisyui: {
    themes: ['light', 'dark', 'cupcake'],
  },
  theme: {
    extend: {
      textShadow: {
        sm: '0 1px 2px var(--tw-shadow-color)',
        DEFAULT: '0 2px 4px var(--tw-shadow-color)',
        lg: '0 8px 16px var(--tw-shadow-color)',
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('daisyui'),
    plugin(function({ matchUtilities, theme }) {
      matchUtilities(
        {
          'text-shadow': (value) => ({
            textShadow: value,
          }),
        },
        { values: theme('textShadow') }
      )
    }),
  ],
}

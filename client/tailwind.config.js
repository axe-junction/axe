/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,jsx,ts,tsx}",
    "./presentation/components/**/*.{js,jsx,ts,tsx}",
  ], 
  presets: [require("nativewind/preset")],
  theme: {
    extend: {
      colors: {
        primary: '#5404FF', 
      },
    },
  },
  plugins: [],
};

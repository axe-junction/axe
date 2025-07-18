/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: "#6316DB",
        secondary: "#FFFFFF",
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "-apple-system", "sans-serif"],
      },
    },
  },
  plugins: [],
};

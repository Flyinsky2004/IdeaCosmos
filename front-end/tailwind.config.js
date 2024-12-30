/** @type {import('tailwindcss').Config} */
export default {
  darkMode:'class',
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      backgroundImage: {
        'glow-effect': `
          radial-gradient(circle, rgba(96,165,250,0.5), rgba(139,92,246,0.4), transparent),
          radial-gradient(circle, rgba(59,130,246,0.1), rgba(168,85,247,0.4), transparent),
          radial-gradient(circle, rgba(96,165,250,0.2), rgba(139,92,246,0.15), transparent)
        `
      },
    },
  },
  plugins: [],
}


/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#0069AA',
        secondary: '#FFA64D',
        gold: '#FFD700',
        silver: '#C0C0C0',
        bronze: '#CD7F32',
        'xp-green': '#58CC02',
        gem: '#1CB0F6',
      },
      fontFamily: {
        sans: ['Nunito', 'Inter', 'sans-serif'],
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-8px)' },
        },
        'glow-pulse': {
          '0%, 100%': { boxShadow: '0 0 5px rgba(0,105,170,0.4), 0 0 20px rgba(0,105,170,0.2)' },
          '50%': { boxShadow: '0 0 15px rgba(0,105,170,0.8), 0 0 40px rgba(0,105,170,0.4)' },
        },
        pop: {
          '0%': { transform: 'scale(0.95)', opacity: '0' },
          '70%': { transform: 'scale(1.03)' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
        'slide-up': {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
      animation: {
        float: 'float 3s ease-in-out infinite',
        'glow-pulse': 'glow-pulse 2s ease-in-out infinite',
        pop: 'pop 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards',
        shimmer: 'shimmer 2.5s linear infinite',
        'slide-up': 'slide-up 0.5s ease-out forwards',
      },
      screens: {
        'xs': '475px',
      },
    },
  },
  plugins: [],
}

